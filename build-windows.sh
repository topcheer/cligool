#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "🚀 CliGool Windows 构建脚本"
echo "================================"
echo ""

mkdir -p bin

WINDOWS_ARCHES=(
  "amd64"
  "arm64"
)

for GOARCH in "${WINDOWS_ARCHES[@]}"; do
  OUTPUT="bin/cligool-windows-${GOARCH}.exe"
  echo "🔨 构建 windows/${GOARCH}..."
  CGO_ENABLED=0 GOOS=windows GOARCH="$GOARCH" go build -o "$OUTPUT" ./cmd/client
  echo "   ✅ 成功: $OUTPUT"
  ls -lh "$OUTPUT" | awk '{print "   📦 大小: " $5}'
  echo ""
done

echo "🎉 Windows 构建完成！"
echo ""
echo "生成的二进制文件:"
ls -lh bin/cligool-windows-*.exe | awk '{print "  " $9, "-", $5}'
echo ""
echo "💡 使用方法:"
echo "  AMD64: ./bin/cligool-windows-amd64.exe -server https://cligool.ystone.us"
echo "  ARM64: ./bin/cligool-windows-arm64.exe -server https://cligool.ystone.us"
