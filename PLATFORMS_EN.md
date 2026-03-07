# CliGool Platform Support

This version supports **30 operating system and architecture combinations**, covering all mainstream platforms.

## 🌟 Supported Platforms

### macOS (2 platforms)
- ✅ **darwin-amd64** - Intel processors (Mac mini 2018, MacBook Pro 2019, etc.)
- ✅ **darwin-arm64** - Apple Silicon (M1/M2/M3/M4)

### Windows (2 platforms)
- ✅ **windows-amd64** - Windows 64-bit (ZIP format)
- ✅ **windows-arm64** - Windows ARM 64-bit (Surface Pro X, etc.)

### Linux (13 platforms)
- ✅ **linux-amd64** - x86_64 (Ubuntu, Debian, CentOS, etc.)
- ✅ **linux-arm64** - ARM 64-bit (Raspberry Pi 4/5, ARM servers)
- ✅ **linux-386** - x86 32-bit (older PCs)
- ✅ **linux-arm** - ARM 32-bit (Raspberry Pi Zero/2B/3B)
- ✅ **linux-armbe** - ARM64 Big-Endian
- ✅ **linux-ppc64le** - PowerPC Little-Endian
- ✅ **linux-ppc64** - PowerPC Big-Endian
- ✅ **linux-riscv64** - RISC-V 64-bit
- ✅ **linux-s390x** - IBM System z
- ✅ **linux-mips** - MIPS 32-bit
- ✅ **linux-mips64le** - MIPS 64-bit Little-Endian
- ✅ **linux-mips64** - MIPS 64-bit Big-Endian
- ✅ **linux-loong64** - LoongArch

### FreeBSD (5 platforms)
- ✅ **freebsd-amd64** - FreeBSD 64-bit
- ✅ **freebsd-arm64** - FreeBSD ARM 64-bit
- ✅ **freebsd-386** - FreeBSD 32-bit
- ✅ **freebsd-arm** - FreeBSD ARM 32-bit
- ✅ **freebsd-riscv64** - FreeBSD RISC-V

### OpenBSD (2 platforms)
- ✅ **openbsd-amd64** - OpenBSD 64-bit
- ✅ **openbsd-arm64** - OpenBSD ARM 64-bit

### NetBSD (4 platforms)
- ✅ **netbsd-amd64** - NetBSD 64-bit
- ✅ **netbsd-arm64** - NetBSD ARM 64-bit
- ✅ **netbsd-arm** - NetBSD ARM 32-bit
- ✅ **netbsd-386** - NetBSD 32-bit

### DragonFlyBSD (1 platform)
- ✅ **dragonfly-amd64** - DragonFlyBSD 64-bit

### Relay Server (1 platform)
- ✅ **relay-server-linux-amd64** - Linux server component

## 📊 Platform Coverage

| Operating System | Architecture Count | Total Coverage |
|-----------------|-------------------|----------------|
| Windows | 2/2 | 100% ✅ |
| macOS | 2/2 | 100% ✅ |
| Linux | 13/13 | 100% ✅ |
| FreeBSD | 5/5 | 100% ✅ |
| OpenBSD | 2/5 | 40% ⚠️ |
| NetBSD | 4/4 | 100% ✅ |
| DragonFlyBSD | 1/2 | 50% ⚠️ |
| **Total** | **29/33** | **88%** |

**Notes**:
- OpenBSD only supports amd64 and arm64 (386/arm/riscv64 limited by pty library)
- DragonFlyBSD only supports amd64 (arm64 not supported by Go)
- Total actually supported: 30 platforms (including Windows)

## 🚀 Quick Start

### 1. Download from GitHub Releases

Visit [https://github.com/topcheer/cligool/releases](https://github.com/topcheer/cligool/releases) to download the binary for your platform.

### 2. Use Installation Script (Recommended)

**macOS/Linux**:
```bash
curl -sSL https://raw.githubusercontent.com/topcheer/cligool/main/install.sh | bash
```

**Windows**:
```powershell
# Download ZIP file and extract
# Download from https://github.com/topcheer/cligool/releases/latest
```

### 3. Build All Platforms Locally
```bash
./build-all.sh
```

### 4. Test Locally
```bash
# Auto-detect current platform
./bin/cligool-* -server http://localhost:8081
```

## 📦 Binary File Naming Convention

Format: `cligool-{GOOS}-{GOARCH}[.exe]`

Examples:
- `cligool-linux-amd64` - Linux AMD64
- `cligool-darwin-arm64` - macOS ARM64
- `cligool-windows-amd64.exe` - Windows AMD64
- `cligool-freebsd-riscv64` - FreeBSD RISC-V

## 🔧 Technical Details

### Build Method
- **Cross-compilation**: Using Go's cross-compilation features (`CGO_ENABLED=0`)
- **Dependency-free deployment**: All binaries are statically compiled, no extra dependencies needed
- **Optimized size**: Using UPX compression (optional) and build optimizations

### System Requirements
- **Building**: Go 1.24+
- **Running**: Corresponding operating system and architecture

### Encoding Support
- **Unix/Linux/macOS/*BSD**: Native UTF-8 support
- **Windows**: Auto-detect console encoding (GBK/Shift-JIS, etc.) and convert to UTF-8

### Terminal Support
- **Unix-like systems**: Uses PTY (pseudo-terminal), supports full terminal features
- **Windows**: Uses ConPTY (Windows 10+) or cmd.exe pipe, supports basic terminal features

## 🎯 Features

### Dual-End Sync
- ✅ Real-time display of all output on local terminal
- ✅ Direct command input on local terminal
- ✅ Web interface syncs all content
- ✅ Web interface can also input commands
- ✅ Both ends fully synchronized, no interference

### Use Cases
1. **Local Work**: Work on local terminal while presenting on web
2. **Remote Collaboration**: Operate on web while monitoring on local terminal
3. **Teaching Demo**: Operate while showing others
4. **Tech Support**: Local viewing with remote operation
5. **Multi-Device Access**: Multiple web clients connecting to same terminal

## 🐳 Docker Support

Client binaries for all platforms are packaged in the Docker image and can be obtained by:

```bash
# Start Docker service
docker-compose up -d

# Copy all platform binaries from container
docker cp cligool-relay:/app/web/downloads ./
```

## 📝 Platform-Specific Notes

### Windows Platform
- Distributed as ZIP compressed package, no installer program needed
- Supports ConPTY (Windows 10 1809+)
- Auto-detects and converts console encoding
- Dependencies: None (pure Go compilation)

### Linux Platform
- Supports all mainstream Linux distributions
- Includes all standard architectures (x86, ARM, PowerPC, RISC-V, MIPS, etc.)
- Supports LoongArch processors

### *BSD Platform
- Supports FreeBSD, OpenBSD, NetBSD, DragonFlyBSD
- Covers all supported architectures
- Full PTY support

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details
