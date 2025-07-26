# 多阶段构建 - 构建阶段
FROM golang:1.22.5-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chat-mvp main.go

# 运行阶段
FROM alpine:latest

# 安装 ca-certificates 用于 HTTPS 请求
RUN apk --no-cache add ca-certificates

# 创建非 root 用户
RUN adduser -D -s /bin/sh appuser

# 设置工作目录
WORKDIR /home/appuser

# 从构建阶段复制二进制文件和静态资源
COPY --from=builder /app/chat-mvp .
COPY --from=builder /app/static ./static

# 更改文件所有者
RUN chown -R appuser:appuser /home/appuser

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# 启动应用
CMD ["./chat-mvp"]
