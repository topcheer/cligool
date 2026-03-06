#!/bin/bash

# CliGool 真实终端启动脚本
# 解决PTY权限问题

echo "🔍 检测终端环境..."

# 检查是否在真实终端中
if [ ! -t 0 ]; then
    echo "❌ 错误：不在真实终端环境中"
    echo ""
    echo "💡 解决方法："
    echo "   1. 打开 Terminal.app 或 iTerm2"
    echo "   2. cd到当前目录: cd $(pwd)"
    echo "   3. 运行: ./bin/cligool-client"
    echo ""
    exit 1
fi

# 检查TERM环境变量
if [ -z "$TERM" ]; then
    echo "❌ 错误：TERM环境变量未设置"
    echo "   请在真实终端中运行此程序"
    exit 1
fi

echo "✅ 终端环境检查通过"
echo "   TERM: $TERM"
echo ""

# 检查PTY设备权限
if [ ! -c /dev/ptmx ]; then
    echo "❌ 错误：/dev/ptmx不存在或不是字符设备"
    exit 1
fi

echo "✅ PTY设备检查通过"
echo ""

# 启动客户端
echo "🚀 启动CliGool客户端..."
echo ""

# 直接运行客户端
./bin/cligool-simple -connect-only "$@"