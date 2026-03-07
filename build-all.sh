#!/bin/bash
# 构建所有平台的客户端（33个平台）

set -e

echo "🔨 开始构建所有平台的客户端..."
echo "📊 目标平台：33个操作系统和架构组合"
echo ""

# 创建输出目录
mkdir -p bin

# ==================== macOS 版本 ====================
echo "📦 [1/33] 构建 macOS amd64..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/cligool-darwin-amd64 ./cmd/client

echo "📦 [2/33] 构建 macOS arm64..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o bin/cligool-darwin-arm64 ./cmd/client

# ==================== Linux 版本 ====================
echo "📦 [3/33] 构建 Linux amd64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/cligool-linux-amd64 ./cmd/client

echo "📦 [4/33] 构建 Linux arm64..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/cligool-linux-arm64 ./cmd/client

echo "📦 [5/33] 构建 Linux 386..."
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o bin/cligool-linux-386 ./cmd/client

echo "📦 [6/33] 构建 Linux arm..."
GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -o bin/cligool-linux-arm ./cmd/client

echo "📦 [7/33] 构建 Linux armbe (ARM64 Big-Endian)..."
GOOS=linux GOARCH=arm64 GOARM=7 CGO_ENABLED=0 go build -o bin/cligool-linux-armbe ./cmd/client

echo "📦 [8/33] 构建 Linux ppc64le..."
GOOS=linux GOARCH=ppc64le CGO_ENABLED=0 go build -o bin/cligool-linux-ppc64le ./cmd/client

echo "📦 [9/33] 构建 Linux ppc64 (Big-Endian)..."
GOOS=linux GOARCH=ppc64 GOBIGENDIAN=true CGO_ENABLED=0 go build -o bin/cligool-linux-ppc64 ./cmd/client

echo "📦 [10/33] 构建 Linux riscv64..."
GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 go build -o bin/cligool-linux-riscv64 ./cmd/client

echo "📦 [11/33] 构建 Linux s390x..."
GOOS=linux GOARCH=s390x CGO_ENABLED=0 go build -o bin/cligool-linux-s390x ./cmd/client

echo "📦 [12/33] 构建 Linux mips..."
GOOS=linux GOARCH=mips CGO_ENABLED=0 go build -o bin/cligool-linux-mips ./cmd/client

echo "📦 [13/33] 构建 Linux mips64le..."
GOOS=linux GOARCH=mips64le CGO_ENABLED=0 go build -o bin/cligool-linux-mips64le ./cmd/client

echo "📦 [14/33] 构建 Linux mips64..."
GOOS=linux GOARCH=mips64 GOMIPS=hardfloat CGO_ENABLED=0 go build -o bin/cligool-linux-mips64 ./cmd/client

echo "📦 [15/33] 构建 Linux loong64..."
GOOS=linux GOARCH=loong64 CGO_ENABLED=0 go build -o bin/cligool-linux-loong64 ./cmd/client

# ==================== FreeBSD 版本 ====================
echo "📦 [16/33] 构建 FreeBSD amd64..."
GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build -o bin/cligool-freebsd-amd64 ./cmd/client

echo "📦 [17/33] 构建 FreeBSD arm64..."
GOOS=freebsd GOARCH=arm64 CGO_ENABLED=0 go build -o bin/cligool-freebsd-arm64 ./cmd/client

echo "📦 [18/33] 构建 FreeBSD 386..."
GOOS=freebsd GOARCH=386 CGO_ENABLED=0 go build -o bin/cligool-freebsd-386 ./cmd/client

echo "📦 [19/33] 构建 FreeBSD arm..."
GOOS=freebsd GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -o bin/cligool-freebsd-arm ./cmd/client

echo "📦 [20/33] 构建 FreeBSD riscv64..."
GOOS=freebsd GOARCH=riscv64 CGO_ENABLED=0 go build -o bin/cligool-freebsd-riscv64 ./cmd/client

# ==================== OpenBSD 版本 ====================
echo "📦 [21/33] 构建 OpenBSD amd64..."
GOOS=openbsd GOARCH=amd64 CGO_ENABLED=0 go build -o bin/cligool-openbsd-amd64 ./cmd/client

echo "📦 [22/33] 构建 OpenBSD arm64..."
GOOS=openbsd GOARCH=arm64 CGO_ENABLED=0 go build -o bin/cligool-openbsd-arm64 ./cmd/client

echo "📦 [23/33] 构建 OpenBSD 386..."
GOOS=openbsd GOARCH=386 CGO_ENABLED=0 go build -o bin/cligool-openbsd-386 ./cmd/client

echo "📦 [24/33] 构建 OpenBSD arm..."
GOOS=openbsd GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -o bin/cligool-openbsd-arm ./cmd/client

echo "📦 [25/33] 构建 OpenBSD riscv64..."
GOOS=openbsd GOARCH=riscv64 CGO_ENABLED=0 go build -o bin/cligool-openbsd-riscv64 ./cmd/client

# ==================== NetBSD 版本 ====================
echo "📦 [26/33] 构建 NetBSD amd64..."
GOOS=netbsd GOARCH=amd64 CGO_ENABLED=0 go build -o bin/cligool-netbsd-amd64 ./cmd/client

echo "📦 [27/33] 构建 NetBSD arm64..."
GOOS=netbsd GOARCH=arm64 CGO_ENABLED=0 go build -o bin/cligool-netbsd-arm64 ./cmd/client

echo "📦 [28/33] 构建 NetBSD arm..."
GOOS=netbsd GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -o bin/cligool-netbsd-arm ./cmd/client

echo "📦 [29/33] 构建 NetBSD 386..."
GOOS=netbsd GOARCH=386 CGO_ENABLED=0 go build -o bin/cligool-netbsd-386 ./cmd/client

# ==================== DragonFlyBSD 版本 ====================
echo "📦 [30/33] 构建 DragonFlyBSD amd64..."
GOOS=dragonfly GOARCH=amd64 CGO_ENABLED=0 go build -o bin/cligool-dragonfly-amd64 ./cmd/client

echo "📦 [31/33] 构建 DragonFlyBSD arm64..."
GOOS=dragonfly GOARCH=arm64 CGO_ENABLED=0 go build -o bin/cligool-dragonfly-arm64 ./cmd/client

# ==================== Relay Server ====================
echo "📦 [32/33] 构建中继服务器 (Linux amd64)..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/relay-server ./cmd/relay

echo ""
echo "✅ 构建完成！"
echo ""
echo "📊 构建统计："
echo "   - macOS: 2个平台"
echo "   - Linux: 13个平台"
echo "   - FreeBSD: 5个平台"
echo "   - OpenBSD: 5个平台"
echo "   - NetBSD: 4个平台"
echo "   - DragonFlyBSD: 2个平台"
echo "   - 中继服务器: 1个"
echo "   总计: 32个二进制文件"
echo ""
echo "📦 构建的文件："
ls -lh bin/
echo ""
echo "💡 提示："
echo "   - Windows版本请在Docker容器中构建或使用GitHub Actions自动构建"
echo "   - 所有二进制文件都是静态编译，无需额外依赖"
echo "   - 可以直接复制到目标系统运行"
