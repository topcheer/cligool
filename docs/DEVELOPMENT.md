# CliGool 开发指南

## 🛠️ 开发环境设置

### 前置要求

- Go 1.21+
- Docker & Docker Compose
- Git
- (可选) Node.js 18+ (Web界面开发)

### 环境配置

1. **克隆仓库**
   ```bash
   git clone https://github.com/cligool/cligool.git
   cd cligool
   ```

2. **安装Go依赖**
   ```bash
   go mod download
   go mod tidy
   ```

3. **配置环境变量**
   ```bash
   cp .env.example .env
   # 编辑 .env 文件，设置必要的配置
   ```

4. **启动开发环境**
   ```bash
   make dev
   ```

## 📁 项目结构

```
cligool/
├── cmd/                    # 命令行工具
│   ├── relay/             # 中继服务器入口
│   └── client/            # CLI客户端入口
├── internal/              # 内部包
│   ├── relay/             # 中继服务逻辑
│   ├── client/            # 客户端逻辑
│   ├── terminal/          # 终端处理
│   ├── auth/              # 认证授权
│   └── database/          # 数据库操作
├── web/                   # Web界面
│   ├── dist/             # 构建产物
│   └── index.html        # 主页面
├── deployments/           # 部署配置
│   └── nginx.conf        # Nginx配置
├── scripts/              # 脚本文件
│   ├── setup.sh         # 快速设置
│   └── build.sh         # 构建脚本
├── docs/                # 文档
├── Makefile             # 构建命令
└── docker-compose.yml   # Docker配置
```

## 🧩 核心组件

### 1. 中继服务 (Relay Service)

**位置**: `internal/relay/`

**功能**:
- WebSocket连接管理
- 消息路由和转发
- 会话管理
- 权限控制

**关键文件**:
- `relay.go` - 核心中继逻辑
- `session.go` - 会话管理
- `websocket.go` - WebSocket处理

### 2. CLI客户端 (Client)

**位置**: `internal/client/`

**功能**:
- 本地终端包装
- WebSocket通信
- PTY管理
- 终端I/O转发

**关键文件**:
- `shell.go` - Shell命令创建
- `connection.go` - 连接管理

### 3. 终端处理 (Terminal)

**位置**: `internal/terminal/`

**功能**:
- 终端参数处理
- Shell配置
- 终端大小管理

### 4. 数据库 (Database)

**位置**: `internal/database/`

**功能**:
- PostgreSQL操作
- Redis缓存
- 数据模型定义

## 🔄 开发工作流

### 1. 功能开发流程

```bash
# 1. 创建功能分支
git checkout -b feature/your-feature-name

# 2. 进行开发
# ... 编写代码 ...

# 3. 本地测试
make test
make run-relay

# 4. 构建验证
make build

# 5. 提交代码
git add .
git commit -m "feat: add your feature description"
git push origin feature/your-feature-name

# 6. 创建Pull Request
```

### 2. 代码规范

**Go代码规范**:
- 遵循 `gofmt` 格式化
- 使用 `golint` 检查代码
- 添加必要的注释和文档
- 编写单元测试

**提交信息规范**:
```
<type>(<scope>): <subject>

<body>

<footer>
```

类型 (type):
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构代码
- `test`: 测试相关
- `chore`: 构建/工具相关

### 3. 测试

**单元测试**:
```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./internal/relay

# 查看覆盖率
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**集成测试**:
```bash
# 启动测试环境
docker-compose -f docker-compose.test.yml up -d

# 运行集成测试
go test -tags=integration ./...

# 清理测试环境
docker-compose -f docker-compose.test.yml down
```

## 🚀 构建和部署

### 本地构建

```bash
# 构建所有组件
make build

# 构建特定组件
make build-relay
make build-client

# 跨平台构建
./scripts/build.sh --all-platforms
```

### Docker构建

```bash
# 构建Docker镜像
make docker-build

# 启动服务
make docker-up

# 查看日志
make docker-logs
```

### 生产部署

1. **环境配置**
   ```bash
   # 编辑生产环境变量
   nano .env.production
   ```

2. **构建镜像**
   ```bash
   docker build -t cligool:latest .
   ```

3. **部署服务**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

## 🐛 调试技巧

### 日志调试

```go
// 使用结构化日志
log.Printf("Debug info: %+v", data)

// 设置日志级别
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

### 性能分析

```bash
# CPU性能分析
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...

# 分析结果
go tool pprof cpu.prof
go tool pprof mem.prof
```

### WebSocket调试

```javascript
// 浏览器控制台
const ws = new WebSocket('ws://localhost:8080/api/terminal/test');
ws.onmessage = (event) => console.log('Received:', event.data);
```

## 📊 监控和指标

### 健康检查

```bash
# 检查服务状态
curl http://localhost:8080/api/health
```

### 性能指标

```bash
# 查看连接数
curl http://localhost:8080/api/metrics

# 查看会话统计
curl http://localhost:8080/api/sessions
```

## 🔧 常见开发任务

### 添加新的API端点

1. 在 `cmd/relay/main.go` 中添加路由
2. 在 `internal/relay/` 中实现处理函数
3. 添加数据库操作（如需要）
4. 编写测试
5. 更新API文档

### 添加WebSocket消息类型

1. 在 `relay.go` 中定义新的消息类型
2. 实现消息处理逻辑
3. 更新客户端消息处理
4. 添加错误处理
5. 测试消息流程

### 修改数据库Schema

1. 在 `database.go` 中更新迁移脚本
2. 更新数据模型
3. 修改相关查询函数
4. 测试数据库迁移
5. 更新API文档

## 📝 文档维护

### 更新API文档

```bash
# 编辑API文档
nano docs/API.md

# 更新使用示例
nano docs/USAGE.md
```

### 更新配置文档

```bash
# 编辑配置文档
nano docs/CONFIG.md

# 更新环境变量说明
nano .env.example
```

## 🤝 贡献指南

1. **Fork仓库**
2. **创建功能分支**
3. **编写代码和测试**
4. **确保代码通过所有测试**
5. **提交Pull Request**

### Pull Request检查清单

- [ ] 代码符合项目规范
- [ ] 添加了必要的测试
- [ ] 所有测试通过
- [ ] 更新了相关文档
- [ ] 提交信息清晰明确

## 🔗 相关资源

- [Go文档](https://golang.org/doc/)
- [WebSocket协议](https://tools.ietf.org/html/rfc6455)
- [PTY编程](https://github.com/creack/pty)
- [xterm.js文档](https://xtermjs.org/)

## 📧 联系方式

- **Issues**: [GitHub Issues](https://github.com/cligool/cligool/issues)
- **Discussions**: [GitHub Discussions](https://github.com/cligool/cligool/discussions)
- **Email**: support@cligool.com