# CliGool 快速开始指南

## 🚀 5分钟快速部署

### 前置条件
- ✅ Docker和Docker Compose已安装
- ✅ Cloudflare账号（免费版即可）
- ✅ 一个域名（托管在Cloudflare）

### 第一步：部署应用
```bash
# 1. 克隆项目
git clone https://github.com/cligool/cligool.git
cd cligool

# 2. 一键部署
./scripts/deploy.sh
```

### 第二步：配置Cloudflare Tunnel
```bash
# 1. 安装cloudflared（如果没有）
brew install cloudflare/cloudflare/cloudflared  # macOS
# 或
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared-linux-amd64.deb  # Linux

# 2. 运行配置脚本
./scripts/cloudflare-tunnel.sh
```

### 第三步：开始使用
```bash
# 1. 构建客户端
make build-client

# 2. 在目标机器上启动客户端
./bin/cligool-client -server https://your-domain.com -session test

# 3. 在浏览器中访问
# https://your-domain.com/?session=test
```

## 🎯 典型使用场景

### 场景1：远程访问家庭电脑
```bash
# 在家里电脑上
cd cligool
./scripts/deploy.sh          # 启动服务
./scripts/cloudflare-tunnel.sh  # 配置域名

# 在公司/学校的电脑上
./bin/cligool-client -server https://home.yourdomain.com -session home-pc

# 在浏览器中访问
# https://home.yourdomain.com/?session=home-pc
```

### 场景2：团队协作调试
```bash
# 在服务器上
./bin/cligool-client -server https://team.yourdomain.com -session debug-session

# 团队成员访问同一URL
# https://team.yourdomain.com/?session=debug-session
```

### 场景3：教学演示
```bash
# 老师的电脑
./bin/cligool-client -server https://demo.yourdomain.com -session lecture-101

# 学生访问
# https://demo.yourdomain.com/?session=lecture-101
```

## 🔧 常用命令

```bash
# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 重新构建
docker-compose build
```

## 📱 客户端使用

```bash
# 基础连接
./bin/cligool-client -server https://your-domain.com

# 指定会话ID
./bin/cligool-client -server https://your-domain.com -session my-session

# 使用特定Shell
./bin/cligool-client -server https://your-domain.com -shell /bin/zsh

# 调试模式
./bin/cligool-client -server https://your-domain.com -debug
```

## 🌐 Web界面功能

- ✅ 实时终端显示
- ✅ 完整的终端特性（颜色、光标等）
- ✅ 多用户协作
- ✅ 自动终端大小适配
- ✅ 快捷键支持（Ctrl+C复制等）

## 🔒 安全特性

- 🔐 **零暴露攻击面** - 无需开放端口到公网
- 🛡️ **DDoS保护** - Cloudflare防护网络
- 🔑 **自动HTTPS** - 无需配置证书
- 👥 **访问控制** - 可配置Cloudflare Access
- 🌍 **全球CDN** - 低延迟访问

## ❓ 常见问题

**Q: 为什么选择Cloudflare Tunnel？**
A: 零配置HTTPS、自动DDoS保护、无需开放端口、全球CDN加速。

**Q: 可以在没有域名的环境下使用吗？**
A: 需要一个托管在Cloudflare的域名，但可以使用免费的子域名。

**Q: 支持多少并发用户？**
A: 取决于服务器配置，一般可以支持数百个并发连接。

**Q: 数据安全吗？**
A: 所有通信通过Cloudflare的TLS加密，且你的服务器不直接暴露在公网。

## 📚 更多文档

- **完整配置**: [docs/CONFIG.md](docs/CONFIG.md)
- **部署指南**: [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)
- **使用手册**: [docs/USAGE.md](docs/USAGE.md)
- **开发指南**: [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)

## 🎉 开始使用

现在你已经准备好了！选择一个场景开始使用CliGool吧：

1. **个人使用** - 远程访问家里的电脑
2. **团队协作** - 与同事一起调试服务器
3. **教学演示** - 向学生展示命令行操作
4. **技术支持** - 远程协助他人解决电脑问题

---

*有问题？查看[文档](docs/)或提交[Issue](https://github.com/cligool/cligool/issues)*