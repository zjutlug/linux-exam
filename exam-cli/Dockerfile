# Build stage
FROM golang:1.24-bookworm AS builder

WORKDIR /build
# 拷贝源代码
COPY . .

# 下载依赖并构建
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o exam-cli

# Final stage
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y bash && \
    rm -rf /var/lib/apt/lists/*

# 创建必要的目录
RUN mkdir -p /challenge

WORKDIR /challenge

# 从构建阶段拷贝二进制文件
COPY --from=builder /build/exam-cli /usr/local/bin/exam-cli

# 确保二进制文件有执行权限
RUN chmod +x /usr/local/bin/exam-cli

# 设置默认的环境变量
ENV EXAM_API_URL=http://127.0.0.1:8080

ENTRYPOINT ["exam-cli","login"]
