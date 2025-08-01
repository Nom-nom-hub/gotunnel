package tunnel

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Handler implements the tunnel handler interface
type Handler struct {
	tunnelManager *TunnelManager
	logger        *logrus.Logger
}

// NewHandler creates a new tunnel handler
func NewHandler(tunnelManager *TunnelManager, logger *logrus.Logger) *Handler {
	return &Handler{
		tunnelManager: tunnelManager,
		logger:        logger,
	}
}

// HandleTunnel handles a new tunnel connection from a client
func (h *Handler) HandleTunnel(ctx context.Context, tunnel *Tunnel) error {
	h.logger.WithFields(logrus.Fields{
		"subdomain": tunnel.Subdomain,
		"id":        tunnel.ID,
	}).Info("New tunnel connection established")

	// Add tunnel to manager
	h.tunnelManager.AddTunnel(tunnel)
	defer func() {
		h.tunnelManager.RemoveTunnel(tunnel.Subdomain)
		tunnel.Close()
		h.logger.WithField("subdomain", tunnel.Subdomain).Info("Tunnel connection closed")
	}()

	// Keep connection alive and handle incoming data
	buffer := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Set read timeout
			if err := tunnel.ClientConn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
				return fmt.Errorf("failed to set read deadline: %w", err)
			}

			n, err := tunnel.ClientConn.Read(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Timeout, continue
					tunnel.UpdateLastSeen()
					continue
				}
				if err == io.EOF {
					return nil
				}
				return fmt.Errorf("failed to read from tunnel: %w", err)
			}

			if n > 0 {
				tunnel.UpdateLastSeen()
				h.logger.WithField("subdomain", tunnel.Subdomain).Debugf("Received %d bytes from tunnel", n)
			}
		}
	}
}

// HandleHTTPRequest handles incoming HTTP requests and forwards them to the appropriate tunnel
func (h *Handler) HandleHTTPRequest(ctx context.Context, subdomain string, conn net.Conn) error {
	h.logger.WithField("subdomain", subdomain).Info("Handling HTTP request")

	// Get the tunnel for this subdomain
	tunnel, exists := h.tunnelManager.GetTunnel(subdomain)
	if !exists {
		h.logger.WithField("subdomain", subdomain).Warn("Tunnel not found")
		return fmt.Errorf("tunnel not found for subdomain: %s", subdomain)
	}

	// Create a connection to handle the request
	connection := NewConnection(conn, tunnel.ClientConn)
	defer connection.Close()

	// Forward the request to the tunnel
	return h.forwardHTTPRequest(connection)
}

// forwardHTTPRequest forwards HTTP requests between client and server
func (h *Handler) forwardHTTPRequest(conn *Connection) error {
	// Create channels for bidirectional communication
	errChan := make(chan error, 2)

	// Forward from client to tunnel
	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.logger.WithField("error", r).Error("Panic in client to tunnel forwarding")
				errChan <- fmt.Errorf("panic in forwarding: %v", r)
			}
		}()

		written, err := io.Copy(conn.Server, conn.Client)
		if err != nil {
			h.logger.WithError(err).Debug("Error forwarding from client to tunnel")
		}
		h.logger.WithField("bytes_written", written).Debug("Forwarded from client to tunnel")
		errChan <- err
	}()

	// Forward from tunnel to client
	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.logger.WithField("error", r).Error("Panic in tunnel to client forwarding")
				errChan <- fmt.Errorf("panic in forwarding: %v", r)
			}
		}()

		written, err := io.Copy(conn.Client, conn.Server)
		if err != nil {
			h.logger.WithError(err).Debug("Error forwarding from tunnel to client")
		}
		h.logger.WithField("bytes_written", written).Debug("Forwarded from tunnel to client")
		errChan <- err
	}()

	// Wait for either direction to complete
	err1 := <-errChan
	err2 := <-errChan

	// Return the first error, or nil if both completed successfully
	if err1 != nil {
		return err1
	}
	return err2
}

// HandleRawTCP handles raw TCP connections
func (h *Handler) HandleRawTCP(ctx context.Context, subdomain string, conn net.Conn) error {
	h.logger.WithField("subdomain", subdomain).Info("Handling raw TCP connection")

	// Get the tunnel for this subdomain
	tunnel, exists := h.tunnelManager.GetTunnel(subdomain)
	if !exists {
		h.logger.WithField("subdomain", subdomain).Warn("Tunnel not found")
		return fmt.Errorf("tunnel not found for subdomain: %s", subdomain)
	}

	// Create a connection to handle the request
	connection := NewConnection(conn, tunnel.ClientConn)
	defer connection.Close()

	// Forward the raw TCP data
	return h.forwardRawTCP(connection)
}

// forwardRawTCP forwards raw TCP data between client and server
func (h *Handler) forwardRawTCP(conn *Connection) error {
	// Create channels for bidirectional communication
	errChan := make(chan error, 2)

	// Forward from client to tunnel
	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.logger.WithField("error", r).Error("Panic in client to tunnel forwarding")
				errChan <- fmt.Errorf("panic in forwarding: %v", r)
			}
		}()

		written, err := io.Copy(conn.Server, conn.Client)
		if err != nil {
			h.logger.WithError(err).Debug("Error forwarding from client to tunnel")
		}
		h.logger.WithField("bytes_written", written).Debug("Forwarded from client to tunnel")
		errChan <- err
	}()

	// Forward from tunnel to client
	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.logger.WithField("error", r).Error("Panic in tunnel to client forwarding")
				errChan <- fmt.Errorf("panic in forwarding: %v", r)
			}
		}()

		written, err := io.Copy(conn.Client, conn.Server)
		if err != nil {
			h.logger.WithError(err).Debug("Error forwarding from tunnel to client")
		}
		h.logger.WithField("bytes_written", written).Debug("Forwarded from tunnel to client")
		errChan <- err
	}()

	// Wait for either direction to complete
	err1 := <-errChan
	err2 := <-errChan

	// Return the first error, or nil if both completed successfully
	if err1 != nil {
		return err1
	}
	return err2
}

// DetectProtocol detects if the connection is HTTP or raw TCP
func (h *Handler) DetectProtocol(conn net.Conn) (string, error) {
	// Create a buffered reader to peek at the data
	reader := bufio.NewReader(conn)
	
	// Peek at the first few bytes to detect protocol
	peekBytes, err := reader.Peek(8)
	if err != nil {
		return "", fmt.Errorf("failed to peek connection: %w", err)
	}

	// Check if it looks like HTTP
	peekStr := strings.ToUpper(string(peekBytes))
	if strings.HasPrefix(peekStr, "GET ") || 
	   strings.HasPrefix(peekStr, "POST ") || 
	   strings.HasPrefix(peekStr, "PUT ") || 
	   strings.HasPrefix(peekStr, "DELETE ") || 
	   strings.HasPrefix(peekStr, "HEAD ") || 
	   strings.HasPrefix(peekStr, "OPTIONS ") || 
	   strings.HasPrefix(peekStr, "PATCH ") {
		return "http", nil
	}

	return "tcp", nil
}

// CreateTLSConfig creates a TLS configuration
func CreateTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	if certFile == "" || keyFile == "" {
		return nil, fmt.Errorf("cert file and key file are required for TLS")
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// CreateClientTLSConfig creates a TLS configuration for client connections
func CreateClientTLSConfig(skipVerify bool) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: skipVerify,
		MinVersion:         tls.VersionTLS12,
	}
} 