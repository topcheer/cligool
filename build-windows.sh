#!/bin/bash

echo "🚀 CliGool 多平台构建脚本"
echo "================================"
echo ""

# 创建输出目录
mkdir -p bin

echo "📦 开始构建..."
echo ""

# 构建不同平台的版本
PLATFORMS=(
    "windows/amd64"
    "linux/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="${PLATFORM%/*}"
    GOARCH="${PLATFORM#*/}"
    OUTPUT="bin/cligool-${GOOS}-${GOARCH}"

    if [ "$GOOS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi

    echo "🔨 构建 ${GOOS}/${GOARCH}..."
    if GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$OUTPUT" ./cmd/client; then
        echo "   ✅ 成功: $OUTPUT"
        ls -lh "$OUTPUT" | awk '{print "   📦 大小: " $5}'
    else
        echo "   ❌ 失败: $PLATFORM"
    fi
    echo ""
done

echo "🎉 构建完成！"
echo ""
echo "生成的二进制文件:"
ls -lh bin/ | grep cligool | awk '{print "  " $9, "-", $5}'

echo ""
echo "💡 使用方法:"
echo "  Windows:   bin/cligool-windows-amd64.exe -server https://cligool.zty8.cn"
echo "  Linux:     ./bin/cligool-linux-amd64 -server https://cligool.zty8.cn"
echo "  macOS x64: ./bin/cligool-darwin-amd64 -server https://cligool.zty8.cn"
echo "  macOS ARM: ./bin/cligool-darwin-arm64 -server https://cligool.zty8.cn"
