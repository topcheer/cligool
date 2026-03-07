# CliGool 平台支持

本版本支持**30个操作系统和架构组合**，覆盖所有主流平台。

## 🌟 已构建的平台

### macOS (2个)
- ✅ **darwin-amd64** - Intel处理器（Mac mini 2018, MacBook Pro 2019等）
- ✅ **darwin-arm64** - Apple Silicon（M1/M2/M3/M4）

### Windows (2个)
- ✅ **windows-amd64** - Windows 64位（ZIP格式）
- ✅ **windows-arm64** - Windows ARM 64位（Surface Pro X等）

### Linux (13个)
- ✅ **linux-amd64** - x86_64（Ubuntu, Debian, CentOS等）
- ✅ **linux-arm64** - ARM 64位（树莓派4/5, ARM服务器）
- ✅ **linux-386** - x86 32位（老式PC）
- ✅ **linux-arm** - ARM 32位（树莓派Zero/2B/3B）
- ✅ **linux-armbe** - ARM64 Big-Endian
- ✅ **linux-ppc64le** - PowerPC Little-Endian
- ✅ **linux-ppc64** - PowerPC Big-Endian
- ✅ **linux-riscv64** - RISC-V 64位
- ✅ **linux-s390x** - IBM System z
- ✅ **linux-mips** - MIPS 32位
- ✅ **linux-mips64le** - MIPS 64位 Little-Endian
- ✅ **linux-mips64** - MIPS 64位 Big-Endian
- ✅ **linux-loong64** - 龙芯（LoongArch）

### FreeBSD (5个)
- ✅ **freebsd-amd64** - FreeBSD 64位
- ✅ **freebsd-arm64** - FreeBSD ARM 64位
- ✅ **freebsd-386** - FreeBSD 32位
- ✅ **freebsd-arm** - FreeBSD ARM 32位
- ✅ **freebsd-riscv64** - FreeBSD RISC-V

### OpenBSD (2个)
- ✅ **openbsd-amd64** - OpenBSD 64位
- ✅ **openbsd-arm64** - OpenBSD ARM 64位

### NetBSD (4个)
- ✅ **netbsd-amd64** - NetBSD 64位
- ✅ **netbsd-arm64** - NetBSD ARM 64位
- ✅ **netbsd-arm** - NetBSD ARM 32位
- ✅ **netbsd-386** - NetBSD 32位

### DragonFlyBSD (1个)
- ✅ **dragonfly-amd64** - DragonFlyBSD 64位

### 中继服务器 (1个)
- ✅ **relay-server-linux-amd64** - Linux服务器组件

## 📊 平台覆盖率

| 操作系统 | 架构覆盖数 | 总覆盖率 |
|---------|----------|---------|
| Windows | 2/2 | 100% ✅ |
| macOS | 2/2 | 100% ✅ |
| Linux | 13/13 | 100% ✅ |
| FreeBSD | 5/5 | 100% ✅ |
| OpenBSD | 2/5 | 40% ⚠️ |
| NetBSD | 4/4 | 100% ✅ |
| DragonFlyBSD | 1/2 | 50% ⚠️ |
| **总计** | **29/33** | **88%** |

**说明**：
- OpenBSD仅支持amd64和arm64（386/arm/riscv64受pty库限制）
- DragonFlyBSD仅支持amd64（arm64不被Go支持）
- 总计实际支持：30个平台（包括Windows）

## 🚀 快速开始

### 1. 从GitHub Releases下载

访问 [https://github.com/topcheer/cligool/releases](https://github.com/topcheer/cligool/releases) 下载对应平台的二进制文件。

### 2. 使用安装脚本（推荐）

**macOS/Linux**:
```bash
curl -sSL https://raw.githubusercontent.com/topcheer/cligool/main/install.sh | bash
```

**Windows**:
```powershell
# 下载ZIP文件并解压
# 从 https://github.com/topcheer/cligool/releases/latest 下载
```

### 3. 本地构建所有平台
```bash
./build-all.sh
```

### 4. 本地测试
```bash
# 自动检测当前平台
./bin/cligool-* -server http://localhost:8081
```

## 📦 二进制文件命名规则

格式：`cligool-{GOOS}-{GOARCH}[.exe]`

示例：
- `cligool-linux-amd64` - Linux AMD64
- `cligool-darwin-arm64` - macOS ARM64
- `cligool-windows-amd64.exe` - Windows AMD64
- `cligool-freebsd-riscv64` - FreeBSD RISC-V

## 🔧 技术细节

### 构建方式
- **交叉编译**：使用Go的交叉编译功能（`CGO_ENABLED=0`）
- **无依赖部署**：所有二进制都是静态编译，无需额外依赖
- **优化体积**：使用UPX压缩（可选）和编译优化

### 系统要求
- **构建**：Go 1.24+
- **运行**：对应的操作系统和架构

### 编码支持
- **Unix/Linux/macOS/*BSD**：UTF-8原生支持
- **Windows**：自动检测控制台编码（GBK/Shift-JIS等）并转换为UTF-8

### 终端支持
- **Unix-like系统**：使用PTY（伪终端），支持完整的终端特性
- **Windows**：使用ConPTY（Windows 10+）或cmd.exe管道，支持基本终端特性

## 🎯 特性

### 双端同步
- ✅ 本地终端实时显示所有输出
- ✅ 本地终端可以直接输入命令
- ✅ Web终端同步显示所有内容
- ✅ Web终端也可以输入命令
- ✅ 两边完全同步，互不干扰

### 使用场景
1. **本地工作**：在本地终端操作，Web端查看演示
2. **远程协作**：在Web端操作，本地终端监控
3. **教学演示**：一边操作一边给其他人看
4. **技术支持**：本地查看，远程操作
5. **多设备访问**：多个Web客户端同时连接同一终端

## 🐳 Docker支持

所有平台的客户端二进制文件都打包在Docker镜像中，可以通过以下方式获取：

```bash
# 启动Docker服务
docker-compose up -d

# 从容器中复制所有平台的二进制文件
docker cp cligool-relay:/app/web/downloads ./
```

## 📝 平台特定说明

### Windows平台
- 使用ZIP压缩包分发，无需安装程序
- 支持ConPTY（Windows 10 1809+）
- 自动检测和控制台编码转换
- 依赖：无（纯Go编译）

### Linux平台
- 支持所有主流Linux发行版
- 包含所有标准架构（x86, ARM, PowerPC, RISC-V, MIPS等）
- 支持龙芯处理器（LoongArch）

### *BSD平台
- 支持FreeBSD、OpenBSD、NetBSD、DragonFlyBSD
- 覆盖所有支持的架构
- 完整的PTY支持

## 📄 许可证
MIT License - 详见 [LICENSE](LICENSE) 文件
