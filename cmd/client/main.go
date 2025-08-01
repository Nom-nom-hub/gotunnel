package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ogrok/gotunnel/pkg/tunnel"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var (
	version = "1.0.0"
	logger  *logrus.Logger
)

// Config represents the client configuration
type Config struct {
	ServerAddr string `yaml:"server_addr" json:"server_addr"`
	Subdomain  string `yaml:"subdomain" json:"subdomain"`
	LocalPort  int    `yaml:"local_port" json:"local_port"`
	LocalHost  string `yaml:"local_host" json:"local_host"`
	AuthToken  string `yaml:"auth_token" json:"auth_token"`
	UseTLS     bool   `yaml:"use_tls" json:"use_tls"`
	SkipVerify bool   `yaml:"skip_verify" json:"skip_verify"`
}

func main() {
	app := &cli.App{
		Name:    "gotunnel-client",
		Version: version,
		Usage:   "Self-hosted tunnel client (ngrok alternative)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "server",
				Aliases:  []string{"s"},
				Required: true,
				Usage:    "Tunnel server address (e.g., tunnel.example.com)",
			},
			&cli.StringFlag{
				Name:     "subdomain",
				Aliases:  []string{"d"},
				Required: true,
				Usage:    "Subdomain for the tunnel (e.g., myapp)",
			},
			&cli.IntFlag{
				Name:     "local-port",
				Aliases:  []string{"p"},
				Required: true,
				Usage:    "Local port to forward",
			},
			&cli.StringFlag{
				Name:    "local-host",
				Value:   "localhost",
				Usage:   "Local host to forward",
			},
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Required: true,
				Usage:    "Authentication token",
			},
			&cli.BoolFlag{
				Name:    "tls",
				Value:   true,
				Usage:   "Use TLS for connection",
			},
			&cli.BoolFlag{
				Name:    "skip-verify",
				Usage:   "Skip TLS certificate verification",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Configuration file path",
			},
			&cli.StringFlag{
				Name:    "log-level",
				Value:   "info",
				Usage:   "Log level (debug, info, warn, error)",
			},
		},
		Action: runClient,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runClient(c *cli.Context) error {
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

	// Load configuration
	config, err := loadConfig(c)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create client
	client := NewClient(config, logger)

	// Start client
	logger.WithFields(logrus.Fields{
		"server":     config.ServerAddr,
		"subdomain":  config.Subdomain,
		"local_port": config.LocalPort,
		"local_host": config.LocalHost,
	}).Info("Starting tunnel client")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Received shutdown signal, stopping client...")
		cancel()
	}()

	return client.Start(ctx)
}

// loadConfig loads configuration from command line flags or config file
func loadConfig(c *cli.Context) (*Config, error) {
	config := &Config{}

	// Check if config file is provided
	if configFile := c.String("config"); configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with command line flags
	if c.IsSet("server") {
		config.ServerAddr = c.String("server")
	}
	if c.IsSet("subdomain") {
		config.Subdomain = c.String("subdomain")
	}
	if c.IsSet("local-port") {
		config.LocalPort = c.Int("local-port")
	}
	if c.IsSet("local-host") {
		config.LocalHost = c.String("local-host")
	}
	if c.IsSet("token") {
		config.AuthToken = c.String("token")
	}
	if c.IsSet("tls") {
		config.UseTLS = c.Bool("tls")
	}
	if c.IsSet("skip-verify") {
		config.SkipVerify = c.Bool("skip-verify")
	}

	// Validate required fields
	if config.ServerAddr == "" {
		return nil, fmt.Errorf("server address is required")
	}
	if config.Subdomain == "" {
		return nil, fmt.Errorf("subdomain is required")
	}
	if config.LocalPort == 0 {
		return nil, fmt.Errorf("local port is required")
	}
	if config.AuthToken == "" {
		return nil, fmt.Errorf("auth token is required")
	}

	return config, nil
}

// Client represents the tunnel client
type Client struct {
	config *Config
	logger *logrus.Logger
	conn   *websocket.Conn
}

