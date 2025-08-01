#!/usr/bin/env python3
"""
GoTunnel Python SDK
A simple and powerful Python client for GoTunnel
"""

import asyncio
import json
import logging
import time
from dataclasses import dataclass
from typing import Dict, List, Optional, Any
import websockets
import aiohttp
import ssl

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


@dataclass
class TunnelConfig:
    """Configuration for a tunnel"""
    server_url: str
    token: str
    subdomain: str
    local_host: str = "localhost"
    local_port: int = 3000
    protocol: str = "http"  # "http" or "tcp"
    timeout: int = 30
    retry_count: int = 3


@dataclass
class Tunnel:
    """Represents a tunnel connection"""
    id: str
    subdomain: str
    public_url: str
    local_host: str
    local_port: int
    protocol: str
    status: str
    created_at: str
    bytes_sent: int = 0
    bytes_recv: int = 0


class GoTunnelClient:
    """GoTunnel client for Python"""
    
    def __init__(self, config: TunnelConfig):
        self.config = config
        self.websocket = None
        self.session = None
        
    async def connect(self) -> None:
        """Connect to the GoTunnel server"""
        # Convert HTTP URL to WebSocket URL
        ws_url = self.config.server_url.replace("https://", "wss://").replace("http://", "ws://")
        ws_url = f"{ws_url}/tunnel"
        
        # Prepare headers
        headers = {}
        if self.config.token:
            headers["Authorization"] = f"Bearer {self.config.token}"
        
        # Create SSL context
        ssl_context = ssl.create_default_context()
        ssl_context.check_hostname = False
        ssl_context.verify_mode = ssl.CERT_NONE
        
        # Connect to WebSocket
        self.websocket = await websockets.connect(
            ws_url,
            extra_headers=headers,
            ssl=ssl_context,
            ping_interval=30,
            ping_timeout=10
        )
        
        logger.info(f"Connected to GoTunnel server: {self.config.server_url}")
    
    async def create_tunnel(self) -> Tunnel:
        """Create a new tunnel"""
        if not self.websocket:
            raise Exception("Not connected to server")
        
        # Prepare tunnel request
        request = {
            "action": "create_tunnel",
            "subdomain": self.config.subdomain,
            "local_host": self.config.local_host,
            "local_port": self.config.local_port,
            "protocol": self.config.protocol
        }
        
        # Send request
        await self.websocket.send(json.dumps(request))
        
        # Read response
        response_data = await self.websocket.recv()
        response = json.loads(response_data)
        
        # Check for errors
        if response.get("status") == "error":
            raise Exception(f"Server error: {response.get('message')}")
        
        # Parse tunnel data
        tunnel_data = response.get("tunnel", {})
        
        tunnel = Tunnel(
            id=tunnel_data.get("id", ""),
            subdomain=self.config.subdomain,
            public_url=tunnel_data.get("public_url", ""),
            local_host=self.config.local_host,
            local_port=self.config.local_port,
            protocol=self.config.protocol,
            status="active",
            created_at=time.strftime("%Y-%m-%d %H:%M:%S")
        )
        
        logger.info(f"Tunnel created: {tunnel.public_url}")
        return tunnel
    
    async def start_tunnel(self, tunnel: Tunnel) -> None:
        """Start the tunnel and begin forwarding traffic"""
        if not self.websocket:
            raise Exception("Not connected to server")
        
        logger.info(f"Starting tunnel: {tunnel.public_url} -> {tunnel.local_host}:{tunnel.local_port}")
        
        # Start forwarding traffic
        asyncio.create_task(self._forward_traffic(tunnel))
        
        # Keep connection alive with heartbeats
        while True:
            try:
                heartbeat = {
                    "action": "heartbeat",
                    "tunnel_id": tunnel.id
                }
                await self.websocket.send(json.dumps(heartbeat))
                await asyncio.sleep(30)
            except websockets.exceptions.ConnectionClosed:
                logger.info("Connection closed")
                break
            except Exception as e:
                logger.error(f"Heartbeat error: {e}")
                break
    
    async def _forward_traffic(self, tunnel: Tunnel) -> None:
        """Forward traffic between tunnel and local service"""
        while True:
            try:
                # Read message from server
                message_data = await self.websocket.recv()
                message = json.loads(message_data)
                
                action = message.get("action")
                if action == "forward_request":
                    asyncio.create_task(self._handle_forward_request(message, tunnel))
                elif action == "tunnel_closed":
                    logger.info("Tunnel closed by server")
                    break
                else:
                    logger.debug(f"Unknown action: {action}")
                    
            except websockets.exceptions.ConnectionClosed:
                logger.info("Connection closed")
                break
            except Exception as e:
                logger.error(f"Error reading message: {e}")
                continue
    
    async def _handle_forward_request(self, message: Dict[str, Any], tunnel: Tunnel) -> None:
        """Handle a forward request from the server"""
        try:
            # Extract request data
            request_data = message.get("data", {})
            
            # Connect to local service
            reader, writer = await asyncio.open_connection(
                tunnel.local_host, tunnel.local_port
            )
            
            try:
                # Forward the request based on protocol
                if tunnel.protocol == "http":
                    await self._forward_http_request(message, reader, writer)
                else:
                    await self._forward_tcp_request(message, reader, writer)
            finally:
                writer.close()
                await writer.wait_closed()
                
        except Exception as e:
            logger.error(f"Failed to handle forward request: {e}")
            await self._send_error_response(message, str(e))
    
    async def _forward_http_request(self, message: Dict[str, Any], reader: asyncio.StreamReader, writer: asyncio.StreamWriter) -> None:
        """Forward an HTTP request"""
        # Implementation for HTTP forwarding
        # This would parse the HTTP request and forward it to the local service
        logger.debug("Forwarding HTTP request")
        
        # For now, just log the action
        request_data = message.get("data", {})
        logger.info(f"HTTP request: {request_data}")
    
    async def _forward_tcp_request(self, message: Dict[str, Any], reader: asyncio.StreamReader, writer: asyncio.StreamWriter) -> None:
        """Forward a TCP request"""
        # Implementation for TCP forwarding
        # This would forward raw TCP data
        logger.debug("Forwarding TCP request")
        
        # For now, just log the action
        request_data = message.get("data", {})
        logger.info(f"TCP request: {request_data}")
    
    async def _send_error_response(self, original_message: Dict[str, Any], error_msg: str) -> None:
        """Send an error response to the server"""
        response = {
            "action": "error_response",
            "request_id": original_message.get("request_id"),
            "error": error_msg
        }
        
        try:
            await self.websocket.send(json.dumps(response))
        except Exception as e:
            logger.error(f"Failed to send error response: {e}")
    
    async def close(self) -> None:
        """Close the tunnel connection"""
        if self.websocket:
            await self.websocket.close()
    
    async def list_tunnels(self) -> List[Tunnel]:
        """List all tunnels for the authenticated user"""
        # Implementation for listing tunnels
        # This would make an HTTP request to the API
        return []
    
    async def delete_tunnel(self, tunnel_id: str) -> None:
        """Delete a tunnel"""
        # Implementation for deleting a tunnel
        # This would make an HTTP request to the API
        pass
    
    async def get_tunnel_stats(self, tunnel_id: str) -> Dict[str, Any]:
        """Get statistics for a tunnel"""
        # Implementation for getting tunnel statistics
        # This would make an HTTP request to the API
        return {}


