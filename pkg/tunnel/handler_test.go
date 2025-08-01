package tunnel

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestNewTunnelManager(t *testing.T) {
	tm := NewTunnelManager()
	if tm == nil {
		t.Fatal("NewTunnelManager returned nil")
	}
}

func TestTunnelManager_AddTunnel(t *testing.T) {
	tm := NewTunnelManager()
	tunnel := &Tunnel{
		ID:         "test-id",
		Subdomain:  "test",
		CreatedAt:  time.Now(),
		LastSeen:   time.Now(),
	}

	tm.AddTunnel(tunnel)

	// Check if tunnel was added
	retrieved, exists := tm.GetTunnel("test")
	if !exists {
		t.Fatal("Tunnel was not added")
	}
	if retrieved.ID != "test-id" {
		t.Fatal("Tunnel ID mismatch")
	}
}

func TestTunnelManager_RemoveTunnel(t *testing.T) {
	tm := NewTunnelManager()
	tunnel := &Tunnel{
		ID:         "test-id",
		Subdomain:  "test",
		CreatedAt:  time.Now(),
		LastSeen:   time.Now(),
	}

	tm.AddTunnel(tunnel)
	tm.RemoveTunnel("test")

	// Check if tunnel was removed
	_, exists := tm.GetTunnel("test")
	if exists {
		t.Fatal("Tunnel was not removed")
	}
}

func TestTunnelManager_ListTunnels(t *testing.T) {
	tm := NewTunnelManager()
	tunnel1 := &Tunnel{
		ID:         "test-id-1",
		Subdomain:  "test1",
		CreatedAt:  time.Now(),
		LastSeen:   time.Now(),
	}
	tunnel2 := &Tunnel{
		ID:         "test-id-2",
		Subdomain:  "test2",
		CreatedAt:  time.Now(),
		LastSeen:   time.Now(),
	}

	tm.AddTunnel(tunnel1)
	tm.AddTunnel(tunnel2)

	tunnels := tm.ListTunnels()
	if len(tunnels) != 2 {
		t.Fatalf("Expected 2 tunnels, got %d", len(tunnels))
	}
}

func TestNewHandler(t *testing.T) {
	tm := NewTunnelManager()
	logger := logrus.New()
	handler := NewHandler(tm, logger)

	if handler == nil {
		t.Fatal("NewHandler returned nil")
	}
}

func TestDefaultTunnelConfig(t *testing.T) {
	config := DefaultTunnelConfig()
	if config == nil {
		t.Fatal("DefaultTunnelConfig returned nil")
	}
	if config.LocalHost != "localhost" {
		t.Fatal("Default local host should be localhost")
	}
	if !config.UseTLS {
		t.Fatal("Default TLS should be enabled")
	}
}

func TestDefaultServerConfig(t *testing.T) {
	config := DefaultServerConfig()
	if config == nil {
		t.Fatal("DefaultServerConfig returned nil")
	}
	if config.Port != 443 {
		t.Fatal("Default port should be 443")
	}
	if !config.UseTLS {
		t.Fatal("Default TLS should be enabled")
	}
}

func TestTunnel_Close(t *testing.T) {
	tunnel := &Tunnel{
		ID:         "test-id",
		Subdomain:  "test",
		CreatedAt:  time.Now(),
		LastSeen:   time.Now(),
	}

	// Test closing an already closed tunnel
	err := tunnel.Close()
	if err != nil {
		t.Fatalf("Expected no error when closing tunnel, got %v", err)
	}

	// Test closing again
	err = tunnel.Close()
	if err != nil {
		t.Fatalf("Expected no error when closing already closed tunnel, got %v", err)
	}
}

func TestTunnel_IsClosed(t *testing.T) {
	tunnel := &Tunnel{
		ID:         "test-id",
		Subdomain:  "test",
		CreatedAt:  time.Now(),
		LastSeen:   time.Now(),
	}

	if tunnel.IsClosed() {
		t.Fatal("New tunnel should not be closed")
	}

	tunnel.Close()

	if !tunnel.IsClosed() {
		t.Fatal("Closed tunnel should be marked as closed")
	}
}

func TestTunnel_UpdateLastSeen(t *testing.T) {
	tunnel := &Tunnel{
		ID:         "test-id",
		Subdomain:  "test",
		CreatedAt:  time.Now(),
		LastSeen:   time.Now(),
	}

	oldLastSeen := tunnel.LastSeen
	time.Sleep(1 * time.Millisecond) // Ensure time difference
	tunnel.UpdateLastSeen()

	if tunnel.LastSeen.Equal(oldLastSeen) {
		t.Fatal("LastSeen should be updated")
	}
}

func TestCreateClientTLSConfig(t *testing.T) {
	config := CreateClientTLSConfig(false)
	if config == nil {
		t.Fatal("CreateClientTLSConfig returned nil")
	}
	if config.InsecureSkipVerify {
		t.Fatal("InsecureSkipVerify should be false when skipVerify is false")
	}

	config = CreateClientTLSConfig(true)
	if !config.InsecureSkipVerify {
		t.Fatal("InsecureSkipVerify should be true when skipVerify is true")
	}
} 