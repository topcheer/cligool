# CliGool Docker Deployment Guide

This document describes how to deploy CliGool relay server using Docker.

## 🐳 Image Information

### Image Registry

- **GitHub Container Registry**: `ghcr.io/topcheer/cligool`
- **Supported Architectures**: linux/amd64, linux/arm64
- **Contents**: Relay server + client binaries for all 33 platforms

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

# Run container (requires PostgreSQL and Redis first)
docker run -d \
  --name cligool-relay \
  -p 8081:8080 \
  -e DATABASE_URL=postgres://cligool:cligool123@postgres:5432/cligool?sslmode=disable \
  -e REDIS_URL=redis://redis:6379 \
  ghcr.io/topcheer/cligool:latest
```

### Method 3: Complete Environment (with databases)

```bash
# Create network
docker network create cligool-network

# Start PostgreSQL
docker run -d \
  --name cligool-postgres \
  --network cligool-network \
  -e POSTGRES_DB=cligool \
  -e POSTGRES_USER=cligool \
  -e POSTGRES_PASSWORD=cligool123 \
  postgres:15-alpine

# Start Redis
docker run -d \
  --name cligool-redis \
  --network cligool-network \
  redis:7-alpine

# Start CliGool relay server
docker run -d \
  --name cligool-relay \
  --network cligool-network \
  -p 8081:8080 \
  -e DATABASE_URL=postgres://cligool:cligool123@cligool-postgres:5432/cligool?sslmode=disable \
  -e REDIS_URL=redis://cligool-redis:6379 \
  ghcr.io/topcheer/cligool:latest
```

## ⚙️ Environment Variables

### Database Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | - |
| `REDIS_URL` | Redis connection string | - |

### Server Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `RELAY_HOST` | Listening address | 0.0.0.0 |
| `RELAY_PORT` | Listening port | 8080 |
| `ENABLE_AUTO_HTTPS` | Auto-enable HTTPS | false |

### Example

```bash
docker run -d \
  --name cligool-relay \
  -p 8081:8080 \
  -e DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=disable \
  -e REDIS_URL=redis://host:6379 \
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

1. **Don't use default passwords in production**
2. **Use HTTPS**: Configure reverse proxy and enable SSL
3. **Restrict network access**: Use firewall to restrict database access
4. **Update images regularly**: `docker pull ghcr.io/topcheer/cligool:latest`
5. **Backup PostgreSQL data**: Regularly backup PostgreSQL data volumes

## 🐛 Troubleshooting

### Container won't start

```bash
# Check logs
docker logs cligool-relay

# Check database connection
docker exec cligool-relay ping -c 3 cligool-postgres
```

### Cannot access web interface

```bash
# Check port mapping
docker ps | grep cligool-relay

# Check firewall
sudo ufw status
```

### Database connection failed

```bash
# Check database container
docker ps | grep postgres

# Test database connection
docker exec cligool-relay sh -c 'apk add postgresql-client && psql $DATABASE_URL'
```

## 📚 More Information

- [Project Homepage](https://github.com/topcheer/cligool)
- [Quick Start Guide](https://github.com/topcheer/cligool/blob/main/QUICKSTART.md)
- [Usage Guide](https://github.com/topcheer/cligool/blob/main/USAGE_GUIDE.md)
- [Configuration Guide](https://github.com/topcheer/cligool/blob/main/CONFIG_GUIDE.md)

## 🌏 Language

- [English](DOCKER_EN.md) (This file)
- [中文](DOCKER.md)
