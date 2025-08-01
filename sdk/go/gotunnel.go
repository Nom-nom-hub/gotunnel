package gotunnel

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Client represents a GoTunnel client
type Client struct {
	serverURL string
	token     string
	conn      *websocket.Conn
	config    *Config
}

// Config represents client configuration
type Config struct {
	ServerURL  string
	Token      string
	Subdomain  string
	LocalHost  string
	LocalPort  int
	Protocol   string // "http", "tcp"
	Timeout    time.Duration
	RetryCount int
}

// Tunnel represents a tunnel connection
type Tunnel struct {
	ID         string    `json:"id"`
	Subdomain  string    `json:"subdomain"`
	PublicURL  string    `json:"public_url"`
	LocalHost  string    `json:"local_host"`
	LocalPort  int       `json:"local_port"`
	Protocol   string    `json:"protocol"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	BytesSent  int64     `json:"bytes_sent"`
	BytesRecv  int64     `json:"bytes_recv"`
}

// NewClient creates a new GoTunnel client
func NewClient(config *Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}

	return &Client{
		serverURL: config.ServerURL,
		token:     config.Token,
		config:    config,
	}
}

// Connect establishes a connection to the GoTunnel server
func (c *Client) Connect(ctx context.Context) error {
	// Parse server URL
	u, err := url.Parse(c.serverURL)
	if err != nil {
		return fmt.Errorf("invalid server URL: %w", err)
	}

	// Convert to WebSocket URL
	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	u.Path = "/tunnel"

	// Create WebSocket dialer
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		HandshakeTimeout: c.config.Timeout,
	}

	// Add authentication header
	headers := http.Header{}
	if c.token != "" {
		headers.Set("Authorization", "Bearer "+c.token)
	}

	// Connect to server
	conn, _, err := dialer.Dial(u.String(), headers)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	c.conn = conn
	logrus.Infof("Connected to GoTunnel server: %s", c.serverURL)

	return nil
}

// CreateTunnel creates a new tunnel
func (c *Client) CreateTunnel(ctx context.Context) (*Tunnel, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("not connected to server")
	}

	// Prepare tunnel request
	request := map[string]interface{}{
		"action":     "create_tunnel",
		"subdomain":  c.config.Subdomain,
		"local_host": c.config.LocalHost,
		"local_port": c.config.LocalPort,
		"protocol":   c.config.Protocol,
	}

	// Send request
	if err := c.conn.WriteJSON(request); err != nil {
		return nil, fmt.Errorf("failed to send tunnel request: %w", err)
	}

	// Read response
	var response map[string]interface{}
	if err := c.conn.ReadJSON(&response); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if status, ok := response["status"].(string); ok && status == "error" {
		return nil, fmt.Errorf("server error: %v", response["message"])
	}

	// Parse tunnel data
	tunnelData, ok := response["tunnel"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	// Convert to Tunnel struct
	tunnel := &Tunnel{
		Subdomain:  c.config.Subdomain,
		LocalHost:  c.config.LocalHost,
		LocalPort:  c.config.LocalPort,
		Protocol:   c.config.Protocol,
		Status:     "active",
		CreatedAt:  time.Now(),
	}

	if id, ok := tunnelData["id"].(string); ok {
		tunnel.ID = id
	}
	if publicURL, ok := tunnelData["public_url"].(string); ok {
		tunnel.PublicURL = publicURL
	}

	logrus.Infof("Tunnel created: %s", tunnel.PublicURL)
	return tunnel, nil
}

// StartTunnel starts the tunnel and begins forwarding traffic
func (c *Client) StartTunnel(ctx context.Context, tunnel *Tunnel) error {
	if c.conn == nil {
		return fmt.Errorf("not connected to server")
	}

	logrus.Infof("Starting tunnel: %s -> %s:%d", tunnel.PublicURL, tunnel.LocalHost, tunnel.LocalPort)

	// Start forwarding traffic
	go c.forwardTraffic(ctx, tunnel)

	// Keep connection alive
	for {
		select {
		case <-ctx.Done():
			return c.Close()
		default:
			// Send heartbeat
			heartbeat := map[string]interface{}{
				"action": "heartbeat",
				"tunnel_id": tunnel.ID,
			}
			if err := c.conn.WriteJSON(heartbeat); err != nil {
				return fmt.Errorf("failed to send heartbeat: %w", err)
			}
			time.Sleep(30 * time.Second)
		}
	}
}

// forwardTraffic forwards traffic between the tunnel and local service
func (c *Client) forwardTraffic(ctx context.Context, tunnel *Tunnel) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Read message from server
			var message map[string]interface{}
			if err := c.conn.ReadJSON(&message); err != nil {
				logrus.Errorf("Failed to read message: %v", err)
				continue
			}

			// Handle different message types
			if action, ok := message["action"].(string); ok {
				switch action {
				case "forward_request":
					go c.handleForwardRequest(ctx, message, tunnel)
				case "tunnel_closed":
					logrus.Info("Tunnel closed by server")
					return
				default:
					logrus.Debugf("Unknown action: %s", action)
				}
			}
		}
	}
}

// handleForwardRequest handles a forward request from the server
func (c *Client) handleForwardRequest(ctx context.Context, message map[string]interface{}, tunnel *Tunnel) {
	// Extract request data
	requestData, ok := message["data"].(map[string]interface{})
	if !ok {
		logrus.Error("Invalid request data")
		return
	}

	// Connect to local service
	localAddr := fmt.Sprintf("%s:%d", tunnel.LocalHost, tunnel.LocalPort)
	localConn, err := net.Dial("tcp", localAddr)
	if err != nil {
		logrus.Errorf("Failed to connect to local service: %v", err)
		c.sendErrorResponse(message, "Failed to connect to local service")
		return
	}
	defer localConn.Close()

	// Forward the request
	if tunnel.Protocol == "http" {
		c.forwardHTTPRequest(ctx, message, localConn)
	} else {
		c.forwardTCPRequest(ctx, message, localConn)
	}
}

// forwardHTTPRequest forwards an HTTP request
func (c *Client) forwardHTTPRequest(ctx context.Context, message map[string]interface{}, localConn net.Conn) {
	// Implementation for HTTP forwarding
	// This would parse the HTTP request and forward it to the local service
	logrus.Debug("Forwarding HTTP request")
}

// forwardTCPRequest forwards a TCP request
func (c *Client) forwardTCPRequest(ctx context.Context, message map[string]interface{}, localConn net.Conn) {
	// Implementation for TCP forwarding
	// This would forward raw TCP data
	logrus.Debug("Forwarding TCP request")
}

// sendErrorResponse sends an error response to the server
func (c *Client) sendErrorResponse(originalMessage map[string]interface{}, errorMsg string) {
	response := map[string]interface{}{
		"action": "error_response",
		"request_id": originalMessage["request_id"],
		"error": errorMsg,
	}

	if err := c.conn.WriteJSON(response); err != nil {
		logrus.Errorf("Failed to send error response: %v", err)
	}
}

// Close closes the tunnel connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ListTunnels lists all tunnels for the authenticated user
func (c *Client) ListTunnels(ctx context.Context) ([]*Tunnel, error) {
	// Implementation for listing tunnels
	// This would make an HTTP request to the API
	return []*Tunnel{}, nil
}

// DeleteTunnel deletes a tunnel
func (c *Client) DeleteTunnel(ctx context.Context, tunnelID string) error {
	// Implementation for deleting a tunnel
	// This would make an HTTP request to the API
	return nil
}

// GetTunnelStats gets statistics for a tunnel
func (c *Client) GetTunnelStats(ctx context.Context, tunnelID string) (map[string]interface{}, error) {
	// Implementation for getting tunnel statistics
	// This would make an HTTP request to the API
	return map[string]interface{}{}, nil
}

// Convenience functions for common use cases

// HTTP creates a new HTTP tunnel client
func HTTP(serverURL, token, subdomain string, localPort int) *Client {
	return NewClient(&Config{
		ServerURL:  serverURL,
		Token:      token,
		Subdomain:  subdomain,
		LocalHost:  "localhost",
		LocalPort:  localPort,
		Protocol:   "http",
		Timeout:    30 * time.Second,
		RetryCount: 3,
	})
}

// TCP creates a new TCP tunnel client
func TCP(serverURL, token, subdomain string, localPort int) *Client {
	return NewClient(&Config{
		ServerURL:  serverURL,
		Token:      token,
		Subdomain:  subdomain,
		LocalHost:  "localhost",
		LocalPort:  localPort,
		Protocol:   "tcp",
		Timeout:    30 * time.Second,
		RetryCount: 3,
	})
}

// QuickStart provides a simple way to start a tunnel
func QuickStart(ctx context.Context, serverURL, token, subdomain string, localPort int) (*Tunnel, error) {
	client := HTTP(serverURL, token, subdomain, localPort)
	
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}
	defer client.Close()

	tunnel, err := client.CreateTunnel(ctx)
	if err != nil {
		return nil, err
	}

	return tunnel, client.StartTunnel(ctx, tunnel)
} 