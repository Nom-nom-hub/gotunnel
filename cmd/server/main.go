package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ogrok/gotunnel/pkg/auth"
	"github.com/ogrok/gotunnel/pkg/tunnel"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/gorilla/websocket"
)

var (
	version = "1.0.0"
	logger  *logrus.Logger
)

func main() {
	app := &cli.App{
		Name:    "gotunnel-server",
		Version: version,
		Usage:   "Self-hosted tunnel server (ngrok alternative)",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   443,
				Usage:   "Port to listen on",
			},
			&cli.StringFlag{
				Name:    "cert",
				Aliases: []string{"c"},
				Usage:   "TLS certificate file",
			},
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "TLS private key file",
			},
			&cli.BoolFlag{
				Name:    "tls",
				Value:   true,
				Usage:   "Enable TLS",
			},
			&cli.StringSliceFlag{
				Name:    "allowed-tokens",
				Usage:   "Allowed authentication tokens",
			},
			&cli.StringFlag{
				Name:    "log-level",
				Value:   "info",
				Usage:   "Log level (debug, info, warn, error)",
			},
		},
		Action: runServer,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runServer(c *cli.Context) error {
	// Setup logging
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	level, err := logrus.ParseLevel(c.String("log-level"))
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	logger.SetLevel(level)

	// Create tunnel manager and handler
	tunnelManager := tunnel.NewTunnelManager()
	handler := tunnel.NewHandler(tunnelManager, logger)

	// Create authentication handler
	authHandler := auth.NewSimpleAuth()
	for _, token := range c.StringSlice("allowed-tokens") {
		authHandler.AddAllowedToken(token)
	}

	// Create server configuration
	config := tunnel.DefaultServerConfig()
	config.Port = c.Int("port")
	config.UseTLS = c.Bool("tls")
	config.TLSCertFile = c.String("cert")
	config.TLSKeyFile = c.String("key")

	// Create server
	server := NewServer(config, handler, authHandler, logger)

	// Start server
	logger.WithFields(logrus.Fields{
		"port": config.Port,
		"tls":  config.UseTLS,
	}).Info("Starting tunnel server")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Received shutdown signal, stopping server...")
		cancel()
	}()

	return server.Start(ctx)
}

// Server represents the tunnel server
type Server struct {
	config      *tunnel.ServerConfig
	handler     *tunnel.Handler
	authHandler *auth.SimpleAuth
	logger      *logrus.Logger
	httpServer  *http.Server
}

// NewServer creates a new tunnel server
func NewServer(config *tunnel.ServerConfig, handler *tunnel.Handler, authHandler *auth.SimpleAuth, logger *logrus.Logger) *Server {
	return &Server{
		config:      config,
		handler:     handler,
		authHandler: authHandler,
		logger:      logger,
	}
}

// Start starts the server
func (s *Server) Start(ctx context.Context) error {
	// Create listener
	listener, err := s.createListener()
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	defer listener.Close()

	s.logger.WithField("address", listener.Addr()).Info("Server listening")

	// Start HTTP server
	s.httpServer = &http.Server{
		Handler:      s.createHTTPHandler(),
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			s.logger.WithError(err).Error("HTTP server error")
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.logger.WithError(err).Error("Error during server shutdown")
	}

	s.logger.Info("Server stopped")
	return nil
}

// createListener creates the network listener
func (s *Server) createListener() (net.Listener, error) {
	addr := fmt.Sprintf(":%d", s.config.Port)

	if s.config.UseTLS {
		if s.config.TLSCertFile == "" || s.config.TLSKeyFile == "" {
			return nil, fmt.Errorf("TLS certificate and key files are required when TLS is enabled")
		}

		tlsConfig, err := tunnel.CreateTLSConfig(s.config.TLSCertFile, s.config.TLSKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create TLS config: %w", err)
		}

		return tls.Listen("tcp", addr, tlsConfig)
	}

	return net.Listen("tcp", addr)
}

// createHTTPHandler creates the HTTP handler
func (s *Server) createHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	// Handle tunnel client connections
	mux.HandleFunc("/tunnel", s.handleTunnelConnection)

	// Handle incoming HTTP requests
	mux.HandleFunc("/", s.handleIncomingRequest)

	return mux
}

