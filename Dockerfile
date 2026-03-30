# 第一阶段：构建
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates tzdata

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译服务器二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/wordzero-server \
    ./cmd/server

# 第二阶段：运行时镜像
FROM alpine:3.20

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/wordzero-server /app/wordzero-server

# 设置时区为上海
ENV TZ=Asia/Shanghai

# 暴露服务端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动服务
ENTRYPOINT ["/app/wordzero-server"]
