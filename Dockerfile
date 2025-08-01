# Multi-stage build for GoTunnel Server
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gotunnel-server ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S gotunnel && \
    adduser -u 1001 -S gotunnel -G gotunnel

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/gotunnel-server .

# Create certificates directory
RUN mkdir -p /certs && chown -R gotunnel:gotunnel /app /certs

# Switch to non-root user
USER gotunnel

# Expose port
EXPOSE 443

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:443/ || exit 1

# Default command
ENTRYPOINT ["./gotunnel-server"]

# Default arguments
CMD ["--port", "443", "--cert", "/certs/cert.pem", "--key", "/certs/key.pem"] 