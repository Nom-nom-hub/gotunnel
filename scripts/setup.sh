#!/bin/bash

# GoTunnel Setup Script
# This script helps you set up GoTunnel for development or production

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21 or later."
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go version $GO_VERSION found"
}

# Check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_warning "Docker is not installed. Docker is optional but recommended for production deployment."
        return 1
    fi
    
    print_success "Docker found"
    return 0
}

# Generate certificates
generate_certs() {
    print_status "Generating TLS certificates..."
    
    if [ ! -d "certs" ]; then
        mkdir -p certs
    fi
    
    # Generate self-signed certificate
    openssl req -x509 -newkey rsa:4096 \
        -keyout certs/key.pem \
        -out certs/cert.pem \
        -days 365 \
        -nodes \
        -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
    
    print_success "Certificates generated in certs/ directory"
}

# Build binaries
build_binaries() {
    print_status "Building GoTunnel binaries..."
    
    # Create build directory
    mkdir -p build
    
    # Build server
    print_status "Building server..."
    go build -o build/gotunnel-server ./cmd/server
    
    # Build client
    print_status "Building client..."
    go build -o build/gotunnel-client ./cmd/client
    
    print_success "Binaries built successfully in build/ directory"
}

# Create configuration files
create_configs() {
    print_status "Creating configuration files..."
    
    # Create client config
    cat > config.yaml << EOF
# GoTunnel Client Configuration
server_addr: "localhost:8080"
subdomain: "myapp"
local_port: 3000
local_host: "localhost"
auth_token: "dev-token"
use_tls: false
skip_verify: true
EOF
    
    print_success "Configuration files created"
}

# Create startup scripts
create_scripts() {
    print_status "Creating startup scripts..."
    
    # Server startup script
    cat > start-server.sh << 'EOF'
#!/bin/bash
# GoTunnel Server Startup Script

echo "Starting GoTunnel Server..."

# Check if certificates exist
if [ ! -f "certs/cert.pem" ] || [ ! -f "certs/key.pem" ]; then
    echo "Error: TLS certificates not found. Run setup.sh first."
    exit 1
fi

# Start server
./build/gotunnel-server \
    --port 8080 \
    --cert certs/cert.pem \
    --key certs/key.pem \
    --allowed-tokens "dev-token" \
    --log-level info
EOF

    # Client startup script
    cat > start-client.sh << 'EOF'
#!/bin/bash
# GoTunnel Client Startup Script

echo "Starting GoTunnel Client..."

# Check if config exists
if [ ! -f "config.yaml" ]; then
    echo "Error: config.yaml not found. Run setup.sh first."
    exit 1
fi

# Start client
./build/gotunnel-client \
    --config config.yaml \
    --log-level info
EOF

    # Make scripts executable
    chmod +x start-server.sh start-client.sh
    
    print_success "Startup scripts created"
}

# Create Docker setup
setup_docker() {
    if ! check_docker; then
        return
    fi
    
    print_status "Setting up Docker configuration..."
    
    # Create .env file for Docker
    cat > .env << EOF
# GoTunnel Docker Environment Variables
ALLOWED_TOKENS=your-secret-token-here
SERVER_PORT=443
LOG_LEVEL=info
EOF
    
    print_success "Docker configuration created"
}

# Main setup function
main() {
    echo "=========================================="
    echo "GoTunnel Setup Script"
    echo "=========================================="
    echo ""
    
    # Check prerequisites
    print_status "Checking prerequisites..."
    check_go
    check_docker
    
    # Generate certificates
    generate_certs
    
    # Build binaries
    build_binaries
    
    # Create configuration files
    create_configs
    
    # Create startup scripts
    create_scripts
    
    # Setup Docker
    setup_docker
    
    echo ""
    echo "=========================================="
    print_success "Setup completed successfully!"
    echo "=========================================="
    echo ""
    echo "Next steps:"
    echo "1. Start the server: ./start-server.sh"
    echo "2. Start the client: ./start-client.sh"
    echo "3. Your local service will be available at: http://myapp.localhost:8080"
    echo ""
    echo "For production deployment:"
    echo "1. Update config.yaml with your server details"
    echo "2. Use proper TLS certificates"
    echo "3. Use strong authentication tokens"
    echo ""
    echo "For Docker deployment:"
    echo "1. docker-compose up -d"
    echo ""
}

# Run main function
main "$@" 