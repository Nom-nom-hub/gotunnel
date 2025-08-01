# 🚀 GoTunnel

> **The open-source ngrok killer for developers who want control**

GoTunnel is a self-hosted, secure, and blazingly fast tunneling solution that gives you complete control over your tunnel infrastructure. No more limits, no more data collection, no more vendor lock-in.

## ✨ Features

- **🔒 Self-hosted** - Run your own tunnel server
- **⚡ Blazingly fast** - WebSocket-based communication
- **🔄 Unlimited tunnels** - No artificial limits
- **🖥️ GUI Dashboard** - Beautiful Electron interface
- **🔧 Easy setup** - One-click installation
- **🌐 Cross-platform** - Windows, macOS, Linux
- **📊 Real-time monitoring** - Live logs and metrics
- **💾 Profile system** - Save and load configurations

## 🚀 Quick Start

### Option 1: GUI Dashboard (Recommended)

1. **Install dependencies:**
   ```bash
   cd dashboard
   npm install
   npm start
   ```

2. **Start the server:**
   - Click "Start Server" in the GUI
   - Configure your port and token

3. **Create a tunnel:**
   - Enter your subdomain (e.g., "myapp")
   - Set local port (e.g., 3000)
   - Click "Create Tunnel"

### Option 2: Command Line

1. **Start the server:**
   ```bash
   ./gotunnel-server.exe --port 8080 --allowed-tokens "your-token" --tls=false
   ```

2. **Create a tunnel:**
   ```bash
   ./gotunnel-client.exe --server localhost:8080 --subdomain myapp --local-port 3000 --token "your-token" --tls=false
   ```

3. **Access your app:**
   - Local: `http://localhost:3000`
   - Tunnel: `http://myapp.localhost:8080`

## 📦 Installation

### Windows
```powershell
# Run as Administrator
powershell -ExecutionPolicy Bypass -File installer/install.ps1
```

### Manual Installation
1. Download the latest release
2. Extract the binaries
3. Run the server and client

## 🏗️ Project Structure

```
gotunnel/
├── cmd/                    # Go application entry points
│   ├── server/            # Tunnel server
│   └── client/            # Tunnel client
├── pkg/                   # Go packages
│   ├── tunnel/            # Core tunneling logic
│   └── auth/              # Authentication
├── dashboard/             # Electron GUI application
│   ├── main.js           # Main process
│   ├── index.html        # Renderer process
│   └── package.json      # Dependencies
├── website/              # Marketing website
├── installer/            # Installation scripts
├── examples/             # Usage examples
├── scripts/              # Build and deployment scripts
├── Dockerfile           # Docker container
├── docker-compose.yml   # Docker orchestration
├── Makefile             # Build automation
└── README.md            # This file
```

## 🔧 Configuration

Create a `config.yaml` file in your home directory:

```yaml
server:
  port: 8080
  tls: false
  allowed_tokens:
    - "your-secret-token"

client:
  server: "localhost:8080"
  token: "your-secret-token"
  tls: false
  log_level: "info"

dashboard:
  enabled: true
  port: 3000
```

## 🎯 Use Cases

- **Development** - Share local development servers
- **Testing** - Expose test environments
- **Demos** - Show your work to clients
- **Webhooks** - Receive external callbacks
- **IoT** - Access devices behind firewalls
- **Microservices** - Connect distributed services

## 🆚 vs ngrok

| Feature | GoTunnel | ngrok |
|---------|----------|-------|
| **Self-hosted** | ✅ Yes | ❌ No |
| **Unlimited tunnels** | ✅ Yes | ❌ No |
| **Free forever** | ✅ Yes | ❌ No |
| **Open source** | ✅ Yes | ❌ No |
| **WebSocket support** | ✅ Yes | ❌ No |
| **GUI dashboard** | ✅ Yes | ❌ No |
| **No data collection** | ✅ Yes | ❌ No |
| **Custom domains** | ✅ Yes | ❌ No |

## 🛠️ Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- npm

### Build from source
```bash
# Build Go binaries
make build

# Install dashboard dependencies
cd dashboard
npm install

# Start development
npm run dev
```

### Docker
```bash
# Build and run with Docker
docker-compose up -d
```

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📞 Support

- 📧 Email: support@gotunnel.dev
- 💬 Discord: [Join our community](https://discord.gg/gotunnel)
- 📖 Docs: [docs.gotunnel.dev](https://docs.gotunnel.dev)
- 🐛 Issues: [GitHub Issues](https://github.com/gotunnel/gotunnel/issues)

---

**Made with ❤️ by developers, for developers**

*The open-source ngrok killer for developers who want control.* 