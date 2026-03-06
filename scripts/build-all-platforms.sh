#!/bin/bash

set -e

echo "🚀 CliGool 多平台构建脚本"
echo "========================================"
echo ""

# 输出目录
DOWNLOAD_DIR="web/downloads"
mkdir -p "$DOWNLOAD_DIR"

# 平台列表
PLATFORMS=(
    "windows/amd64"
    "linux/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

echo "📦 开始构建各个平台版本..."
echo ""

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="${PLATFORM%/*}"
    GOARCH="${PLATFORM#*/}"

    # 设置输出文件名
    if [ "$GOOS" = "windows" ]; then
        BINARY="cligool-windows-amd64.exe"
    else
        BINARY="cligool-${GOOS}-${GOARCH}"
    fi

    echo "🔨 构建 ${GOOS}/${GOARCH}..."

    # 构建二进制文件
    if GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0 go build -o "$DOWNLOAD_DIR/$BINARY" ./cmd/client; then
        SIZE=$(ls -lh "$DOWNLOAD_DIR/$BINARY" | awk '{print $5}')

        # Windows版本打包成zip
        if [ "$GOOS" = "windows" ]; then
            echo "   📦 打包Windows版本为ZIP..."
            cd "$DOWNLOAD_DIR"
            zip -q "cligool-windows-amd64.zip" "$BINARY"
            rm -f "$BINARY" # 删除原始exe文件
            cd - > /dev/null
            echo "   ✅ 成功: cligool-windows-amd64.zip (${SIZE})"
        else
            echo "   ✅ 成功: $BINARY (${SIZE})"
        fi
    else
        echo "   ❌ 失败: $PLATFORM"
    fi
done

echo ""
echo "🎉 构建完成！"
echo ""
echo "📦 生成的文件:"
ls -lh "$DOWNLOAD_DIR/" | awk 'NR>1 {printf "   %-30s %s\n", $9, $5}'

echo ""
echo "📊 版本信息:"
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
echo "   版本: $VERSION"
echo "   构建时间: $(date '+%Y-%m-%d %H:%M:%S')"

echo ""
echo "✅ 文件已准备就绪，可以启动服务器提供下载！"
