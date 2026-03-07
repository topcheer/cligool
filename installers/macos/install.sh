#!/bin/bash
# CliGool macOS 安装脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
VERSION="${CLIGOOL_VERSION:-latest}"
DOWNLOAD_BASE="https://github.com/topcheer/cligool/releases"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.cligool"

echo "🚀 CliGool macOS 安装程序"
echo "========================="
echo ""

# 检测架构
echo "🔍 检测系统架构..."
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)
        BINARY="cligool-darwin-amd64"
        ARCH_NAME="Intel (x86_64)"
        ;;
    arm64)
        BINARY="cligool-darwin-arm64"
        ARCH_NAME="Apple Silicon (ARM64)"
        ;;
    *)
        echo -e "${RED}❌ 错误：不支持的架构 $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}✅ 检测到架构：$ARCH_NAME${NC}"
echo "📦 将下载二进制文件：$BINARY"
echo ""

# 检查是否有 sudo 权限
if [ "$EUID" -ne 0 ]; then
    echo -e "${YELLOW}⚠️  注意：安装到 $INSTALL_DIR 需要 sudo 权限${NC}"
    SUDO="sudo"
else
    SUDO=""
fi

# 构建下载 URL
if [ "$VERSION" = "latest" ]; then
    DOWNLOAD_URL="$DOWNLOAD_BASE/latest/download/$BINARY"
else
    DOWNLOAD_URL="$DOWNLOAD_BASE/download/$VERSION/$BINARY"
fi

echo "📥 下载地址：$DOWNLOAD_URL"
echo ""

# 创建临时目录
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# 下载二进制文件
echo "⬇️  正在下载 CliGool..."
if command -v curl &> /dev/null; then
    curl -L --progress-bar "$DOWNLOAD_URL" -o "$TEMP_DIR/cligool"
elif command -v wget &> /dev/null; then
    wget -q --show-progress "$DOWNLOAD_URL" -O "$TEMP_DIR/cligool"
else
    echo -e "${RED}❌ 错误：需要 curl 或 wget 来下载文件${NC}"
    exit 1
fi

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ 下载失败${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 下载完成${NC}"
echo ""

# 安装二进制文件
echo "📦 安装到 $INSTALL_DIR..."
$SUDO install -m 755 "$TEMP_DIR/cligool" "$INSTALL_DIR/cligool"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 安装成功${NC}"
else
    echo -e "${RED}❌ 安装失败${NC}"
    exit 1
fi

echo ""

# 创建配置目录（如果不存在）
if [ ! -d "$CONFIG_DIR" ]; then
    echo "📝 创建配置目录：$CONFIG_DIR"
    mkdir -p "$CONFIG_DIR"
fi

# 创建默认配置文件（如果不存在）
CONFIG_FILE="$HOME/.cligool.json"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "📝 创建默认配置文件：$CONFIG_FILE"
    cat > "$CONFIG_FILE" << EOF
{
  "server": "https://cligool.zty8.cn",
  "cols": 0,
  "rows": 0
}
EOF
    echo -e "${GREEN}✅ 配置文件已创建${NC}"
else
    echo "ℹ️  配置文件已存在，跳过创建"
fi

echo ""

# 验证安装
echo "🔍 验证安装..."
if command -v cligool &> /dev/null; then
    echo -e "${GREEN}✅ CliGool 已成功安装到系统${NC}"
    echo "📍 安装位置：$INSTALL_DIR/cligool"
    echo ""
    echo "🎉 安装完成！"
    echo ""
    echo "📖 使用方法："
    echo "   cligool                    # 使用默认配置启动"
    echo "   cligool -help              # 查看帮助"
    echo "   cligool -cmd claude        # 运行 AI CLI 工具"
    echo ""
    echo "📚 更多信息："
    echo "   - 配置文件：$CONFIG_FILE"
    echo "   - 文档：https://github.com/topcheer/cligool"
    echo ""
    echo "💡 提示："
    echo "   - 如需自动补全，运行: eval \"\$(cligool --completion)\""
    echo "   - 如需升级，重新运行此安装脚本"
else
    echo -e "${YELLOW}⚠️  警告：CLI 命令不可用${NC}"
    echo ""
    echo "可能的原因："
    echo "1. $INSTALL_DIR 不在 PATH 中"
    echo "2. 需要重新登录或重启终端"
    echo ""
    echo "解决方案："
    echo "1. 检查 PATH：echo \$PATH"
    echo "2. 手动运行：$INSTALL_DIR/cligool"
    echo "3. 添加到 shell 配置文件（~/.zshrc 或 ~/.bash_profile）："
    echo "   export PATH=\"\$PATH:$INSTALL_DIR\""
fi

echo ""
echo "🗑️  卸载方法："
echo "   $SUDO rm -f $INSTALL_DIR/cligool"
echo "   rm -rf $CONFIG_DIR"
echo "   rm -f $CONFIG_FILE"
