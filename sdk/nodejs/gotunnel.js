/**
 * GoTunnel Node.js SDK
 * A simple and powerful Node.js client for GoTunnel
 */

const WebSocket = require('ws');
const https = require('https');
const http = require('http');
const net = require('net');

class TunnelConfig {
    constructor(options = {}) {
        this.serverUrl = options.serverUrl || 'https://tunnel.gotunnel.com';
        this.token = options.token || '';
        this.subdomain = options.subdomain || '';
        this.localHost = options.localHost || 'localhost';
        this.localPort = options.localPort || 3000;
        this.protocol = options.protocol || 'http'; // 'http' or 'tcp'
        this.timeout = options.timeout || 30000;
        this.retryCount = options.retryCount || 3;
    }
}

class Tunnel {
    constructor(data = {}) {
        this.id = data.id || '';
        this.subdomain = data.subdomain || '';
        this.publicUrl = data.publicUrl || '';
        this.localHost = data.localHost || 'localhost';
        this.localPort = data.localPort || 3000;
        this.protocol = data.protocol || 'http';
        this.status = data.status || 'active';
        this.createdAt = data.createdAt || new Date().toISOString();
        this.bytesSent = data.bytesSent || 0;
        this.bytesRecv = data.bytesRecv || 0;
    }
}

class GoTunnelClient {
    constructor(config) {
        this.config = config;
        this.websocket = null;
        this.isConnected = false;
        this.heartbeatInterval = null;
    }

    async connect() {
        return new Promise((resolve, reject) => {
            try {
                // Convert HTTP URL to WebSocket URL
                let wsUrl = this.config.serverUrl
                    .replace('https://', 'wss://')
                    .replace('http://', 'ws://');
                wsUrl = `${wsUrl}/tunnel`;

                // Prepare headers
                const headers = {};
                if (this.config.token) {
                    headers['Authorization'] = `Bearer ${this.config.token}`;
                }

                // Create WebSocket connection
                this.websocket = new WebSocket(wsUrl, {
                    headers,
                    rejectUnauthorized: false
                });

                this.websocket.on('open', () => {
                    console.log(`Connected to GoTunnel server: ${this.config.serverUrl}`);
                    this.isConnected = true;
                    resolve();
                });

                this.websocket.on('error', (error) => {
                    console.error('WebSocket error:', error);
                    reject(error);
                });

                this.websocket.on('close', () => {
                    console.log('Connection closed');
                    this.isConnected = false;
                    if (this.heartbeatInterval) {
                        clearInterval(this.heartbeatInterval);
                    }
                });

            } catch (error) {
                reject(error);
            }
        });
    }

    async createTunnel() {
        if (!this.isConnected) {
            throw new Error('Not connected to server');
        }

        return new Promise((resolve, reject) => {
            // Prepare tunnel request
            const request = {
                action: 'create_tunnel',
                subdomain: this.config.subdomain,
                localHost: this.config.localHost,
                localPort: this.config.localPort,
                protocol: this.config.protocol
            };

            // Send request
            this.websocket.send(JSON.stringify(request));

            // Handle response
            const messageHandler = (data) => {
                try {
                    const response = JSON.parse(data);
                    
                    // Check for errors
                    if (response.status === 'error') {
                        reject(new Error(`Server error: ${response.message}`));
                        return;
                    }

                    // Parse tunnel data
                    const tunnelData = response.tunnel || {};
                    
                    const tunnel = new Tunnel({
                        id: tunnelData.id || '',
                        subdomain: this.config.subdomain,
                        publicUrl: tunnelData.publicUrl || '',
                        localHost: this.config.localHost,
                        localPort: this.config.localPort,
                        protocol: this.config.protocol,
                        status: 'active',
                        createdAt: new Date().toISOString()
                    });

                    console.log(`Tunnel created: ${tunnel.publicUrl}`);
                    
                    // Remove message handler
                    this.websocket.removeListener('message', messageHandler);
                    resolve(tunnel);

                } catch (error) {
                    reject(error);
                }
            };

            this.websocket.on('message', messageHandler);
        });
    }

    async startTunnel(tunnel) {
        if (!this.isConnected) {
            throw new Error('Not connected to server');
        }

        console.log(`Starting tunnel: ${tunnel.publicUrl} -> ${tunnel.localHost}:${tunnel.localPort}`);

        // Start forwarding traffic
        this._startForwarding(tunnel);

        // Start heartbeat
        this.heartbeatInterval = setInterval(() => {
            if (this.isConnected) {
                const heartbeat = {
                    action: 'heartbeat',
                    tunnelId: tunnel.id
                };
                this.websocket.send(JSON.stringify(heartbeat));
            }
        }, 30000);

        // Keep connection alive
        return new Promise((resolve, reject) => {
            this.websocket.on('close', () => {
                if (this.heartbeatInterval) {
                    clearInterval(this.heartbeatInterval);
                }
                resolve();
            });

            this.websocket.on('error', (error) => {
                reject(error);
            });
        });
    }