# Convenience functions for common use cases

def http_tunnel(server_url: str, token: str, subdomain: str, local_port: int) -> GoTunnelClient:
    """Create a new HTTP tunnel client"""
    config = TunnelConfig(
        server_url=server_url,
        token=token,
        subdomain=subdomain,
        local_port=local_port,
        protocol="http"
    )
    return GoTunnelClient(config)


def tcp_tunnel(server_url: str, token: str, subdomain: str, local_port: int) -> GoTunnelClient:
    """Create a new TCP tunnel client"""
    config = TunnelConfig(
        server_url=server_url,
        token=token,
        subdomain=subdomain,
        local_port=local_port,
        protocol="tcp"
    )
    return GoTunnelClient(config)


async def quick_start(server_url: str, token: str, subdomain: str, local_port: int) -> Tunnel:
    """Quick start a tunnel"""
    client = http_tunnel(server_url, token, subdomain, local_port)
    
    await client.connect()
    
    try:
        tunnel = await client.create_tunnel()
        await client.start_tunnel(tunnel)
        return tunnel
    finally:
        await client.close()


# Example usage
async def main():
    """Example usage of the GoTunnel Python SDK"""
    
    # Configuration
    server_url = "https://tunnel.gotunnel.com"
    token = "your-auth-token"
    subdomain = "myapp"
    local_port = 3000
    
    # Create client
    client = http_tunnel(server_url, token, subdomain, local_port)
    
    try:
        # Connect to server
        await client.connect()
        
        # Create tunnel
        tunnel = await client.create_tunnel()
        print(f"Tunnel created: {tunnel.public_url}")
        
        # Start tunnel
        await client.start_tunnel(tunnel)
        
    except Exception as e:
        logger.error(f"Error: {e}")
    finally:
        await client.close()


if __name__ == "__main__":
    asyncio.run(main()) 