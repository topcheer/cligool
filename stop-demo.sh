#!/bin/bash

# CliGool 停止演示脚本

echo "🛑 停止 CliGool 演示环境"
echo "======================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 停止CLI客户端
if [ -f "/tmp/cligool-demo-client.pid" ]; then
    CLIENT_PID=$(cat /tmp/cligool-demo-client.pid)
    if ps -p $CLIENT_PID > /dev/null 2>&1; then
        echo -e "${YELLOW}正在停止CLI客户端 (PID: $CLIENT_PID)...${NC}"
        kill $CLIENT_PID
        rm /tmp/cligool-demo-client.pid
        echo -e "${GREEN}✅ CLI客户端已停止${NC}"
    else
        echo -e "${YELLOW}⚠️  CLI客户端未运行${NC}"
        rm /tmp/cligool-demo-client.pid 2>/dev/null
    fi
else
    echo -e "${YELLOW}⚠️  未找到客户端PID文件${NC}"
fi

# 停止所有cligool-simple进程
pkill -f cligool-simple 2>/dev/null && echo -e "${GREEN}✅ 已清理所有CLI客户端进程${NC}" || true

echo ""

# 停止Docker服务
echo -e "${YELLOW}正在停止Docker服务...${NC}"
if docker-compose down 2>/dev/null; then
    echo -e "${GREEN}✅ Docker服务已停止${NC}"
else
    echo -e "${RED}❌ 停止Docker服务失败${NC}"
fi

echo ""

# 清理日志文件（可选）
echo -e "${YELLOW}是否清理日志文件？ (y/n)${NC}"
read -r response
if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    rm /tmp/cligool-demo-client.log 2>/dev/null || true
    echo -e "${GREEN}✅ 日志文件已清理${NC}"
fi

echo ""
echo -e "${GREEN}🎉 清理完成！${NC}"
echo ""
echo -e "${YELLOW}重新启动演示：${NC}"
echo "./demo.sh"
echo ""