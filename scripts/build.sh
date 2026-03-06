#!/bin/bash

# CliGool 构建脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔨 CliGool 构建脚本"
echo "=================="

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go 未安装，请先安装 Go 1.21+${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${GREEN}✅ Go 版本: $GO_VERSION${NC}"

# 创建构建目录
echo "📁 创建构建目录..."
mkdir -p bin

# 安装依赖
echo "📦 安装依赖..."
go mod download
go mod tidy

# 构建函数
build_binary() {
    local name=$1
    local path=$2
    local output=$3

    echo -e "${YELLOW}🔨 构建 $name...${NC}"
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o "$output" "$path"

    if [ $? -eq 0 ]; then
        size=$(du -h "$output" | cut -f1)
        echo -e "${GREEN}✅ $name 构建成功 (大小: $size)${NC}"
    else
        echo -e "${RED}❌ $name 构建失败${NC}"
        exit 1
    fi
}

# 构建中继服务器
build_binary "中继服务器" "./cmd/relay" "./bin/relay-server"

# 构建CLI客户端
build_binary "CLI客户端" "./cmd/client" "./bin/cligool-client"

# 跨平台构建（可选）
if [ "$1" == "--all-platforms" ]; then
    echo ""
    echo "🌐 跨平台构建..."

    platforms=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")

    for platform in "${platforms[@]}"; do
        IFS='/' read -ra PARTS <<< "$platform"
        goos="${PARTS[0]}"
        goarch="${PARTS[1]}"

        output_name="./bin/relay-server-${goos}-${goarch}"
        if [ $goos == "windows" ]; then
            output_name="${output_name}.exe"
        fi

        echo -e "${YELLOW}🔨 构建 $platform...${NC}"
        CGO_ENABLED=0 GOOS=$goos GOARCH=$goarch go build -ldflags="-s -w" -o "$output_name" "./cmd/relay"

        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✅ $platform 构建成功${NC}"
        else
            echo -e "${RED}❌ $platform 构建失败${NC}"
        fi
    done
fi

echo ""
echo -e "${GREEN}🎉 构建完成！${NC}"
echo ""
echo "📦 构建产物:"
ls -lh bin/

echo ""
echo "🚀 快速启动:"
echo "   中继服务器: ./bin/relay-server"
echo "   CLI客户端:  ./bin/cligool-client -server http://localhost:8080"