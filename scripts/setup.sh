#!/bin/bash

# CliGool 快速设置脚本

set -e

echo "🚀 CliGool 快速设置"
echo "=================="

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 创建必要的目录
echo "📁 创建项目目录..."
mkdir -p bin logs

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "⚠️  Go 未安装，将使用Docker构建"
else
    echo "✅ 检测到 Go $(go version | awk '{print $3}')"
fi

# 复制环境变量文件
if [ ! -f .env ]; then
    echo "📝 创建环境变量文件..."
    cp .env.example .env

    # 生成随机密码
    DB_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)
    JWT_SECRET=$(openssl rand -base64 32)

    # 更新.env文件
    sed -i.bak "s/cligool123/$DB_PASSWORD/g" .env
    sed -i.bak "s/your-super-secret-jwt-key-change-this/$JWT_SECRET/g" .env
    rm .env.bak

    echo "✅ 环境变量文件已创建 (.env)"
    echo "⚠️  请根据需要修改配置文件"
else
    echo "✅ 环境变量文件已存在"
fi

# 构建Docker镜像
echo "🐳 构建Docker镜像..."
docker-compose build

# 启动服务
echo "🚀 启动服务..."
docker-compose up -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 5

# 检查服务状态
echo "📊 服务状态:"
docker-compose ps

# 显示访问信息
echo ""
echo "✅ 设置完成！"
echo ""
echo "📍 本地访问:"
echo "   Web界面: http://localhost:8080"
echo "   API: http://localhost:8080/api"
echo ""
echo "🌐 Cloudflare Tunnel配置:"
echo "   1. 安装cloudflared: brew install cloudflare/cloudflare/cloudflared"
echo "   2. 创建tunnel: cloudflared tunnel create cligool"
echo "   3. 配置DNS: cloudflared tunnel route dns <tunnel-id> your-domain.com"
echo "   4. 启动tunnel: cloudflared tunnel run <tunnel-id>"
echo ""
echo "🔧 管理命令:"
echo "   查看日志: docker-compose logs -f"
echo "   停止服务: docker-compose down"
echo "   重启服务: docker-compose restart"
echo ""
echo "📖 更多配置信息请查看 docs/CONFIG.md"

# 询问是否立即查看日志
read -p "是否查看实时日志? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker-compose logs -f
fi