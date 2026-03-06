.PHONY: all build clean test deps docker docker-build docker-up docker-down docker-logs help

# 默认目标
all: build

# 安装依赖
deps:
	@echo "安装依赖..."
	go mod download
	go mod tidy

# 构建所有组件
build: deps
	@echo "构建中继服务器..."
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin/relay-server ./cmd/relay
	@echo "构建CLI客户端..."
	CGO_ENABLED=0 go build -o bin/cligool-client ./cmd/client
	@echo "构建完成！"

# 构建中继服务器
build-relay:
	@echo "构建中继服务器..."
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin/relay-server ./cmd/relay

# 构建CLI客户端
build-client:
	@echo "构建CLI客户端..."
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin/cligool ./cmd/client

# 构建跨平台版本
build-all-platforms:
	@echo "构建所有平台版本..."
	@./build-windows.sh

# 构建Windows版本
build-windows:
	@echo "构建Windows版本..."
	mkdir -p bin
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/cligool.exe ./cmd/client
	@echo "Windows版本构建完成: bin/cligool.exe"

# 构建Linux版本
build-linux:
	@echo "构建Linux版本..."
	mkdir -p bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/cligool-linux ./cmd/client
	@echo "Linux版本构建完成: bin/cligool-linux"

# 构建macOS版本
build-macos:
	@echo "构建macOS版本..."
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/cligool-macos ./cmd/client
	@echo "macOS版本构建完成: bin/cligool-macos"

# 运行中继服务器
run-relay: build-relay
	@echo "启动中继服务器..."
	./bin/relay-server

# 运行CLI客户端
run-client: build-client
	@echo "启动CLI客户端..."
	./bin/cligool-client

# 测试
test:
	@echo "运行测试..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -f relay-server cligool-client
	@echo "清理完成！"

# Docker构建
docker-build:
	@echo "构建Docker镜像..."
	docker-compose build

# Docker启动
docker-up:
	@echo "启动Docker服务..."
	docker-compose up -d

# Docker停止
docker-down:
	@echo "停止Docker服务..."
	docker-compose down

# Docker日志
docker-logs:
	@echo "查看Docker日志..."
	docker-compose logs -f

# Docker重启
docker-restart: docker-down docker-up

# 初始化数据库
init-db:
	@echo "初始化数据库..."
	docker-compose up -d postgres redis
	@echo "等待数据库启动..."
	sleep 5
	@echo "数据库初始化完成！"

# 开发环境启动
dev: init-db
	@echo "启动开发环境..."
	docker-compose up

# 生成SSL证书
gen-certs:
	@echo "生成SSL证书..."
	mkdir -p certs
	openssl req -x509 -newkey rsa:4096 -keyout certs/key.pem -out certs/cert.pem -days 365 -nodes -subj "/CN=localhost"
	@echo "证书生成完成！"

# 帮助信息
help:
	@echo "CliGool 构建系统"
	@echo ""
	@echo "使用方法:"
	@echo "  make build              - 构建所有组件"
	@echo "  make build-relay        - 构建中继服务器"
	@echo "  make build-client       - 构建CLI客户端"
	@echo "  make build-all-platforms - 构建所有平台版本"
	@echo "  make build-windows      - 构建Windows版本"
	@echo "  make build-linux        - 构建Linux版本"
	@echo "  make build-macos        - 构建macOS版本"
	@echo "  make run-relay          - 运行中继服务器"
	@echo "  make run-client         - 运行CLI客户端"
	@echo "  make test               - 运行测试"
	@echo "  make clean              - 清理构建文件"
	@echo "  make docker-build       - 构建Docker镜像"
	@echo "  make docker-up          - 启动Docker服务"
	@echo "  make docker-down        - 停止Docker服务"
	@echo "  make dev                - 启动开发环境"
	@echo "  make gen-certs          - 生成SSL证书"
	@echo ""
	@echo "平台支持:"
	@echo "  ✅ Windows   - 完全支持（有限功能）"
	@echo "  ✅ Linux     - 完全支持（所有功能）"
	@echo "  ✅ macOS     - 完全支持（所有功能）"