#!/bin/bash
# 构建所有平台的安装包

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "🔨 CliGool 安装包构建工具"
echo "=============================="
echo ""

# 获取版本号
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
echo "📦 版本: $VERSION"
echo ""

# 显示菜单
show_menu() {
    echo "请选择要构建的安装包："
    echo "1) macOS 安装脚本"
    echo "2) Linux 安装脚本"
    echo "3) Windows Inno Setup 脚本（需要 Windows）"
    echo "4) 所有平台（除 Windows）"
    echo "5) 全部"
    echo "6) 退出"
    echo ""
    echo -n "请输入选项 [1-6]: "
}

# 构建 macOS 安装脚本
build_macos() {
    echo "📦 构建 macOS 安装脚本..."
    chmod +x "$SCRIPT_DIR/macos/install.sh"
    echo -e "${GREEN}✅ macOS 安装脚本已准备就绪${NC}"
    echo "   文件: installers/macos/install.sh"
    echo ""
}

# 构建 Linux 安装脚本
build_linux() {
    echo "📦 构建 Linux 安装脚本..."
    chmod +x "$SCRIPT_DIR/linux/install.sh"
    echo -e "${GREEN}✅ Linux 安装脚本已准备就绪${NC}"
    echo "   文件: installers/linux/install.sh"
    echo ""
}

# 构建 Windows 安装脚本
build_windows() {
    echo "📦 准备 Windows Inno Setup 脚本..."
    chmod +x "$SCRIPT_DIR/windows/build.sh"

    if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
        # Windows 环境
        echo "🔨 在 Windows 上编译安装程序..."
        cd "$SCRIPT_DIR/windows"
        ./build.sh
    else
        # 非 Windows 环境
        echo -e "${YELLOW}⚠️  Windows 安装包需要在 Windows 上编译${NC}"
        echo "   Inno Setup 脚本已准备就绪"
        echo "   文件: installers/windows/cligool.iss"
        echo ""
        echo "编译方法："
        echo "1. 在 Windows 上安装 Inno Setup"
        echo "2. 运行: iscc installers/windows/cligool.iss"
        echo "   或: installers/windows/build.sh"
    fi
    echo ""
}

# 构建所有平台（除 Windows）
build_all_except_windows() {
    echo "🔨 构建所有平台（除 Windows）..."
    build_macos
    build_linux
}

# 构建全部
build_all() {
    echo "🔨 构建所有平台..."
    build_macos
    build_linux
    build_windows
}

# 主循环
while true; do
    show_menu
    read -r choice

    case $choice in
        1)
            build_macos
            ;;
        2)
            build_linux
            ;;
        3)
            build_windows
            ;;
        4)
            build_all_except_windows
            ;;
        5)
            build_all
            ;;
        6)
            echo "👋 再见！"
            exit 0
            ;;
        *)
            echo -e "${RED}❌ 无效选项，请重试${NC}"
            echo ""
            ;;
    esac

    echo "按 Enter 继续..."
    read
done
