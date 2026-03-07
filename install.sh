#!/bin/bash
# CliGool 快速安装脚本
# 自动检测操作系统并执行相应的安装

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}"
echo "╔═══════════════════════════════════════════════════════════╗"
echo "║                    🚀 CliGool 安装程序                    ║"
echo "║                   远程终端解决方案                         ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo ""

# 检测操作系统
detect_os() {
    case "$(uname -s)" in
        Darwin*)
            echo "macOS"
            ;;
        Linux*)
            echo "linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo "windows"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# macOS 安装
install_macos() {
    echo "🍎 检测到 macOS 系统"
    echo ""

    # 检查是否有 curl 或 wget
    if command -v curl &> /dev/null; then
        DOWNLOADER="curl -sSL"
    elif command -v wget &> /dev/null; then
        DOWNLOADER="wget -qO-"
    else
        echo -e "${RED}❌ 错误：需要 curl 或 wget${NC}"
        exit 1
    fi

    echo "📥 下载并运行 macOS 安装脚本..."
    echo ""

    # 使用本地脚本（如果存在）
    if [ -f "installers/macos/install.sh" ]; then
        bash installers/macos/install.sh
    else
        # 使用在线脚本
        $DOWNLOADER https://raw.githubusercontent.com/topcheer/cligool/main/installers/macos/install.sh | bash
    fi
}

# Linux 安装
install_linux() {
    echo "🐧 检测到 Linux 系统"
    echo ""

    # 检查是否有 curl 或 wget
    if command -v curl &> /dev/null; then
        DOWNLOADER="curl -sSL"
    elif command -v wget &> /dev/null; then
        DOWNLOADER="wget -qO-"
    else
        echo -e "${RED}❌ 错误：需要 curl 或 wget${NC}"
        exit 1
    fi

    echo "📥 下载并运行 Linux 安装脚本..."
    echo ""

    # 使用本地脚本（如果存在）
    if [ -f "installers/linux/install.sh" ]; then
        bash installers/linux/install.sh
    else
        # 使用在线脚本
        $DOWNLOADER https://raw.githubusercontent.com/topcheer/cligool/main/installers/linux/install.sh | bash
    fi
}

# Windows 安装
install_windows() {
    echo "🪟 检测到 Windows 系统"
    echo ""

    echo -e "${YELLOW}Windows 安装方法：${NC}"
    echo ""
    echo "1. 下载安装程序："
    echo "   https://github.com/topcheer/cligool/releases/latest/download/cligool-setup.exe"
    echo ""
    echo "2. 双击运行安装程序"
    echo "   右键 -> '以管理员身份运行' 获得完整功能"
    echo ""
    echo "3. 或使用 PowerShell 下载并安装："
    echo ""
    echo "   # 下载安装程序"
    echo "   Invoke-WebRequest -Uri 'https://github.com/topcheer/cligool/releases/latest/download/cligool-setup.exe' -OutFile 'cligool-setup.exe'"
    echo "   # 运行安装程序"
    echo "   .\cligool-setup.exe"
    echo ""
}

# 未知系统
install_unknown() {
    echo -e "${RED}❌ 错误：无法检测操作系统${NC}"
    echo ""
    echo "支持的操作系统："
    echo "  - macOS 10.15+"
    echo "  - Linux (所有主流发行版)"
    echo "  - Windows 7+"
    echo ""
    echo "手动下载地址："
    echo "  https://github.com/topcheer/cligool/releases"
}

# 主安装流程
main() {
    OS=$(detect_os)

    case "$OS" in
        macOS)
            install_macos
            ;;
        linux)
            install_linux
            ;;
        windows)
            install_windows
            ;;
        *)
            install_unknown
            exit 1
            ;;
    esac
}

# 显示帮助信息
show_help() {
    cat << EOF
用法: $0 [选项]

选项:
  -h, --help     显示此帮助信息
  -v, --version  显示版本信息
  -d, --debug    启用调试模式

环境变量:
  CLIGOOL_VERSION  指定要安装的版本（默认: latest）

示例:
  # 安装最新版本
  $0

  # 安装特定版本
  CLIGOOL_VERSION=v1.0.0 $0

  # 启用调试模式
  $0 --debug

更多信息:
  https://github.com/topcheer/cligool
EOF
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--version)
            echo "CliGool 安装脚本 v1.0.0"
            exit 0
            ;;
        -d|--debug)
            set -x
            shift
            ;;
        *)
            echo -e "${RED}❌ 未知选项: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 运行主安装流程
main
