# CliGool 安装包实现总结

## ✅ 已实现的功能

### 1. 通用安装脚本
- **文件**: `install.sh` (项目根目录)
- **功能**: 自动检测操作系统并执行相应的安装
- **使用方法**:
  ```bash
  curl -sSL https://raw.githubusercontent.com/topcheer/cligool/main/install.sh | bash
  ```

### 2. macOS 安装包
- **文件**: `installers/macos/install.sh`
- **功能**:
  - ✅ 自动检测架构 (Intel/Apple Silicon)
  - ✅ 下载正确的二进制文件
  - ✅ 安装到 `/usr/local/bin/cligool`
  - ✅ 创建配置文件 `~/.cligool.json`
  - ✅ 添加到 PATH
  - ✅ 完整的卸载说明

### 3. Linux 安装包
- **文件**: `installers/linux/install.sh`
- **功能**:
  - ✅ 自动检测架构 (amd64, arm64, 386, arm, ppc64le, riscv64, s390x, mips64le)
  - ✅ 下载正确的二进制文件
  - ✅ 安装到 `/usr/local/bin/cligool`
  - ✅ 创建配置文件 `~/.cligool.json`
  - ✅ 创建桌面快捷方式
  - ✅ 添加到 PATH
  - ✅ 完整的卸载说明

### 4. Windows 安装包 (Inno Setup)
- **文件**: `installers/windows/cligool.iss`
- **功能**:
  - ✅ **双模式支持**:
    - 管理员模式: 安装到 Program Files，系统级 PATH
    - 普通用户模式: 安装到用户目录，用户级 PATH
  - ✅ 自动权限检测
  - ✅ 动态安装路径选择
  - ✅ 创建快捷方式（开始菜单、桌面、快速启动）
  - ✅ 右键菜单集成（仅管理员）
  - ✅ 添加到 PATH（系统级或用户级）
  - ✅ 配置文件初始化
  - ✅ 完整卸载支持
  - ✅ 多语言支持（中文/英文）
  - ✅ 自定义安装界面

## 📂 文件结构

```
cligool/
├── install.sh                           # 通用快速安装脚本
├── installers/                          # 安装包目录
│   ├── DESIGN.md                        # 完整设计文档
│   ├── README.md                        # 安装包使用说明
│   ├── build-all.sh                     # 构建所有安装包
│   ├── macos/
│   │   └── install.sh                   # macOS 安装脚本
│   ├── linux/
│   │   └── install.sh                   # Linux 安装脚本
│   └── windows/
│       ├── cligool.iss                  # Inno Setup 脚本
│       └── build.sh                     # Windows 构建脚本
```

## 🚀 使用方法

### 方法1：快速安装（推荐）

**macOS/Linux**:
```bash
curl -sSL https://raw.githubusercontent.com/topcheer/cligool/main/install.sh | bash
```

### 方法2：使用平台特定脚本

**macOS**:
```bash
curl -sSL https://raw.githubusercontent.com/topcheer/cligool/main/installers/macos/install.sh | bash
```

**Linux**:
```bash
curl -sSL https://raw.githubusercontent.com/topcheer/cligool/main/installers/linux/install.sh | bash
```

**Windows**:
```powershell
# 下载安装程序
Invoke-WebRequest -Uri "https://github.com/topcheer/cligool/releases/latest/download/cligool-setup.exe" -OutFile "cligool-setup.exe"

# 运行安装程序（推荐管理员模式）
.\cligool-setup.exe
```

### 方法3：本地安装

**macOS/Linux**:
```bash
cd installers/macos  # 或 installers/linux
sudo ./install.sh
```

**Windows**:
```powershell
# 需要先安装 Inno Setup
cd installers\windows
iscc cligool.iss
```

## 🎯 Windows 双模式详解

### 管理员模式
- **触发方式**: 以管理员身份运行安装程序
- **安装位置**: `C:\Program Files\CliGool`
- **PATH**: 系统级 PATH（所有用户可用）
- **功能**:
  - ✅ 完整功能
  - ✅ 右键菜单集成（"在 CliGool 中打开"）
  - ✅ 所有用户开始菜单
  - ✅ 系统级配置
- **适用场景**:
  - 个人电脑
  - 公司电脑（IT部门部署）
  - 需要完整功能的场景

### 普通用户模式
- **触发方式**: 以普通用户身份运行安装程序
- **安装位置**: `%LOCALAPPDATA%\CliGool`
- **PATH**: 用户级 PATH（仅当前用户）
- **功能**:
  - ✅ 基础功能
  - ✅ 用户开始菜单
  - ✅ 用户配置
  - ❌ 右键菜单集成（需要管理员权限）
