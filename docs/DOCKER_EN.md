# CliGool Docker Deployment Guide

This document describes how to deploy CliGool relay server using Docker.

## 🐳 Image Information

### Image Registry

- **GitHub Container Registry**: `ghcr.io/topcheer/cligool`
- **Supported Architectures**: linux/amd64, linux/arm64
- **Contents**: Relay server + client binaries for all 30 platforms

### Available Tags

- `latest` - Latest stable version
- `v1.1.0`, `v1.2.0`, etc. - Specific versions
- `v1` - Latest major version
- `v1.1` - Latest minor version

## 🚀 Quick Start

### Method 1: Docker Compose (Recommended)

```bash
# Clone repository
git clone https://github.com/topcheer/cligool.git
cd cligool

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f relay-server

# Access web interface
open http://localhost:8081
```

### Method 2: Docker Commands

```bash
# Pull image
docker pull ghcr.io/topcheer/cligool:latest

# Run container
docker run -d \
  --name cligool-relay \
  -p 8081:8080 \
  ghcr.io/topcheer/cligool:latest
```

## ⚙️ Environment Variables

### Server Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `RELAY_HOST` | Listening address | 0.0.0.0 |
| `RELAY_PORT` | Listening port | 8080 |

### Example

```bash
docker run -d \
  --name cligool-relay \
  -p 8081:8080 \
  -e RELAY_HOST=0.0.0.0 \
  -e RELAY_PORT=8080 \
  ghcr.io/topcheer/cligool:latest
```

## 📦 Multi-Architecture Support

The Docker image supports the following architectures:

- **linux/amd64** - x86_64 servers (Intel/AMD)
- **linux/arm64** - ARM 64-bit servers (AWS Graviton, Apple Silicon, etc.)

Docker will automatically pull the appropriate image for your system architecture:

```bash
# Automatically select correct architecture
docker pull ghcr.io/topcheer/cligool:latest

# View image architecture
docker inspect ghcr.io/topcheer/cligool:latest | grep Architecture
```

## 🔧 Building Custom Images

### Local Build (Single Architecture)

```bash
# Use standard Dockerfile
docker build -t cligool:local .

# Or use multi-architecture Dockerfile
docker build -f Dockerfile.multiarch -t cligool:local .
```

### Local Build (Multi-Architecture)

```bash
# Enable buildx
docker buildx create --use

# Build multi-architecture image
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t cligool:multiarch \
  -f Dockerfile.multiarch \
  --push \
  .
```

## 📊 Container Management

### View Logs

```bash
# View all logs
docker logs cligool-relay

# Real-time log viewing
docker logs -f cligool-relay

# View last 100 lines
docker logs --tail 100 cligool-relay
```

### Enter Container

```bash
# Use shell to enter container
docker exec -it cligool-relay sh

# Check server status
docker exec cligool-relay wget -qO- http://localhost:8080/api/health
```

### Restart Container

```bash
# Restart container
docker restart cligool-relay

# Stop container
docker stop cligool-relay

# Start container
docker start cligool-relay
```

### Cleanup

```bash
# Stop and remove container
docker stop cligool-relay && docker rm cligool-relay

# Remove image
docker rmi ghcr.io/topcheer/cligool:latest

# Clean all related resources (including data volumes)
docker-compose down -v
```

## 🌐 Reverse Proxy Configuration

### Nginx Example

```nginx
server {
    listen 80;
    server_name cligool.example.com;

    location / {
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support
        proxy_read_timeout 86400;
    }
}
```

### Caddy Example

```
cligool.example.com {
    reverse_proxy localhost:8081
}
```

## 🔒 Security Recommendations

1. **Use HTTPS**: Configure reverse proxy and enable SSL
2. **Restrict network access**: Use firewall to restrict access
3. **Update images regularly**: `docker pull ghcr.io/topcheer/cligool:latest`
4. **Monitor logs**: Regularly check container logs

## 🐛 Troubleshooting

### Container won't start

```bash
# Check logs
docker logs cligool-relay

# Check container status
docker ps -a | grep cligool-relay
```

### Cannot access web interface

```bash
# Check port mapping
docker ps | grep cligool-relay

# Check firewall
sudo ufw status
```

### Connection Issues

```bash
# Check container status
docker ps | grep cligool-relay

# Test health endpoint
docker exec cligool-relay wget -qO- http://localhost:8080/api/health
```

## 📚 More Information

- [Project Homepage](https://github.com/topcheer/cligool)
- [Quick Start Guide](https://github.com/topcheer/cligool/blob/main/QUICKSTART.md)
- [Usage Guide](https://github.com/topcheer/cligool/blob/main/USAGE_GUIDE.md)
- [Configuration Guide](https://github.com/topcheer/cligool/blob/main/CONFIG_GUIDE.md)

## 🌏 Language

- [English](DOCKER_EN.md) (This file)
- [中文](DOCKER.md)
