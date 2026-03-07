# CliGool - Cross-Platform Remote Terminal Solution

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platforms](https://img.shields.io/badge/platforms-30-blue)](#-supported-platforms)
[![Docker](https://img.shields.io/badge/docker-multi--arch-blue.svg)](https://ghcr.io/topcheer/cligool)
[![Demo](https://img.shields.io/badge/demo-online-success.svg)](https://cligool.zty8.cn/)

A cross-platform remote terminal solution based on Go and WebSocket, supporting 30 operating systems and architectures.

**[🚀 Live Demo](https://cligool.zty8.cn/)** | **[📥 Download Clients](https://cligool.zty8.cn/)**

## ✨ Key Features

- 🌍 **Cross-Platform Support**: 30 operating systems and architectures (Windows, Linux, macOS, *BSD, etc.)
- ⚡ **Low Latency**: Real-time WebSocket communication with millisecond response
- 🔒 **Secure Connection**: End-to-end encrypted communication, supporting HTTPS/WSS
- 💎 **Real PTY**: Full terminal feature support (colors, cursor control, etc.)
- 🤖 **AI CLI Tools**: Perfect support for Claude, Gemini, Aider, and other AI CLI tools
- 👥 **Multi-User Collaboration**: Multiple users can connect to the same terminal session
- 🚀 **Ready to Use**: One-click Docker deployment, no complex configuration needed
- 🎨 **Modern Web Interface**: Professional terminal UI based on xterm.js
- 💓 **Heartbeat Keep-Alive**: Bidirectional heartbeat mechanism with auto-reconnect and dead connection cleanup
- 🖥️ **Windows ConPTY**: Windows version uses ConPTY with feature parity to Unix systems

## 🏗️ System Architecture

```
CLI Client            Relay Server          Web Browser
┌──────────────┐                    ┌──────────┐                  ┌─────────────┐
│ Real PTY     │──WebSocket──▶     │ Message  │◀──WebSocket───   │  xterm.js   │
│ Environment  │                    │ Forwarder│                  │  Terminal   │
│              │                    │          │                  │  Interface  │
│ 30 Platforms │                    │          │                  │             │
└──────────────┘                    └──────────┘                  └─────────────┘
```

**Core Features**:
- ✅ **Real PTY**: CLI client provides a complete terminal environment
- ✅ **Message Forwarding**: Relay server handles WebSocket message routing
- ✅ **Independent Interface**: Web interface based on xterm.js, works with any browser
- ✅ **Stateless Design**: In-memory session management, no database required
- ✅ **Ready to Use**: Single container deployment, minimal configuration needed

## 🌍 Supported Platforms

### Windows (2 platforms)
- Windows amd64 (Intel/AMD 64-bit)
- Windows arm64 (Surface Pro X and other ARM devices)

### Linux (13 platforms)
- Linux amd64 (Ubuntu, Debian, CentOS and other 64-bit systems)
- Linux arm64 (Raspberry Pi 4/5, ARM servers)
- Linux 386 (32-bit x86 systems)
- Linux arm (Raspberry Pi and other 32-bit ARM devices)
- Linux armbe (ARM64 Big-Endian)
- Linux ppc64le (PowerPC Little-Endian)
- Linux ppc64 (PowerPC Big-Endian)
- Linux riscv64 (RISC-V architecture)
- Linux s390x (IBM System z mainframes)
- Linux mips (MIPS 32-bit)
- Linux mips64le (MIPS 64-bit Little-Endian)
- Linux mips64 (MIPS 64-bit Big-Endian)
- Linux loong64 (LoongArch)

### *BSD Systems (12 platforms)
- FreeBSD amd64/arm64/386/arm/riscv64
- OpenBSD amd64/arm64
- NetBSD amd64/arm64/arm/386
- DragonFlyBSD amd64

### macOS (2 platforms)
- macOS Intel (Intel processors)
- macOS ARM (Apple M1/M2/M3/M4)

## 🎮 Live Demo

Don't want to deploy locally? Try the online demo first!

### Access Online Demo

👉 **[https://cligool.zty8.cn/](https://cligool.zty8.cn/)**

### Demo Usage Steps

1. **Download Client**: Select the client for your OS from the download page
2. **Run Client**: Launch the downloaded client program
3. **Access Web Terminal**: Use the session ID displayed by the client to access the web interface
4. **Start Experience**: Experience full remote terminal functionality in your browser

**Note**: The online demo is for experience only and may be shut down at any time. It's recommended to deploy your own instance for stable service.

## 🚀 Quick Start

### Method 1: Docker Compose Deployment (Recommended)

```bash
# 1. Clone repository
git clone https://github.com/topcheer/cligool.git
cd cligool

# 2. Production environment (using pre-built image)
docker-compose up -d

# Or: Development environment (build locally)
docker-compose -f docker-compose.dev.yml up -d --build

# 3. Check service status
docker-compose ps
```

Services started:
- **CliGool Relay Server**: http://localhost:8081

Access web interface: **http://localhost:8081**

### Method 2: Using Docker Image

```bash
# Pull image
docker pull ghcr.io/topcheer/cligool:latest

# Run relay server
docker run -d \
  --name cligool-relay \
  -p 8081:8080 \
  ghcr.io/topcheer/cligool:latest
```

### Method 3: Build from Source

```bash
# 1. Clone repository
git clone https://github.com/topcheer/cligool.git
cd cligool

# 2. Build relay server
go build -o relay-server ./cmd/relay

# 3. Run relay server
export RELAY_HOST=0.0.0.0
export RELAY_PORT=8080
./relay-server
```

## 📥 Client Installation

### Automatic Installation (macOS/Linux)

```bash
curl -sSL https://raw.githubusercontent.com/topcheer/cligool/main/install.sh | bash
```

### Manual Installation

1. **Download**: Visit [GitHub Releases](https://github.com/topcheer/cligool/releases)
2. **Select Platform**: Choose the appropriate file for your OS
   - macOS: `cligool-darwin-amd64` or `cligool-darwin-arm64`
   - Linux: `cligool-linux-amd64`, `cligool-linux-arm64`, etc.
   - Windows: `cligool-windows-amd64.zip` or `cligool-windows-arm64.zip`
3. **Install**:
   - **macOS/Linux**:
     ```bash
     chmod +x cligool-*
     sudo mv cligool-* /usr/local/bin/cligool
     ```
   - **Windows**:
     - Unzip the file
     - Add to PATH environment variable
     - Run in terminal: `cligool.exe`

## 💻 Usage

### Start CLI Client

```bash
# Basic usage (connect to local server)
cligool

# Connect to remote server
cligool -server https://your-server.com

# Specify terminal size
cligool -server http://localhost:8081 -cols 120 -rows 36

# Use configuration file
cligool -config ~/.cligool.json
```

### Configuration File

Create `~/.cligool.json` or `./cligool.json`:

```json
{
  "server": "http://localhost:8081",
  "cols": 120,
  "rows": 36
}
```

Configuration priority: `./cligool.json` → `~/.cligool.json` → auto-create `~/.cligool.json` with defaults

### Access Web Terminal

1. Start CLI client and get session ID
2. Open browser and visit: `http://localhost:8081/terminal/{session_id}`
3. Start using!

## 🎯 Use Cases

- **Remote Development**: Access your development terminal from anywhere
- **Pair Programming**: Real-time collaborative coding
- **Technical Support**: Remote troubleshooting and debugging
- **Teaching & Demo**: Share terminal operations with an audience
- **Mobile Development**: Use terminal on mobile devices
- **AI CLI Tools**: Use Claude, Gemini, Aider and other AI tools in web interface

## 🔧 Advanced Configuration

### Environment Variables

- `RELAY_HOST`: Server listening address (default: 0.0.0.0)
- `RELAY_PORT`: Server listening port (default: 8080)

### Command Line Options

```
-server string        Server URL (default: http://localhost:8081)
-config string        Configuration file path
-cols int             Terminal columns (default: auto-detect)
-rows int             Terminal rows (default: auto-detect)
-debug                Enable debug mode
-version              Show version information
```

## 🐳 Docker Deployment

See [Docker Deployment Guide](docs/DOCKER.md) for detailed instructions.

### Quick Start with Docker

```bash
# Pull image
docker pull ghcr.io/topcheer/cligool:latest

# Run with docker-compose
docker-compose up -d

# Access web interface
open http://localhost:8081
```

### Supported Architectures

- linux/amd64 - Intel/AMD 64-bit
- linux/arm64 - ARM 64-bit (AWS Graviton, Apple Silicon)

## 📚 Documentation

- [Quick Start Guide](QUICKSTART.md)
- [Usage Guide](USAGE_GUIDE.md) | [中文版](USAGE_GUIDE_CN.md)
- [Configuration Guide](CONFIG_GUIDE.md)
- [Platform List](PLATFORMS.md)
- [Docker Deployment](docs/DOCKER.md)
- [Command Line Arguments](CMD_ARGS_USAGE.md)

## 🛠️ Development

### Prerequisites

- Go 1.24+
- Docker & Docker Compose (for containerized deployment)

### Build Relay Server

```bash
go build -o relay-server ./cmd/relay
```

### Build Client for All Platforms

```bash
./build-all.sh
```

### Run Tests

```bash
go test ./...
```

### Start Development Environment

```bash
make dev
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

## 🙏 Acknowledgments

- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket library
- [creack/pty](https://github.com/creack/pty) - PTY support for Unix-like systems
- [xterm.js](https://xtermjs.org/) - Terminal emulator in the browser

## 📞 Contact

- GitHub: https://github.com/topcheer/cligool
- Issues: https://github.com/topcheer/cligool/issues
- Live Demo: https://cligool.zty8.cn

---

**If you find CliGool helpful, please give us a ⭐️!**
