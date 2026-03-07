# CliGool 安装包

本目录包含 CliGool 的各种平台安装包和安装脚本。

## 📦 支持的平台

### macOS
- **安装脚本**: `macos/install.sh`
- **安装位置**: `/usr/local/bin/cligool`
- **要求**: macOS 10.15+, sudo 权限

### Linux
- **安装脚本**: `linux/install.sh`
- **安装位置**: `/usr/local/bin/cligool`
- **支持架构**: amd64, arm64, 386, arm, ppc64le, riscv64, s390x, mips64le
- **要求**: sudo 权限

### Windows
- **安装程序**: `windows/cligool.iss` (Inno Setup 脚本)
- **生成文件**: `windows/output/cligool-setup.exe`
- **支持模式**:
  - **管理员模式**: 安装到 `Program Files`，系统级 PATH，完整功能
  - **用户模式**: 安装到 `%LOCALAPPDATA%`，用户级 PATH，基础功能
- **要求**:
  - 构建: Windows + Inno Setup
  - 安装: Windows 7+

## 🚀 快速安装

### macOS
```bash
# 下载并运行安装脚本
curl -sSL https://raw.githubusercontent.com/topcheer/cligool/main/installers/macos/install.sh | bash

# 或本地运行
cd installers/macos
sudo ./install.sh
```

### Linux
```bash
# 下载并运行安装脚本
curl -sSL https://raw.githubusercontent.com/topcheer/cligool/main/installers/linux/install.sh | bash

# 或本地运行
cd installers/linux
sudo ./install.sh
```

### Windows
```powershell
# 1. 下载安装程序
# 从 https://github.com/topcheer/cligool/releases 下载 cligool-setup.exe

# 2. 双击运行安装程序
# 右键 -> "以管理员身份运行" 获得完整功能

# 或使用命令行
cligool-setup.exe /VERYSILENT /SUPPRESSMSGBOXES
```

## 🔨 构建安装包

### 构建所有平台（推荐）
```bash
cd installers
chmod +x build-all.sh
./build-all.sh
```

### 单独构建

#### macOS
```bash
cd installers/macos
chmod +x install.sh
# 脚本已就绪，无需编译
```

#### Linux
```bash
cd installers/linux
chmod +x install.sh
# 脚本已就绪，无需编译
```

#### Windows

**在 Windows 上**:
```powershell
# 方法1：使用构建脚本
cd installers\windows
build.sh

# 方法2：使用 Inno Setup 编译器
iscc cligool.iss
```

**在 macOS/Linux 上**:
```bash
# 仅准备脚本，无法编译 .exe
cd installers/windows
chmod +x build.sh
./build.sh  # 会提示需要在 Windows 上编译
```

## 📋 安装包功能

### macOS/Linux
- ✅ 自动检测架构
- ✅ 下载正确的二进制文件
- ✅ 安装到系统路径
- ✅ 创建配置文件
- ✅ 创建桌面快捷方式（Linux）
- ✅ 添加到 PATH

### Windows (Inno Setup)
- ✅ 自动权限检测
- ✅ 管理员模式 vs 用户模式
- ✅ 动态安装路径选择
- ✅ 添加到 PATH（系统级或用户级）
- ✅ 创建快捷方式（开始菜单、桌面、快速启动）
- ✅ 右键菜单集成（仅管理员）
- ✅ 配置文件初始化
- ✅ 完整卸载支持

## 🎨 Windows 安装模式详解

### 管理员模式
- **安装位置**: `C:\Program Files\CliGool`
- **PATH**: 系统级 PATH（所有用户）
- **功能**:
  - ✅ 完整功能
  - ✅ 右键菜单集成
  - ✅ 所有用户开始菜单
  - ✅ 系统级配置
- **要求**: 以管理员身份运行安装程序

### 用户模式
- **安装位置**: `%LOCALAPPDATA%\CliGool` (通常 `C:\Users\YourName\AppData\Local\CliGool`)
- **PATH**: 用户级 PATH（仅当前用户）
- **功能**:
  - ✅ 基础功能
  - ✅ 用户开始菜单
  - ✅ 用户配置
  - ❌ 右键菜单集成（需要管理员权限）
- **要求**: 无特殊要求

## 🗑️ 卸载

### macOS
```bash
sudo rm -f /usr/local/bin/cligool
rm -rf ~/.cligool
rm -f ~/.cligool.json
```

### Linux
```bash
sudo rm -f /usr/local/bin/cligool
rm -rf ~/.cligool
rm -f ~/.cligool.json
rm -f ~/.local/share/applications/cligool.desktop
```

### Windows
- **方法1**: 控制面板 → 程序和功能 → CliGool → 卸载
- **方法2**: 开始菜单 → CliGool → 卸载 CliGool
- **方法3**: 运行安装目录中的 `unins000.exe`

## ⚙️ 高级配置

### 自定义安装位置

**macOS/Linux**:
```bash
# 下载后手动安装到自定义位置
INSTALL_DIR=/opt/cligool bash install.sh
```

**Windows**:
- 运行安装程序
- 在安装目录页面选择自定义路径

### 仅下载不安装

**macOS/Linux**:
```bash
# 下载到当前目录
curl -LO https://github.com/topcheer/cligool/releases/latest/download/cligool-darwin-arm64
chmod +x cligool-darwin-arm64
./cligool-darwin-arm64
```

**Windows**:
- 使用便携版（如果提供）
- 或下载 .zip 解压后直接运行

## 🐛 故障排除

### macOS/Linux

#### 问题：Permission denied
```bash
# 解决方案：添加执行权限
chmod +x installers/macos/install.sh
sudo ./installers/macos/install.sh
```

#### 问题：command not found after install
```bash
# 检查 PATH
echo $PATH

# 重新加载 shell 配置
source ~/.zshrc  # 或 source ~/.bash_profile

# 或手动指定路径
/usr/local/bin/cligool
```

### Windows

#### 问题：无法添加到 PATH
- **原因**: 权限不足
- **解决方案**:
  1. 以管理员身份运行安装程序
  2. 或手动添加到 PATH:
     - 系统属性 → 高级 → 环境变量
     - 编辑 Path 变量，添加安装目录

#### 问题：右键菜单不显示
- **原因**: 安装时未选择右键菜单集成，或不是管理员模式
- **解决方案**:
  1. 重新运行安装程序（管理员模式）
  2. 勾选"添加右键菜单"选项

#### 问题：安装程序无法运行
- **原因**: Windows Defender 或其他杀毒软件拦截
- **解决方案**:
  1. 临时禁用杀毒软件
  2. 添加到信任列表
  3. 下载时点击"保留"

## 📝 许可证

MIT License - 详见项目根目录的 LICENSE 文件

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

如果发现任何问题或有改进建议：
1. 创建 Issue 描述问题
2. 提供复现步骤
3. 包含日志输出（如有）

## 📚 相关文档

- [项目主页](https://github.com/topcheer/cligool)
- [使用指南](../USAGE_GUIDE_CN.md)
- [配置文件指南](../CONFIG_GUIDE.md)
- [命令行参数使用](../CMD_ARGS_USAGE.md)
