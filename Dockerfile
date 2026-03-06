# 中继服务器Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git make

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建中继服务器
RUN CGO_ENABLED=0 GOOS=linux go build -o relay-server ./cmd/relay

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