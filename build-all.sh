#!/bin/bash
# 构建所有平台的客户端

set -e

echo "🔨 开始构建所有平台的客户端..."

# macOS Intel
echo "📦 构建 macOS amd64..."
GOOS=darwin GOARCH=amd64 go build -o bin/cligool-darwin-amd64 ./cmd/client

# macOS ARM64 (当前平台)
echo "📦 构建 macOS arm64..."
GOOS=darwin GOARCH=arm64 go build -o bin/cligool-darwin-arm64 ./cmd/client

# Linux amd64
echo "📦 构建 Linux amd64..."
GOOS=linux GOARCH=amd64 go build -o bin/cligool-linux-amd64 ./cmd/client

# Linux arm64
echo "📦 构建 Linux arm64..."
GOOS=linux GOARCH=arm64 go build -o bin/cligool-linux-arm64 ./cmd/client

# Linux 386
echo "📦 构建 Linux 386..."
GOOS=linux GOARCH=386 go build -o bin/cligool-linux-386 ./cmd/client

# Linux arm (32位)
echo "📦 构建 Linux arm..."
GOOS=linux GOARCH=arm go build -o bin/cligool-linux-arm ./cmd/client

# Linux ppc64le
echo "📦 构建 Linux ppc64le..."
GOOS=linux GOARCH=ppc64le go build -o bin/cligool-linux-ppc64le ./cmd/client

# Linux riscv64
echo "📦 构建 Linux riscv64..."
GOOS=linux GOARCH=riscv64 go build -o bin/cligool-linux-riscv64 ./cmd/client

# Linux s390x
echo "📦 构建 Linux s390x..."
GOOS=linux GOARCH=s390x go build -o bin/cligool-linux-s390x ./cmd/client

# Linux mips64le
echo "📦 构建 Linux mips64le..."
GOOS=linux GOARCH=mips64le go build -o bin/cligool-linux-mips64le ./cmd/client

# FreeBSD amd64
echo "📦 构建 FreeBSD amd64..."
GOOS=freebsd GOARCH=amd64 go build -o bin/cligool-freebsd-amd64 ./cmd/client

# FreeBSD arm64
echo "📦 构建 FreeBSD arm64..."
GOOS=freebsd GOARCH=arm64 go build -o bin/cligool-freebsd-arm64 ./cmd/client

# OpenBSD amd64
echo "📦 构建 OpenBSD amd64..."
GOOS=openbsd GOARCH=amd64 go build -o bin/cligool-openbsd-amd64 ./cmd/client

# OpenBSD arm64
echo "📦 构建 OpenBSD arm64..."
GOOS=openbsd GOARCH=arm64 go build -o bin/cligool-openbsd-arm64 ./cmd/client

# NetBSD amd64
echo "📦 构建 NetBSD amd64..."
GOOS=netbsd GOARCH=amd64 go build -o bin/cligool-netbsd-amd64 ./cmd/client

# DragonFlyBSD amd64
echo "📦 构建 DragonFlyBSD amd64..."
GOOS=dragonfly GOARCH=amd64 go build -o bin/cligool-dragonfly-amd64 ./cmd/client

echo ""
echo "✅ 构建完成！"
echo ""
echo "📊 构建结果："
ls -lh bin/

echo ""
echo "💡 提示：Windows版本请在Docker容器中构建或使用Docker镜像"

