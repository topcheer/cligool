# CliGool - User Guide

A Go and WebSocket-based cross-platform remote terminal solution supporting 18 operating systems and architectures.

## 🏗️ System Architecture

```
CLI Client          Relay Server          Web Browser
┌──────────────┐              ┌──────────┐              ┌─────────────┐
│ Real PTY     │──WebSocket──▶│ Message  │◀──WebSocket───│  xterm.js    │
│              │              │ Forwarder│              │  Terminal    │
│ 18 Platforms │              │          │              │              │
└──────────────┘              └──────────┘              └─────────────┘
```

**Key Features**:
- ✅ **Real PTY**: CLI client provides complete terminal environment
- ✅ **Message Forwarding**: Relay server routes WebSocket messages
- ✅ **Independent Interface**: Web interface based on xterm.js
- ✅ **Session Management**: Create, delete, and list sessions

## 🚀 Quick Start

### Method 1: Docker Compose (Recommended)

#### 1. Start Relay Server

```bash
# Clone repository
git clone https://github.com/topcheer/cligool.git
cd cligool

# Start all services
docker-compose up -d

# Check service status
docker-compose ps
```

**Expected Output**:
```
NAME                STATUS
cligool-relay       Up (healthy)
cligool-postgres    Up
cligool-redis       Up
```

#### 2. Start CLI Client

**On the machine you want to control remotely**:

```bash
# Download client for your platform
# Get from download page: http://localhost:8081/

# Start client (connect to local server)
./cligool -server http://localhost:8081

# Or connect to remote server
./cligool -server https://your-server.com
```

**Output Example**:
```
╔═══════════════════════════════════════════════════════════╗
║              🚀 CliGool Remote Terminal                    ║
╠═══════════════════════════════════════════════════════════╣
║ 📋 Session ID: abc123-def456-7890-abcd-ef1234567890      ║
║ 🌐 Web Access: http://localhost:8081/session/abc123-...    ║
║ 🔗 Status: 🟢 Connected                                   ║
╚═══════════════════════════════════════════════════════════╝
```

#### 3. Open Web Interface

Visit in your browser:
```
http://localhost:8081/session/[session-id]
```

Or visit the download page:
```
http://localhost:8081/
```

### Method 2: Manual Build

#### 1. Build Relay Server

```bash
# Build server
go build -o bin/relay-server ./cmd/relay

# Run server
./bin/relay-server
```

#### 2. Build Client

```bash
# Build for current platform
go build -o cligool ./cmd/client

# Cross-compile for other platforms
GOOS=linux GOARCH=amd64 go build -o cligool-linux-amd64 ./cmd/client
GOOS=windows GOARCH=amd64 go build -o cligool-windows-amd64.exe ./cmd/client
GOOS=darwin GOARCH=arm64 go build -o cligool-darwin-arm64 ./cmd/client
```

## 🌍 Supported Platforms

### Windows (2 platforms)
- Windows amd64 (Intel/AMD 64-bit)
- Windows arm64 (Surface Pro X and other ARM devices)

### Linux (8 platforms)
- Linux amd64 (Ubuntu, Debian, CentOS, etc.)
- Linux arm64 (Raspberry Pi 4/5, ARM servers)
- Linux 386 (32-bit x86 systems)
- Linux arm (Raspberry Pi and other 32-bit ARM devices)
- Linux ppc64le (PowerPC systems)
- Linux riscv64 (RISC-V architecture)
- Linux s390x (IBM System z mainframes)
- Linux mips64le (MIPS architecture)

### *BSD Systems (6 platforms)
- FreeBSD amd64/arm64
- OpenBSD amd64/arm64
- NetBSD amd64
- DragonFlyBSD amd64

### macOS (2 platforms)
- macOS Intel (Intel processors)
- macOS ARM (Apple M1/M2/M3)

## 💡 Use Cases

### 1. Remote Home Access

```bash
# Start client on home Mac
./cligool-darwin-arm64 -server https://your-server.com

# Connect from office browser
# Use generated session ID
```

### 2. Server Management

```bash
# Run on Linux server
./cligool-linux-amd64 -server https://your-server.com

# Manage from mobile browser
```

### 3. Tech Support

```bash
# Friend's computer has issues
# Have friend download and run client for their platform
# Remotely assist from your browser
```

### 4. Team Collaboration

```bash
# Multiple people connect to same session
# Real-time viewing and control
```

## 🔧 Advanced Features

### Heartbeat Keep-Alive

Bidirectional WebSocket heartbeat mechanism:
- **Server → Client**: Send ping every 30 seconds
- **Client → Server**: Auto-reply with pong
- **Timeout Detection**: Auto-disconnect after 90 seconds

### Auto-Reconnect

Web interface supports auto-reconnection:
- Auto-retry on disconnect
- Exponential backoff strategy
- Manual reconnect button

### Session Management

```bash
# Create new session
curl -X POST http://localhost:8081/api/sessions \
  -H "Content-Type: application/json" \
  -d '{"owner": "user@example.com"}'

# List all sessions
curl http://localhost:8081/api/sessions

# Get session details
curl http://localhost:8081/api/sessions/{session_id}

# Delete session
curl -X DELETE http://localhost:8081/api/sessions/{session_id}
```

## 🔍 Troubleshooting

### Client Connection Failed

**Problem**: `WebSocket connection failed`

**Solution**:
```bash
# 1. Check if server is running
curl http://localhost:8081/api/health

# 2. Check server logs
docker logs cligool-relay

# 3. Check firewall settings
# Ensure port 8081 is accessible
```

### Terminal Not Responding

**Problem**: No output after entering commands

**Solution**:
1. Check WebSocket connection status (browser console)
2. Confirm client process is still running
3. Try refreshing the page to reconnect

### Windows Client Garbled Text

**Problem**: Chinese characters display as gibberish

**Solution**:
- ✅ Auto-fixed: GBK encoding automatically converted to UTF-8
- If issues persist, check terminal encoding settings

### PTY Permission Issues (Linux/macOS)

**Problem**: `Failed to allocate PTY`

**Solution**:
```bash
# Check /dev/ptmx permissions
ls -l /dev/ptmx

# Ensure running in real terminal (not IDE built-in terminal)
# Some IDE terminals may not support PTY
```

## 🚀 Production Deployment

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f relay-server

# Stop services
docker-compose down
```

### Using Cloudflare Tunnel (HTTPS)

1. Install cloudflared
2. Create tunnel: `cloudflared tunnel create cligool`
3. Configure `~/.cloudflared/config.yml`:
```yaml
tunnel: <your-tunnel-id>
credentials-file: /path/to/credentials.json

ingress:
  - hostname: cligool.yourdomain.com
    service: http://localhost:8081
  - service: http_status:404
```

4. Start tunnel: `cloudflared tunnel run`

### Environment Variables

```bash
# Database connection
DATABASE_URL=postgres://user:pass@host:5432/cligool?sslmode=disable

# Redis connection
REDIS_URL=redis://host:6379

# Server configuration
RELAY_HOST=0.0.0.0
RELAY_PORT=8080
```

## 📚 More Documentation

- [README.md](../README.md) - Project Overview
- [DEPLOYMENT.md](DEPLOYMENT.md) - Deployment Guide
- [DEVELOPMENT.md](DEVELOPMENT.md) - Development Guide
- [PTY_TROUBLESHOOTING.md](PTY_TROUBLESHOOTING.md) - PTY Issues

## 🤝 Contributing

Issues and Pull Requests are welcome!

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📝 License

MIT License - see [LICENSE](../LICENSE) file for details