// NewClient creates a new tunnel client
func NewClient(config *Config, logger *logrus.Logger) *Client {
	return &Client{
		config: config,
		logger: logger,
	}
}

// Start starts the client
func (c *Client) Start(ctx context.Context) error {
	// Connect to tunnel server
	if err := c.connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to tunnel server: %w", err)
	}
	defer c.conn.Close()

	c.logger.Info("Connected to tunnel server")

	// Start forwarding traffic
	return c.forwardTraffic(ctx)
}

// connect establishes connection to the tunnel server
func (c *Client) connect(ctx context.Context) error {
	// Build WebSocket URL
	scheme := "wss"
	if !c.config.UseTLS {
		scheme = "ws"
	}

	u := url.URL{
		Scheme: scheme,
		Host:   c.config.ServerAddr,
		Path:   "/tunnel",
		RawQuery: url.Values{
			"subdomain": []string{c.config.Subdomain},
			"token":     []string{c.config.AuthToken},
		}.Encode(),
	}

	c.logger.WithField("url", u.String()).Debug("Connecting to tunnel server")

	// Create dialer
	dialer := websocket.Dialer{
		HandshakeTimeout: 30 * time.Second,
	}

	// Add TLS configuration if needed
	if c.config.UseTLS {
		dialer.TLSClientConfig = tunnel.CreateClientTLSConfig(c.config.SkipVerify)
	}

	// Connect
	conn, _, err := dialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to establish WebSocket connection: %w", err)
	}

	c.conn = conn
	return nil
}

// forwardTraffic forwards traffic between local service and tunnel server
func (c *Client) forwardTraffic(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Read message from tunnel server
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					return fmt.Errorf("unexpected WebSocket close: %w", err)
				}
				return fmt.Errorf("failed to read message from tunnel: %w", err)
			}

			// Forward message to local service
			if err := c.forwardToLocal(message); err != nil {
				c.logger.WithError(err).Error("Failed to forward to local service")
				continue
			}
		}
	}
}

// forwardToLocal forwards data to the local service
func (c *Client) forwardToLocal(data []byte) error {
	// Connect to local service
	localAddr := fmt.Sprintf("%s:%d", c.config.LocalHost, c.config.LocalPort)
	conn, err := net.DialTimeout("tcp", localAddr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to local service: %w", err)
	}
	defer conn.Close()

	// Write data to local service
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("failed to write to local service: %w", err)
	}

	// Read response from local service
	response := make([]byte, 4096)
	n, err := conn.Read(response)
	if err != nil {
		return fmt.Errorf("failed to read from local service: %w", err)
	}

	// Send response back through tunnel
	if err := c.conn.WriteMessage(websocket.BinaryMessage, response[:n]); err != nil {
		return fmt.Errorf("failed to send response through tunnel: %w", err)
	}

	c.logger.WithFields(logrus.Fields{
		"bytes_sent":     len(data),
		"bytes_received": n,
	}).Debug("Forwarded traffic")

	return nil
}

// handleIncomingRequest handles incoming HTTP requests from the tunnel server
func (c *Client) handleIncomingRequest(req *http.Request) (*http.Response, error) {
	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Build local URL
	localURL := fmt.Sprintf("http://%s:%d%s", c.config.LocalHost, c.config.LocalPort, req.URL.Path)
	if req.URL.RawQuery != "" {
		localURL += "?" + req.URL.RawQuery
	}

	// Create new request
	localReq, err := http.NewRequest(req.Method, localURL, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create local request: %w", err)
	}

	// Copy headers
	for name, values := range req.Header {
		for _, value := range values {
			localReq.Header.Add(name, value)
		}
	}

	// Send request to local service
	resp, err := client.Do(localReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to local service: %w", err)
	}

	return resp, nil
}

// createLocalListener creates a local listener for testing
func (c *Client) createLocalListener() (net.Listener, error) {
	localAddr := fmt.Sprintf("%s:%d", c.config.LocalHost, c.config.LocalPort)
	return net.Listen("tcp", localAddr)
} 