    _startForwarding(tunnel) {
        this.websocket.on('message', async (data) => {
            try {
                const message = JSON.parse(data);
                const action = message.action;

                if (action === 'forward_request') {
                    await this._handleForwardRequest(message, tunnel);
                } else if (action === 'tunnel_closed') {
                    console.log('Tunnel closed by server');
                } else {
                    console.debug(`Unknown action: ${action}`);
                }
            } catch (error) {
                console.error('Error processing message:', error);
            }
        });
    }

    async _handleForwardRequest(message, tunnel) {
        try {
            // Extract request data
            const requestData = message.data || {};

            // Connect to local service
            const localConnection = await this._connectToLocal(tunnel.localHost, tunnel.localPort);

            try {
                // Forward the request based on protocol
                if (tunnel.protocol === 'http') {
                    await this._forwardHttpRequest(message, localConnection);
                } else {
                    await this._forwardTcpRequest(message, localConnection);
                }
            } finally {
                localConnection.end();
            }

        } catch (error) {
            console.error('Failed to handle forward request:', error);
            await this._sendErrorResponse(message, error.message);
        }
    }

    _connectToLocal(host, port) {
        return new Promise((resolve, reject) => {
            const connection = net.createConnection(port, host, () => {
                resolve(connection);
            });

            connection.on('error', (error) => {
                reject(error);
            });
        });
    }

    async _forwardHttpRequest(message, localConnection) {
        // Implementation for HTTP forwarding
        // This would parse the HTTP request and forward it to the local service
        console.debug('Forwarding HTTP request');
        
        // For now, just log the action
        const requestData = message.data || {};
        console.log('HTTP request:', requestData);
    }

    async _forwardTcpRequest(message, localConnection) {
        // Implementation for TCP forwarding
        // This would forward raw TCP data
        console.debug('Forwarding TCP request');
        
        // For now, just log the action
        const requestData = message.data || {};
        console.log('TCP request:', requestData);
    }

    async _sendErrorResponse(originalMessage, errorMsg) {
        const response = {
            action: 'error_response',
            requestId: originalMessage.requestId,
            error: errorMsg
        };

        try {
            this.websocket.send(JSON.stringify(response));
        } catch (error) {
            console.error('Failed to send error response:', error);
        }
    }

    async close() {
        if (this.websocket) {
            this.websocket.close();
        }
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
        }
    }

    async listTunnels() {
        // Implementation for listing tunnels
        // This would make an HTTP request to the API
        return [];
    }

    async deleteTunnel(tunnelId) {
        // Implementation for deleting a tunnel
        // This would make an HTTP request to the API
    }

    async getTunnelStats(tunnelId) {
        // Implementation for getting tunnel statistics
        // This would make an HTTP request to the API
        return {};
    }
}

// Convenience functions for common use cases

function httpTunnel(serverUrl, token, subdomain, localPort) {
    const config = new TunnelConfig({
        serverUrl,
        token,
        subdomain,
        localPort,
        protocol: 'http'
    });
    return new GoTunnelClient(config);
}

function tcpTunnel(serverUrl, token, subdomain, localPort) {
    const config = new TunnelConfig({
        serverUrl,
        token,
        subdomain,
        localPort,
        protocol: 'tcp'
    });
    return new GoTunnelClient(config);
}

async function quickStart(serverUrl, token, subdomain, localPort) {
    const client = httpTunnel(serverUrl, token, subdomain, localPort);
    
    await client.connect();
    
    try {
        const tunnel = await client.createTunnel();
        await client.startTunnel(tunnel);
        return tunnel;
    } finally {
        await client.close();
    }
}

// Example usage
async function main() {
    // Configuration
    const serverUrl = 'https://tunnel.gotunnel.com';
    const token = 'your-auth-token';
    const subdomain = 'myapp';
    const localPort = 3000;

    // Create client
    const client = httpTunnel(serverUrl, token, subdomain, localPort);

    try {
        // Connect to server
        await client.connect();

        // Create tunnel
        const tunnel = await client.createTunnel();
        console.log(`Tunnel created: ${tunnel.publicUrl}`);

        // Start tunnel
        await client.startTunnel(tunnel);

    } catch (error) {
        console.error('Error:', error);
    } finally {
        await client.close();
    }
}

// Export classes and functions
module.exports = {
    GoTunnelClient,
    TunnelConfig,
    Tunnel,
    httpTunnel,
    tcpTunnel,
    quickStart
};

// Run example if this file is executed directly
if (require.main === module) {
    main().catch(console.error);
} 