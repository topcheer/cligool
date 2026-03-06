#!/bin/bash

echo "🚀 部署更新到远程服务器"
echo "================================"
echo ""

# 检查是否有远程服务器配置
if [ -z "$REMOTE_SERVER" ]; then
    echo "❌ 错误: 没有设置远程服务器地址"
    echo ""
    echo "使用方法:"
    echo "  REMOTE_SERVER=user@your-server.com ./deploy-to-remote.sh"
    echo ""
    echo "或者手动部署到远程服务器:"
    echo "1. 将当前代码复制到远程服务器"
    echo "2. 在远程服务器上运行:"
    echo "   docker-compose build"
    echo "   docker-compose up -d"
    exit 1
fi

echo "📦 准备部署到: $REMOTE_SERVER"
echo ""

# 在远程服务器上执行部署
ssh "$REMOTE_SERVER" << 'ENDSSH'
    cd /path/to/cligool  # 修改为实际路径
    git pull origin main  # 或使用其他方式更新代码
    docker-compose build
    docker-compose up -d
    echo "✅ 部署完成"
ENDSSH

echo ""
echo "✅ 远程部署完成"
