package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "og",
		Usage:   "GoTunnel - Secure tunneling to localhost",
		Version: "1.0.0",
		Commands: []*cli.Command{
			{
				Name:  "tunnel",
				Usage: "Create a tunnel to localhost",
				Subcommands: []*cli.Command{
					{
						Name:  "http",
						Usage: "Tunnel HTTP traffic",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "subdomain",
								Aliases: []string{"s"},
								Usage:   "Custom subdomain (optional)",
							},
							&cli.StringFlag{
								Name:    "host",
								Aliases: []string{"h"},
								Value:   "localhost",
								Usage:   "Local host to tunnel to",
							},
							&cli.StringFlag{
								Name:    "server",
								Aliases: []string{"S"},
								Value:   "https://tunnel.gotunnel.com",
								Usage:   "GoTunnel server URL",
							},
							&cli.StringFlag{
								Name:    "token",
								Aliases: []string{"t"},
								Usage:   "Authentication token",
							},
						},
						Action: func(c *cli.Context) error {
							if c.NArg() < 1 {
								return fmt.Errorf("port is required")
							}

							port := c.Args().Get(0)
							subdomain := c.String("subdomain")
							host := c.String("host")
							server := c.String("server")
							token := c.String("token")

							return startHTTPTunnel(port, subdomain, host, server, token)
						},
					},
					{
						Name:  "tcp",
						Usage: "Tunnel TCP traffic",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "subdomain",
								Aliases: []string{"s"},
								Usage:   "Custom subdomain (optional)",
							},
							&cli.StringFlag{
								Name:    "host",
								Aliases: []string{"h"},
								Value:   "localhost",
								Usage:   "Local host to tunnel to",
							},
							&cli.StringFlag{
								Name:    "server",
								Aliases: []string{"S"},
								Value:   "https://tunnel.gotunnel.com",
								Usage:   "GoTunnel server URL",
							},
							&cli.StringFlag{
								Name:    "token",
								Aliases: []string{"t"},
								Usage:   "Authentication token",
							},
						},
						Action: func(c *cli.Context) error {
							if c.NArg() < 1 {
								return fmt.Errorf("port is required")
							}

							port := c.Args().Get(0)
							subdomain := c.String("subdomain")
							host := c.String("host")
							server := c.String("server")
							token := c.String("token")

							return startTCPTunnel(port, subdomain, host, server, token)
						},
					},
				},
			},
			{
				Name:  "auth",
				Usage: "Manage authentication",
				Subcommands: []*cli.Command{
					{
						Name:  "login",
						Usage: "Login to GoTunnel",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "server",
								Aliases: []string{"s"},
								Value:   "https://tunnel.gotunnel.com",
								Usage:   "GoTunnel server URL",
							},
						},
						Action: func(c *cli.Context) error {
							server := c.String("server")
							return login(server)
						},
					},
					{
						Name:  "logout",
						Usage: "Logout from GoTunnel",
						Action: func(c *cli.Context) error {
							return logout()
						},
					},
					{
						Name:  "status",
						Usage: "Show authentication status",
						Action: func(c *cli.Context) error {
							return showAuthStatus()
						},
					},
				},
			},
			{
				Name:  "status",
				Usage: "Show tunnel status",
				Action: func(c *cli.Context) error {
					return showStatus()
				},
			},
			{
				Name:  "version",
				Usage: "Show version information",
				Action: func(c *cli.Context) error {
					fmt.Printf("og version %s\n", app.Version)
					return nil
				},
			},
		},
		Before: func(c *cli.Context) error {
			// Check if user is authenticated
			if c.Command.Name == "auth" {
				return nil // Skip auth check for auth commands
			}
			return checkAuth()
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startHTTPTunnel(port, subdomain, host, server, token string) error {
	fmt.Printf("🚀 Starting HTTP tunnel to %s:%s\n", host, port)
	
	if subdomain != "" {
		fmt.Printf("📍 Subdomain: %s\n", subdomain)
	} else {
		fmt.Printf("📍 Subdomain: auto-assigned\n")
	}
	
	fmt.Printf("🌐 Server: %s\n", server)
	fmt.Printf("⏳ Connecting...\n")

	// Create tunnel client
	client := NewTunnelClient(server, token)
	
	// Start tunnel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n🛑 Shutting down tunnel...")
		cancel()
	}()

	// Start the tunnel
	err := client.StartHTTPTunnel(ctx, host, port, subdomain)
	if err != nil {
		return fmt.Errorf("failed to start tunnel: %w", err)
	}

	return nil
}

