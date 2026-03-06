# 中继服务器Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git make zip

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建中继服务器
RUN CGO_ENABLED=0 GOOS=linux go build -o relay-server ./cmd/relay

# 构建所有平台的客户端
RUN mkdir -p web/downloads

# Windows版本
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o web/downloads/cligool-windows-amd64.exe ./cmd/client && \
    cd web/downloads && zip cligool-windows-amd64.zip cligool-windows-amd64.exe && rm -f cligool-windows-amd64.exe && cd /app

RUN CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o web/downloads/cligool-windows-arm64.exe ./cmd/client && \
    cd web/downloads && zip cligool-windows-arm64.zip cligool-windows-arm64.exe && rm -f cligool-windows-arm64.exe && cd /app

# ==================== Windows版本 ====================
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o web/downloads/cligool-windows-amd64.exe ./cmd/client && \
    cd web/downloads && zip cligool-windows-amd64.zip cligool-windows-amd64.exe && rm -f cligool-windows-amd64.exe && cd /app

RUN CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o web/downloads/cligool-windows-arm64.exe ./cmd/client && \
    cd web/downloads && zip cligool-windows-arm64.zip cligool-windows-arm64.exe && rm -f cligool-windows-arm64.exe && cd /app

# ==================== Linux版本 ====================
# Linux amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web/downloads/cligool-linux-amd64 ./cmd/client

# Linux arm64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o web/downloads/cligool-linux-arm64 ./cmd/client

# Linux 386 (32位)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o web/downloads/cligool-linux-386 ./cmd/client

# Linux arm (32位)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o web/downloads/cligool-linux-arm ./cmd/client

# Linux ppc64le (PowerPC)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -o web/downloads/cligool-linux-ppc64le ./cmd/client

# Linux riscv64 (RISC-V)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build -o web/downloads/cligool-linux-riscv64 ./cmd/client

# Linux s390x (IBM System z)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build -o web/downloads/cligool-linux-s390x ./cmd/client

# Linux mips64le (MIPS Little-Endian)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -o web/downloads/cligool-linux-mips64le ./cmd/client

# ==================== macOS版本 ====================
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o web/downloads/cligool-darwin-amd64 ./cmd/client

RUN CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o web/downloads/cligool-darwin-arm64 ./cmd/client

# ==================== FreeBSD版本 ====================
RUN CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -o web/downloads/cligool-freebsd-amd64 ./cmd/client

RUN CGO_ENABLED=0 GOOS=freebsd GOARCH=arm64 go build -o web/downloads/cligool-freebsd-arm64 ./cmd/client

# ==================== OpenBSD版本 ====================
RUN CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build -o web/downloads/cligool-openbsd-amd64 ./cmd/client

RUN CGO_ENABLED=0 GOOS=openbsd GOARCH=arm64 go build -o web/downloads/cligool-openbsd-arm64 ./cmd/client

# ==================== NetBSD版本 ====================
RUN CGO_ENABLED=0 GOOS=netbsd GOARCH=amd64 go build -o web/downloads/cligool-netbsd-amd64 ./cmd/client

# ==================== DragonFlyBSD版本 ====================
RUN CGO_ENABLED=0 GOOS=dragonfly GOARCH=amd64 go build -o web/downloads/cligool-dragonfly-amd64 ./cmd/client

# 最终镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 复制构建的二进制文件
COPY --from=builder /app/relay-server .

# 复制web界面文件
COPY --from=builder /app/web ./web

# 创建必要的目录
RUN mkdir -p logs

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# 运行服务
CMD ["./relay-server"]