// handleTunnelConnection handles tunnel client connections
func (s *Server) handleTunnelConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade to WebSocket for tunnel connection
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.WithError(err).Error("Failed to upgrade connection to WebSocket")
		return
	}

	// Get subdomain from query parameter
	subdomain := r.URL.Query().Get("subdomain")
	if subdomain == "" {
		s.logger.Error("No subdomain provided")
		conn.Close()
		return
	}

	// Get auth token from query parameter
	authToken := r.URL.Query().Get("token")
	if authToken == "" {
		s.logger.Error("No auth token provided")
		conn.Close()
		return
	}

	// Validate auth token
	if !s.authHandler.Authenticate(authToken) {
		s.logger.WithField("subdomain", subdomain).Error("Invalid auth token")
		conn.Close()
		return
	}

	// Create tunnel
	tunnel := &tunnel.Tunnel{
		ID:         generateID(),
		Subdomain:  subdomain,
		ClientConn: &WebSocketConn{conn: conn},
		CreatedAt:  time.Now(),
		LastSeen:   time.Now(),
	}

	// Handle tunnel in background
	go func() {
		if err := s.handler.HandleTunnel(context.Background(), tunnel); err != nil {
			s.logger.WithError(err).Error("Tunnel handler error")
		}
	}()
}

// handleIncomingRequest handles incoming HTTP requests
func (s *Server) handleIncomingRequest(w http.ResponseWriter, r *http.Request) {
	// Extract subdomain from Host header
	host := r.Host
	if host == "" {
		host = r.Header.Get("Host")
	}

	// Parse subdomain from host
	subdomain := s.extractSubdomain(host)
	if subdomain == "" {
		http.Error(w, "Invalid subdomain", http.StatusBadRequest)
		return
	}

	s.logger.WithFields(logrus.Fields{
		"subdomain": subdomain,
		"method":    r.Method,
		"path":      r.URL.Path,
	}).Info("Handling incoming request")

	// Create a connection to handle the request
	conn := &HTTPConn{
		request:  r,
		response: w,
		done:     make(chan struct{}),
	}

	// Handle the request
	if err := s.handler.HandleHTTPRequest(r.Context(), subdomain, conn); err != nil {
		s.logger.WithError(err).Error("Failed to handle HTTP request")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Wait for connection to complete
	<-conn.done
}

// extractSubdomain extracts subdomain from host
func (s *Server) extractSubdomain(host string) string {
	// Remove port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Split by dots
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return ""
	}

	// Return the first part as subdomain
	return parts[0]
}

// generateID generates a unique ID
func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// WebSocketConn wraps WebSocket connection to implement net.Conn
type WebSocketConn struct {
	conn *websocket.Conn
}

func (w *WebSocketConn) Read(b []byte) (n int, err error) {
	_, message, err := w.conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	copy(b, message)
	return len(message), nil
}

func (w *WebSocketConn) Write(b []byte) (n int, err error) {
	err = w.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (w *WebSocketConn) Close() error {
	return w.conn.Close()
}

func (w *WebSocketConn) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *WebSocketConn) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}

func (w *WebSocketConn) SetDeadline(t time.Time) error {
	return w.conn.SetReadDeadline(t)
}

func (w *WebSocketConn) SetReadDeadline(t time.Time) error {
	return w.conn.SetReadDeadline(t)
}

func (w *WebSocketConn) SetWriteDeadline(t time.Time) error {
	return w.conn.SetWriteDeadline(t)
}

// HTTPConn wraps HTTP request/response to implement net.Conn
type HTTPConn struct {
	request  *http.Request
	response http.ResponseWriter
	done     chan struct{}
}

func (h *HTTPConn) Read(b []byte) (n int, err error) {
	// For HTTP, we don't read from the connection
	return 0, nil
}

func (h *HTTPConn) Write(b []byte) (n int, err error) {
	_, err = h.response.Write(b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (h *HTTPConn) Close() error {
	close(h.done)
	return nil
}

func (h *HTTPConn) LocalAddr() net.Addr {
	return nil
}

func (h *HTTPConn) RemoteAddr() net.Addr {
	return nil
}

func (h *HTTPConn) SetDeadline(t time.Time) error {
	return nil
}

func (h *HTTPConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (h *HTTPConn) SetWriteDeadline(t time.Time) error {
	return nil
} 