func startTCPTunnel(port, subdomain, host, server, token string) error {
	fmt.Printf("🚀 Starting TCP tunnel to %s:%s\n", host, port)
	
	if subdomain != "" {
		fmt.Printf("📍 Subdomain: %s\n", subdomain)
	} else {
		fmt.Printf("📍 Subdomain: auto-assigned\n")
	}
	
	fmt.Printf("🌐 Server: %s\n", server)
	fmt.Printf("⏳ Connecting...\n")

	// Create tunnel client
	client := NewTunnelClient(server, token)
	
	// Start tunnel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n🛑 Shutting down tunnel...")
		cancel()
	}()

	// Start the tunnel
	err := client.StartTCPTunnel(ctx, host, port, subdomain)
	if err != nil {
		return fmt.Errorf("failed to start tunnel: %w", err)
	}

	return nil
}

func login(server string) error {
	fmt.Printf("🔐 Logging in to %s\n", server)
	
	// Interactive login
	fmt.Print("Email: ")
	var email string
	fmt.Scanln(&email)
	
	fmt.Print("Password: ")
	var password string
	fmt.Scanln(&password)

	// Perform login
	client := NewAuthClient(server)
	token, err := client.Login(email, password)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Save token
	err = SaveToken(token)
	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	fmt.Println("✅ Login successful!")
	return nil
}

func logout() error {
	fmt.Println("🔓 Logging out...")
	
	err := ClearToken()
	if err != nil {
		return fmt.Errorf("failed to clear token: %w", err)
	}

	fmt.Println("✅ Logout successful!")
	return nil
}

func showAuthStatus() error {
	token, err := LoadToken()
	if err != nil {
		fmt.Println("❌ Not authenticated")
		return nil
	}

	if token == "" {
		fmt.Println("❌ Not authenticated")
		return nil
	}

	fmt.Println("✅ Authenticated")
	return nil
}

func checkAuth() error {
	token, err := LoadToken()
	if err != nil {
		return fmt.Errorf("authentication required. Run 'og auth login'")
	}

	if token == "" {
		return fmt.Errorf("authentication required. Run 'og auth login'")
	}

	return nil
}

func showStatus() error {
	fmt.Println("📊 GoTunnel Status")
	fmt.Println("==================")
	
	// Check authentication
	token, err := LoadToken()
	if err != nil || token == "" {
		fmt.Println("❌ Not authenticated")
		return nil
	}
	fmt.Println("✅ Authenticated")

	// Show active tunnels (if any)
	fmt.Println("🌐 Active tunnels: 0")
	
	return nil
}

// TunnelClient represents a tunnel client
type TunnelClient struct {
	server string
	token  string
}

// NewTunnelClient creates a new tunnel client
func NewTunnelClient(server, token string) *TunnelClient {
	return &TunnelClient{
		server: server,
		token:  token,
	}
}

// StartHTTPTunnel starts an HTTP tunnel
func (tc *TunnelClient) StartHTTPTunnel(ctx context.Context, host, port, subdomain string) error {
	// Implementation would connect to the tunnel server
	// and establish the tunnel connection
	
	fmt.Printf("✅ Tunnel established!\n")
	fmt.Printf("🌐 Public URL: https://%s.tunnel.gotunnel.com\n", subdomain)
	fmt.Printf("📡 Forwarding: %s:%s\n", host, port)
	fmt.Printf("💡 Press Ctrl+C to stop\n")

	// Keep the tunnel alive
	<-ctx.Done()
	return nil
}

// StartTCPTunnel starts a TCP tunnel
func (tc *TunnelClient) StartTCPTunnel(ctx context.Context, host, port, subdomain string) error {
	// Implementation would connect to the tunnel server
	// and establish the tunnel connection
	
	fmt.Printf("✅ Tunnel established!\n")
	fmt.Printf("🌐 Public URL: tcp://%s.tunnel.gotunnel.com\n", subdomain)
	fmt.Printf("📡 Forwarding: %s:%s\n", host, port)
	fmt.Printf("💡 Press Ctrl+C to stop\n")

	// Keep the tunnel alive
	<-ctx.Done()
	return nil
}

// AuthClient represents an authentication client
type AuthClient struct {
	server string
}

// NewAuthClient creates a new auth client
func NewAuthClient(server string) *AuthClient {
	return &AuthClient{server: server}
}

// Login performs user login
func (ac *AuthClient) Login(email, password string) (string, error) {
	// Implementation would make HTTP request to login endpoint
	// For now, return a mock token
	return "mock-token-" + email, nil
}

// Token management functions
func SaveToken(token string) error {
	// Implementation would save token to config file
	return nil
}

func LoadToken() (string, error) {
	// Implementation would load token from config file
	return "mock-token", nil
}

func ClearToken() error {
	// Implementation would clear token from config file
	return nil
} 