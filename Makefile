# GoTunnel Makefile

# Variables
BINARY_SERVER=gotunnel-server
BINARY_CLIENT=gotunnel-client
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

.PHONY: all build clean test deps server client docker-build docker-run help

# Default target
all: clean build

# Build both server and client
build: server client

# Build server only
server:
	@echo "Building server..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_SERVER) ./cmd/server

# Build client only
client:
	@echo "Building client..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_CLIENT) ./cmd/client

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	$(GOGET) github.com/gorilla/websocket
	$(GOGET) github.com/sirupsen/logrus
	$(GOGET) github.com/urfave/cli/v2
	$(GOGET) gopkg.in/yaml.v3

# Generate certificates for testing
certs:
	@echo "Generating test certificates..."
	@mkdir -p certs
	openssl req -x509 -newkey rsa:4096 -keyout certs/key.pem -out certs/cert.pem -days 365 -nodes -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

# Run server in development mode
dev-server: server
	@echo "Starting server in development mode..."
	./$(BUILD_DIR)/$(BINARY_SERVER) --port 8080 --cert certs/cert.pem --key certs/key.pem --allowed-tokens "dev-token" --log-level debug

# Run client in development mode
dev-client: client
	@echo "Starting client in development mode..."
	./$(BUILD_DIR)/$(BINARY_CLIENT) --server localhost:8080 --subdomain test --local-port 3000 --token "dev-token" --skip-verify --log-level debug

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t gotunnel:$(VERSION) .

# Docker run server
docker-run:
	@echo "Running server in Docker..."
	docker run -p 443:443 -v $(PWD)/certs:/certs gotunnel:$(VERSION) ./gotunnel-server --port 443 --cert /certs/cert.pem --key /certs/key.pem --allowed-tokens "docker-token"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Security check
security:
	@echo "Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Create release
release: clean build
	@echo "Creating release..."
	@mkdir -p release
	@tar -czf release/gotunnel-$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR) .
	@echo "Release created: release/gotunnel-$(VERSION)-linux-amd64.tar.gz"

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Clean and build both server and client"
	@echo "  build        - Build both server and client"
	@echo "  server       - Build server only"
	@echo "  client       - Build client only"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  deps         - Download dependencies"
	@echo "  install-deps - Install dependencies"
	@echo "  certs        - Generate test certificates"
	@echo "  dev-server   - Run server in development mode"
	@echo "  dev-client   - Run client in development mode"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run server in Docker"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  security     - Run security checks"
	@echo "  release      - Create release package"
	@echo "  help         - Show this help message" 