import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { tomorrow } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { 
  CodeBracketIcon,
  ServerIcon,
  ComputerDesktopIcon,
  CogIcon
} from '@heroicons/react/24/outline';

const CodeExamples: React.FC = () => {
  const [activeTab, setActiveTab] = useState(0);

  const examples = [
    {
      title: "Server Configuration",
      description: "Complete server setup with TLS and authentication",
      icon: ServerIcon,
      language: "go",
      code: `package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/ogrok/gotunnel/pkg/tunnel"
    "github.com/ogrok/gotunnel/pkg/auth"
)

func main() {
    // Create server configuration
    config := tunnel.DefaultServerConfig()
    config.Port = 443
    config.UseTLS = true
    config.CertFile = "server.crt"
    config.KeyFile = "server.key"
    
    // Initialize authentication
    auth := auth.NewSimpleAuth([]string{"your-secret-token"})
    
    // Create tunnel manager and handler
    manager := tunnel.NewTunnelManager()
    handler := tunnel.NewHandler(manager)
    
    // Create and start server
    server := &Server{
        config:  config,
        auth:    auth,
        handler: handler,
    }
    
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        cancel()
    }()
    
    if err := server.Start(ctx); err != nil {
        log.Fatal(err)
    }
}`,
      highlights: ["Server configuration", "TLS setup", "Authentication", "Graceful shutdown"]
    },
    {
      title: "Client Implementation",
      description: "Client-side tunnel connection and traffic forwarding",
      icon: ComputerDesktopIcon,
      language: "go",
      code: `package main

import (
    "context"
    "log"
    "net/http"
    "net/url"
    
    "github.com/gorilla/websocket"
    "github.com/ogrok/gotunnel/pkg/tunnel"
)

type Client struct {
    config *tunnel.TunnelConfig
    conn   *websocket.Conn
}

func (c *Client) Connect(ctx context.Context) error {
    // Create WebSocket URL
    u := url.URL{
        Scheme: "wss",
        Host:   c.config.ServerAddr,
        Path:   "/tunnel",
    }
    
    // Add query parameters
    q := u.Query()
    q.Set("subdomain", c.config.Subdomain)
    q.Set("token", c.config.AuthToken)
    u.RawQuery = q.Encode()
    
    // Establish WebSocket connection
    conn, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
    if err != nil {
        return err
    }
    
    c.conn = conn
    return nil
}

func (c *Client) ForwardTraffic(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Read message from WebSocket
            _, message, err := c.conn.ReadMessage()
            if err != nil {
                return err
            }
            
            // Forward to local service
            if err := c.forwardToLocal(message); err != nil {
                log.Printf("Forward error: %v", err)
            }
        }
    }
}

func (c *Client) forwardToLocal(data []byte) error {
    // Connect to local service
    resp, err := http.Post(
        fmt.Sprintf("http://%s:%d", c.config.LocalHost, c.config.LocalPort),
        "application/octet-stream",
        bytes.NewReader(data),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // Read response and send back through tunnel
    responseData, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    return c.conn.WriteMessage(websocket.BinaryMessage, responseData)
}`,
      highlights: ["WebSocket connection", "Traffic forwarding", "Local service integration", "Error handling"]
    },
    {
      title: "Configuration",
      description: "YAML configuration for both client and server",
      icon: CogIcon,
      language: "yaml",
      code: `# Server Configuration (server.yaml)
server:
  port: 443
  use_tls: true
  cert_file: "certs/server.crt"
  key_file: "certs/server.key"
  allowed_tokens:
    - "your-secret-token-1"
    - "your-secret-token-2"
  
  # Optional: Rate limiting
  rate_limit:
    requests_per_minute: 1000
    burst_size: 100
  
  # Optional: Logging
  logging:
    level: "info"
    format: "json"
    output: "stdout"

# Client Configuration (config.yaml)
client:
  server_addr: "your-server.com"
  subdomain: "myapp"
  local_port: 3000
  local_host: "localhost"
  auth_token: "your-secret-token"
  use_tls: true
  skip_verify: false
  
  # Optional: Connection settings
  connection:
    timeout: 30s
    keep_alive: true
    max_retries: 3
  
  # Optional: Logging
  logging:
    level: "info"
    format: "text"
    output: "stderr"`,
      highlights: ["Server settings", "Client configuration", "Security options", "Logging setup"]
    },
    {
      title: "Docker Setup",
      description: "Docker Compose configuration for production deployment",
      icon: ServerIcon,
      language: "yaml",
      code: `# docker-compose.yml
version: '3.8'

services:
  gotunnel-server:
    build: .
    container_name: gotunnel-server
    ports:
      - "443:443"
      - "80:80"
    volumes:
      - ./certs:/certs:ro
      - ./logs:/var/log/gotunnel
    environment:
      - ALLOWED_TOKENS=your-secret-token-1,your-secret-token-2
      - LOG_LEVEL=info
      - USE_TLS=true
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    
  # Optional: Nginx reverse proxy
  nginx:
    image: nginx:alpine
    container_name: gotunnel-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./certs:/etc/nginx/certs:ro
    depends_on:
      - gotunnel-server
    restart: unless-stopped

# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \\
    -a -installsuffix cgo \\
    -o gotunnel-server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/gotunnel-server .
RUN chmod +x gotunnel-server

EXPOSE 443
CMD ["./gotunnel-server"]`,
      highlights: ["Service definition", "Volume mounting", "Health checks", "Multi-stage build"]
    }
  ];

  return (
    <section id="examples" className="py-20 bg-dark-800/50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <h2 className="text-4xl md:text-5xl font-bold mb-6">
            Code <span className="text-gradient">Examples</span>
          </h2>
          <p className="text-xl text-gray-300 max-w-3xl mx-auto">
            Real-world examples and configurations to help you get started quickly.
          </p>
        </motion.div>

        {/* Tab Navigation */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.2 }}
          viewport={{ once: true }}
          className="flex justify-center mb-8"
        >
          <div className="flex space-x-2 bg-dark-700 rounded-lg p-1">
            {examples.map((example, index) => (
              <motion.button
                key={index}
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                onClick={() => setActiveTab(index)}
                className={`px-4 py-2 rounded-md text-sm font-medium transition-all duration-200 flex items-center space-x-2 ${
                  activeTab === index
                    ? 'bg-gradient-to-r from-tunnel-500 to-primary-600 text-white'
                    : 'text-gray-400 hover:text-white'
                }`}
              >
                <example.icon className="w-4 h-4" />
                <span>{example.title}</span>
              </motion.button>
            ))}
          </div>
        </motion.div>

        {/* Code Display */}
        <motion.div
          key={activeTab}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="grid grid-cols-1 lg:grid-cols-3 gap-8"
        >
          {/* Description */}
          <div className="lg:col-span-1">
            <motion.div
              initial={{ opacity: 0, x: -30 }}
              whileInView={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.8 }}
              viewport={{ once: true }}
              className="bg-dark-800 rounded-xl p-6 border border-gray-700 h-full"
            >
              <div className="flex items-center space-x-3 mb-4">
                <div className="w-10 h-10 rounded-lg bg-gradient-to-r from-tunnel-500 to-primary-600 flex items-center justify-center">
                  <examples[activeTab].icon className="w-5 h-5 text-white" />
                </div>
                <h3 className="text-xl font-semibold text-white">
                  {examples[activeTab].title}
                </h3>
              </div>
              
              <p className="text-gray-300 mb-6 leading-relaxed">
                {examples[activeTab].description}
              </p>
              
              <div className="space-y-3">
                <h4 className="text-tunnel-300 font-semibold">Key Features:</h4>
                <ul className="space-y-2">
                  {examples[activeTab].highlights.map((highlight, index) => (
                    <motion.li
                      key={index}
                      initial={{ opacity: 0, x: -10 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-center space-x-2 text-gray-300 text-sm"
                    >
                      <div className="w-2 h-2 bg-tunnel-400 rounded-full"></div>
                      <span>{highlight}</span>
                    </motion.li>
                  ))}
                </ul>
              </div>
            </motion.div>
          </div>

          {/* Code */}
          <div className="lg:col-span-2">
            <motion.div
              initial={{ opacity: 0, x: 30 }}
              whileInView={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.8 }}
              viewport={{ once: true }}
              className="bg-dark-800 rounded-xl border border-gray-700 overflow-hidden"
            >
              <div className="bg-dark-700 px-6 py-3 border-b border-gray-700">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <CodeBracketIcon className="w-5 h-5 text-gray-400" />
                    <span className="text-gray-300 font-medium">
                      {examples[activeTab].title}
                    </span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                    <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
                    <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                  </div>
                </div>
              </div>
              
              <div className="relative">
                <SyntaxHighlighter
                  language={examples[activeTab].language}
                  style={tomorrow}
                  customStyle={{
                    margin: 0,
                    padding: '1.5rem',
                    backgroundColor: '#1e293b',
                    fontSize: '14px',
                    lineHeight: '1.6',
                    fontFamily: 'JetBrains Mono, monospace'
                  }}
                  showLineNumbers={true}
                  wrapLines={true}
                >
                  {examples[activeTab].code}
                </SyntaxHighlighter>
              </div>
            </motion.div>
          </div>
        </motion.div>

        {/* Additional Examples */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.4 }}
          viewport={{ once: true }}
          className="mt-16 grid grid-cols-1 md:grid-cols-2 gap-6"
        >
          <motion.div
            whileHover={{ scale: 1.02 }}
            className="bg-dark-800 rounded-xl p-6 border border-gray-700"
          >
            <h3 className="text-lg font-semibold mb-3 text-tunnel-300">Makefile</h3>
            <div className="code-block">
              <pre className="text-gray-300 text-sm">
{`# Build both server and client
build:
	go build -o gotunnel-server ./cmd/server
	go build -o gotunnel-client ./cmd/client

# Generate certificates
certs:
	mkdir -p certs
	openssl req -x509 -newkey rsa:4096 -keyout certs/server.key \\
		-out certs/server.crt -days 365 -nodes \\
		-subj "/C=US/ST=State/L=City/O=Org/CN=localhost"

# Run tests
test:
	go test ./... -v

# Clean build artifacts
clean:
	rm -f gotunnel-server gotunnel-client`}
              </pre>
            </div>
          </motion.div>

          <motion.div
            whileHover={{ scale: 1.02 }}
            className="bg-dark-800 rounded-xl p-6 border border-gray-700"
          >
            <h3 className="text-lg font-semibold mb-3 text-primary-300">Health Check</h3>
            <div className="code-block">
              <pre className="text-gray-300 text-sm">
{`// Health check endpoint
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"tunnels": s.handler.manager.Count(),
		"uptime": time.Since(s.startTime).String(),
	})
}`}
              </pre>
            </div>
          </motion.div>
        </motion.div>
      </div>
    </section>
  );
};

export default CodeExamples; 