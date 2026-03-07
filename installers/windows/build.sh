#!/bin/bash
# Windows 安装包构建脚本

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "🔨 构建 Windows 安装包..."

# 检查是否在 macOS/Linux 上，需要使用 wine/Inno Setup
if [[ "$OSTYPE" == "darwin"* ]] || [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "⚠️  警告：在 macOS/Linux 上构建 Windows 安装包需要 Inno Setup"
    echo "   推荐在 Windows 上使用 Inno Setup 编译器（ISCC）"
    echo ""
    echo "构建方法："
    echo "1. Windows: 安装 Inno Setup，然后运行 iscc cligool.iss"
    echo "2. macOS/Linux: 使用 Wine + Inno Setup（未测试）"
    echo ""
    echo "跳过 Windows 安装包构建"
    exit 0
fi

# Windows 上的构建
if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    # 检查 Inno Setup 是否安装
    if ! command -v iscc &> /dev/null; then
        echo "❌ 错误：未找到 Inno Setup 编译器 (iscc)"
        echo "   请从以下地址下载安装："
        echo "   https://jrsoftware.org/isdl.php"
        exit 1
    fi

    echo "✅ 找到 Inno Setup 编译器"

    # 编译安装程序
    cd "$SCRIPT_DIR"
    iscc cligool.iss

    if [ $? -eq 0 ]; then
        echo "✅ Windows 安装包构建成功！"
        echo "📂 输出: installers/windows/output/cligool-setup.exe"
    else
        echo "❌ 构建失败"
        exit 1
    fi
fi
