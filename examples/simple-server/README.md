# Simple Example Server

This is a simple HTTP server that demonstrates how to use GoTunnel.

## Running the Example

1. **Start the example server:**
   ```bash
   go run main.go
   ```
   This will start a server on port 3000.

2. **Start the GoTunnel server:**
   ```bash
   ./gotunnel-server --port 8080 --cert certs/cert.pem --key certs/key.pem --allowed-tokens "dev-token" --log-level info
   ```

3. **Start the GoTunnel client:**
   ```bash
   ./gotunnel-client --server localhost:8080 --subdomain myapp --local-port 3000 --token "dev-token" --skip-verify --log-level info
   ```

4. **Access your tunneled service:**
   Visit `http://myapp.localhost:8080` to see your local server exposed through the tunnel.

## Features

- Shows request information (method, path, headers)
- Displays connection details
- Provides instructions for using GoTunnel
- Responsive design

## Customization

You can change the port by passing it as an argument:

```bash
go run main.go 8080
```

This will start the server on port 8080 instead of 3000. 