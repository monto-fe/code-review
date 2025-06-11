# Build stage
FROM golang:1.19-alpine AS builder

# 设置 Go 代理（国内推荐 goproxy.cn）
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git

# 复制 go mod 文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY backend/ .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Production stage
FROM alpine:latest

WORKDIR /app

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY .env .

# 暴露端口
EXPOSE 9000 

# 启动前延迟 5 秒，防止出现数据库没有初始化完成导致的失败
CMD ["sh", "-c", "sleep 5 && ./main"]
