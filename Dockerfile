FROM golang:1.23.9 AS builder

WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/http-server ./cmd/http/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/cron-server ./cmd/cron/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/wallet-server ./cmd/wallet/main.go

# 使用更小的基础镜像
FROM alpine:latest

WORKDIR /app

# 安装基本工具
RUN apk --no-cache add ca-certificates tzdata

# 设置时区为上海
ENV TZ=Asia/Shanghai

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/bin/ /app/bin/

# 复制配置文件
COPY --from=builder /app/etc/ /app/etc/

# 创建日志目录
RUN mkdir -p /app/log

# 设置工作目录
WORKDIR /app

# 暴露端口
EXPOSE 2900

# 默认启动HTTP服务
CMD ["/app/bin/http-server", "--conf", "/app/etc/config.yaml"] 