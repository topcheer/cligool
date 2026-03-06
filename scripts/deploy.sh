#!/bin/bash

# CliGool 快速部署脚本
# 可以在任何支持Docker的Linux服务器上运行

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🚀 CliGool 快速部署脚本"
echo "======================"
echo ""

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker 未安装${NC}"
    echo "请先安装Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose 未安装${NC}"
    echo "请先安装Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

echo -e "${GREEN}✅ Docker环境检查通过${NC}"

# 获取项目目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# 创建必要的目录
echo "📁 创建项目目录..."
mkdir -p logs

# 复制环境变量文件
if [ ! -f .env ]; then
    echo "📝 创建环境变量文件..."
    cp .env.example .env

    # 生成随机密码
    DB_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)
    JWT_SECRET=$(openssl rand -base64 32)

    # 更新.env文件
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s/cligool123/$DB_PASSWORD/g" .env
        sed -i '' "s/your-super-secret-jwt-key-change-this/$JWT_SECRET/g" .env
    else
        # Linux
        sed -i "s/cligool123/$DB_PASSWORD/g" .env
        sed -i "s/your-super-secret-jwt-key-change-this/$JWT_SECRET/g" .env
    fi

    echo -e "${GREEN}✅ 环境变量文件已创建${NC}"
else
    echo -e "${YELLOW}⚠️  环境变量文件已存在${NC}"
fi

# 停止现有服务
echo "🛑 停止现有服务..."
docker-compose down 2>/dev/null || true

# 构建镜像
echo "🔨 构建Docker镜像..."
docker-compose build

# 启动服务
echo "🚀 启动服务..."
docker-compose up -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 5

# 检查服务状态
echo ""
echo "📊 服务状态:"
docker-compose ps

# 显示访问信息
echo ""
echo -e "${GREEN}✅ 部署完成！${NC}"
echo ""
echo "📍 本地访问:"
echo "   curl http://localhost:8080/api/health"
echo ""
echo "🌐 下一步 - 配置Cloudflare Tunnel:"
echo "   1. 确保你有Cloudflare账号和域名"
echo "   2. 运行: ./scripts/cloudflare-tunnel.sh"
echo "   3. 或手动配置cloudflared"
echo ""
echo "🔧 管理命令:"
echo "   查看日志: docker-compose logs -f"
echo "   停止服务: docker-compose down"
echo "   重启服务: docker-compose restart"
echo "   查看状态: docker-compose ps"
echo ""
echo "📖 详细文档: docs/CONFIG.md"
echo ""

# 询问是否查看日志
read -p "是否查看实时日志? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker-compose logs -f
fi