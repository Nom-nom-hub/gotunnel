package tunnel

import (
	"context"
	"net"
	"sync"
	"time"
)

// Tunnel represents a client tunnel connection
type Tunnel struct {
	ID          string
	Subdomain   string
	ClientConn  net.Conn
	CreatedAt   time.Time
	LastSeen    time.Time
	mu          sync.RWMutex
	closed      bool
}

// TunnelManager handles multiple tunnel connections
type TunnelManager struct {
	tunnels map[string]*Tunnel
	mu      sync.RWMutex
}

// NewTunnelManager creates a new tunnel manager
func NewTunnelManager() *TunnelManager {
	return &TunnelManager{
		tunnels: make(map[string]*Tunnel),
	}
}

// AddTunnel adds a new tunnel to the manager
func (tm *TunnelManager) AddTunnel(tunnel *Tunnel) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tunnels[tunnel.Subdomain] = tunnel
}

// GetTunnel retrieves a tunnel by subdomain
func (tm *TunnelManager) GetTunnel(subdomain string) (*Tunnel, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	tunnel, exists := tm.tunnels[subdomain]
	return tunnel, exists
}

// RemoveTunnel removes a tunnel from the manager
func (tm *TunnelManager) RemoveTunnel(subdomain string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	delete(tm.tunnels, subdomain)
}

// ListTunnels returns all active tunnels
func (tm *TunnelManager) ListTunnels() []*Tunnel {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	tunnels := make([]*Tunnel, 0, len(tm.tunnels))
	for _, tunnel := range tm.tunnels {
		tunnels = append(tunnels, tunnel)
	}
	return tunnels
}

// Close closes the tunnel and removes it from the manager
func (t *Tunnel) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return nil
	}
	t.closed = true
	return t.ClientConn.Close()
}

// IsClosed checks if the tunnel is closed
func (t *Tunnel) IsClosed() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.closed
}

// UpdateLastSeen updates the last seen timestamp
func (t *Tunnel) UpdateLastSeen() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.LastSeen = time.Now()
}

// TunnelConfig holds configuration for tunnel connections
type TunnelConfig struct {
	ServerAddr    string
	Subdomain     string
	LocalPort     int
	LocalHost     string
	AuthToken     string
	UseTLS        bool
	SkipVerify    bool
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

// DefaultTunnelConfig returns default configuration
func DefaultTunnelConfig() *TunnelConfig {
	return &TunnelConfig{
		LocalHost:    "localhost",
		UseTLS:       true,
		SkipVerify:   false,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port           int
	TLSCertFile    string
	TLSKeyFile     string
	UseTLS         bool
	AllowedOrigins []string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:        443,
		UseTLS:      true,
		ReadTimeout: 30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

// TunnelHandler defines the interface for handling tunnel connections
type TunnelHandler interface {
	HandleTunnel(ctx context.Context, tunnel *Tunnel) error
	HandleHTTPRequest(ctx context.Context, subdomain string, conn net.Conn) error
}

// Connection represents a bidirectional connection between client and server
type Connection struct {
	ID       string
	Client   net.Conn
	Server   net.Conn
	Created  time.Time
	mu       sync.Mutex
	closed   bool
}

// NewConnection creates a new connection
func NewConnection(client, server net.Conn) *Connection {
	return &Connection{
		ID:      generateID(),
		Client:  client,
		Server:  server,
		Created: time.Now(),
	}
}

// Close closes both client and server connections
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	
	var errs []error
	if err := c.Client.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := c.Server.Close(); err != nil {
		errs = append(errs, err)
	}
	
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// IsClosed checks if the connection is closed
func (c *Connection) IsClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

// generateID generates a unique connection ID
func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
} 