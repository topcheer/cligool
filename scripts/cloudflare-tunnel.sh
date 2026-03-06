#!/bin/bash

# Cloudflare Tunnel 快速配置脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🌐 CliGool Cloudflare Tunnel 配置"
echo "================================="

# 检查cloudflared是否安装
if ! command -v cloudflared &> /dev/null; then
    echo -e "${RED}❌ cloudflared 未安装${NC}"
    echo ""
    echo "请安装 cloudflared:"
    echo "  macOS:   brew install cloudflare/cloudflare/cloudflared"
    echo "  Linux:   wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb"
    echo "           sudo dpkg -i cloudflared-linux-amd64.deb"
    echo "  Windows: 下载 https://github.com/cloudflare/cloudflared/releases"
    exit 1
fi

echo -e "${GREEN}✅ cloudflared 已安装${NC}"

# 检查是否已登录
echo ""
echo "📝 检查Cloudflare登录状态..."
if ! cloudflared tunnel info &> /dev/null; then
    echo "请先登录Cloudflare..."
    cloudflared tunnel login
fi

# 获取域名
read -p "请输入你的域名 (例如: cligool.example.com): " DOMAIN

if [ -z "$DOMAIN" ]; then
    echo -e "${RED}❌ 域名不能为空${NC}"
    exit 1
fi

# 创建tunnel
echo ""
echo "🚀 创建Cloudflare Tunnel..."
TUNNEL_NAME="cligool-$(date +%s)"

if cloudflared tunnel create "$TUNNEL_NAME"; then
    TUNNEL_ID=$(cloudflared tunnel list | grep "$TUNNEL_NAME" | awk '{print $1}')
    echo -e "${GREEN}✅ Tunnel创建成功: $TUNNEL_ID${NC}"
else
    echo -e "${RED}❌ Tunnel创建失败${NC}"
    exit 1
fi

# 配置DNS
echo ""
echo "📡 配置DNS记录..."
if cloudflared tunnel route dns "$TUNNEL_ID" "$DOMAIN"; then
    echo -e "${GREEN}✅ DNS记录配置成功${NC}"
else
    echo -e "${RED}❌ DNS记录配置失败${NC}"
    exit 1
fi

# 生成配置文件
echo ""
echo "📝 生成配置文件..."

cat > cloudflare-tunnel.yml <<EOF
# Cloudflare Tunnel 配置
# Tunnel ID: $TUNNEL_ID
# 域名: $DOMAIN

tunnel: $TUNNEL_ID
credentials-file: ~/.cloudflared/${TUNNEL_ID}.json

ingress:
  - hostname: $DOMAIN
    service: http://localhost:8080
  - service: http_status:404
EOF

echo -e "${GREEN}✅ 配置文件已生成: cloudflare-tunnel.yml${NC}"

# 启动tunnel
echo ""
echo "🚀 启动Cloudflare Tunnel..."
echo "按 Ctrl+C 停止tunnel"
echo ""

cloudflared tunnel --config cloudflare-tunnel.yml run

# 如果需要后台运行，用户可以使用:
# nohup cloudflared tunnel --config cloudflare-tunnel.yml run > /dev/null 2>&1 &
# 或使用 systemd/screen/tmux