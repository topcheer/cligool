#!/bin/bash
# 快速测试本地构建的客户端

SESSION_ID="${1:-test-session}"
SERVER_URL="${2:-http://localhost:8081}"

echo "🚀 启动CliGool客户端（本地构建）"
echo "📋 会话ID: $SESSION_ID"
echo "🌐 服务器: $SERVER_URL"
echo ""

# 检测当前平台
PLATFORM=$(uname -s)
ARCH=$(uname -m)

case "$PLATFORM" in
    Darwin)
        if [[ "$ARCH" == "arm64" ]]; then
            CLIENT="bin/cligool-darwin-arm64"
        else
            CLIENT="bin/cligool-darwin-amd64"
        fi
        ;;
    Linux)
        if [[ "$ARCH" == "aarch64" ]]; then
            CLIENT="bin/cligool-linux-arm64"
        elif [[ "$ARCH" == "i386" ]] || [[ "$ARCH" == "i686" ]]; then
            CLIENT="bin/cligool-linux-386"
        elif [[ "$ARCH" == "armv7l" ]]; then
            CLIENT="bin/cligool-linux-arm"
        else
            CLIENT="bin/cligool-linux-amd64"
        fi
        ;;
    *)
        echo "❌ 不支持的平台: $PLATFORM $ARCH"
        exit 1
        ;;
esac

echo "📦 使用客户端: $CLIENT"
echo ""

if [ ! -f "$CLIENT" ]; then
    echo "❌ 客户端不存在: $CLIENT"
    echo "💡 请先运行 ./build-all.sh 构建所有平台"
    exit 1
fi

# 启动客户端
exec "$CLIENT" -server "$SERVER_URL" -session "$SESSION_ID"