- **适用场景**:
  - 受限环境（无管理员权限）
  - 临时使用
  - 公司电脑（无管理员权限）

## 🔧 构建安装包

### 构建所有平台
```bash
cd installers
chmod +x build-all.sh
./build-all.sh
```

### 单独构建
```bash
# macOS/Linux（脚本已就绪，无需编译）
chmod +x installers/macos/install.sh
chmod +x installers/linux/install.sh

# Windows（需要在 Windows 上编译）
cd installers/windows
iscc cligool.iss
```

## 📝 注意事项

### Windows Inno Setup
1. **需要先安装 Inno Setup**
   - 下载地址: https://jrsoftware.org/isdl.php
   - 安装后运行 `iscc` 命令编译

2. **图标文件**（可选，可自定义）
   - `installers/windows/cligool-icon.ico` - 安装程序图标
   - `installers/windows/cligool-setup.bmp` - 安装向导图片
   - `installers/windows/cligool-small.bmp` - 小图标

3. **输出文件**
   - 编译后生成: `installers/windows/output/cligool-setup.exe`

### macOS/Linux
1. **自动下载二进制**
   - 脚本会自动从 GitHub Releases 下载对应平台的二进制文件
   - 支持指定版本: `CLIGOOL_VERSION=v1.0.0 bash install.sh`

2. **需要 sudo 权限**
   - 安装到 `/usr/local/bin` 需要 sudo 权限
   - 配置文件安装到用户目录无需 sudo

## 🎨 安装体验

### macOS/Linux
1. 运行安装脚本
2. 自动检测系统架构
3. 下载对应的二进制文件
4. 安装到系统路径
5. 创建配置文件
6. 验证安装
7. 显示使用说明

### Windows
1. 双击运行安装程序
2. 选择安装模式（管理员/用户）
3. 阅读许可协议
4. 选择安装选项：
   - 创建桌面快捷方式
   - 添加右键菜单（仅管理员）
   - 创建配置文件
5. 安装进度
6. 完成页面（可立即运行）

## 📚 文档

- **设计文档**: `installers/DESIGN.md` - 完整的技术设计
- **使用说明**: `installers/README.md` - 安装包使用指南
- **主README**: `README.md` - 项目主页（需更新）

## 🔄 后续工作

### 待实现
1. **macOS .pkg 安装包**
   - 使用 Packages 创建
   - 提供拖拽安装的 .dmg

2. **Linux .deb/.rpm 包**
   - Debian/Ubuntu: .deb
   - RedHat/CentOS: .rpm

3. **自动更新功能**
   - 版本检查
   - 自动下载更新

4. **安装程序美化**
   - 自定义图标
   - 背景图片
   - 品牌元素

### 可选改进
1. **数字签名**
   - Windows: 代码签名证书
   - macOS: 公证

2. **安装验证**
   - 安装后测试
   - 自动卸载测试

3. **多语言支持**
   - 扩展到更多语言
   - 自动语言检测

## ✅ 测试清单

### macOS
- [x] Intel Mac 安装
- [x] Apple Silicon Mac 安装
- [x] 配置文件创建
- [x] PATH 添加
- [x] 卸载功能

### Linux
- [x] Ubuntu/Debian 安装
- [x] CentOS/RHEL 安装
- [x] 多架构支持
- [x] 配置文件创建
- [x] PATH 添加
- [x] 桌面快捷方式
- [x] 卸载功能

### Windows
- [x] Inno Setup 脚本编写
- [ ] 管理员模式测试（需要在 Windows 上测试）
- [ ] 普通用户模式测试（需要在 Windows 上测试）
- [ ] PATH 添加验证
- [ ] 右键菜单集成测试
- [ ] 卸载功能测试

## 🎉 总结

已成功实现跨平台安装包支持：

1. **✅ 通用快速安装脚本** - 一键安装所有平台
2. **✅ macOS 安装脚本** - 完整的 macOS 支持
3. **✅ Linux 安装脚本** - 支持 8 种架构
4. **✅ Windows Inno Setup** - 双模式安装（管理员/用户）
5. **✅ 完整文档** - 设计文档、使用说明、README

用户现在可以通过简单的命令快速安装 CliGool：
```bash
curl -sSL https://raw.githubusercontent.com/topcheer/cligool/main/install.sh | bash
```
