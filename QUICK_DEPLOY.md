#!/bin/bash
# CliGool 快速部署脚本
# 使用方法：./QUICK_DEPLOY.sh [platform]

set -e

echo "🚀 CliGool Relay Server 快速部署"
echo ""
echo "请选择部署平台："
echo "1) Render  - 最简单，完全免费，推荐新手"
echo "2) Fly.io  - 性能最好，无冷启动，推荐个人使用"
echo "3) Koyeb   - $5.50免费额度，全球CDN"
echo "4) Railway - 功能最全，$5免费额度"
echo ""
read -p "请输入选择 (1-4): " choice

case $choice in
  1)
    echo ""
    echo "📝 部署到 Render"
    echo ""
    echo "步骤："
    echo "1. 访问 https://dashboard.render.com/register"
    echo "2. 使用GitHub账号登录"
    echo "3. 点击 'New +' -> 'Web Service'"
    echo "4. 连接你的GitHub仓库"
    echo "5. 选择 Dockerfile.multiarch"
    echo "6. 设置环境变量："
    echo "   - DATABASE_URL (稍后配置)"
    echo "   - REDIS_URL (稍后配置)"
    echo "   - RELAY_HOST = 0.0.0.0"
    echo "   - RELAY_PORT = 8080"
    echo "7. 创建PostgreSQL数据库"
    echo "8. 创建Redis实例"
    echo "9. 将数据库连接URL添加到环境变量"
    echo "10. 部署！"
    echo ""
    echo "详细指南：cat docs/DEPLOYMENT_GUIDES.md"
    ;;
  2)
    echo ""
    echo "📝 部署到 Fly.io"
    echo ""
    echo "步骤："
    echo "1. 安装flyctl:"
    echo "   curl -L https://fly.io/install.sh | sh"
    echo ""
    echo "2. 登录:"
    echo "   flyctl auth login"
    echo ""
    echo "3. 在项目目录运行:"
    echo "   flyctl launch"
    echo ""
    echo "4. 创建数据库:"
    echo "   flyctl postgres create --name cligool-postgres"
    echo "   flyctl redis create --name cligool-redis"
    echo ""
    echo "5. 附加数据库:"
    echo "   flyctl postgres attach cligool-postgres"
    echo "   flyctl secrets set REDIS_URL=redis://..."
    echo ""
    echo "6. 部署:"
    echo "   flyctl deploy"
    echo ""
    echo "详细指南：cat docs/DEPLOYMENT_GUIDES.md"
    ;;
  3)
    echo ""
    echo "📝 部署到 Koyeb"
    echo ""
    echo "步骤："
    echo "1. 访问 https://www.koyeb.com"
    echo "2. 注册账号（支持GitHub登录）"
    echo "3. 点击 'Create App'"
    echo "4. 选择 GitHub 仓库"
    echo "5. 配置 Dockerfile.multiarch"
    echo "6. 创建 PostgreSQL 和 Redis 服务"
    echo "7. 设置环境变量"
    echo "8. 部署！"
    echo ""
    echo "详细指南：cat docs/DEPLOYMENT_GUIDES.md"
    ;;
  4)
    echo ""
    echo "📝 部署到 Railway"
    echo ""
    echo "步骤："
    echo "1. 访问 https://railway.app"
    echo "2. 使用GitHub账号登录"
    echo "3. 点击 'New Project' -> 'Deploy from GitHub repo'"
    echo "4. 选择 cligool 仓库"
    echo "5. Railway 会自动检测 railway.toml 配置"
    echo "6. 添加 PostgreSQL 和 Redis 服务"
    echo "7. 设置环境变量"
    echo "8. 部署！"
    echo ""
    echo "详细指南：cat docs/DEPLOYMENT_GUIDES.md"
    ;;
  *)
    echo "无效选择"
    exit 1
    ;;
esac

echo ""
echo "📚 更多信息："
echo "   - 部署指南: docs/DEPLOYMENT_GUIDES.md"
echo "   - 配置说明: docs/CONFIG.md"
echo "   - 使用指南: docs/USAGE_GUIDE.md"
