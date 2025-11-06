# Build stage
FROM golang:1.24-bookworm AS builder

WORKDIR /build
# 拷贝源代码
COPY exam-cli .

# 下载依赖并构建
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o exam-cli

# Final stage
FROM ubuntu:20.04

RUN apt-get update && \
    apt install -y software-properties-common && \
    add-apt-repository -y ppa:neurobin/ppa && \
    apt install -y bash uuid-runtime shc gcc openssh-server vim nano screen tmux curl && \
    mkdir /run/sshd && \
    rm -rf /var/lib/apt/lists/*

# 创建必要的目录
RUN mkdir -p /challenge

WORKDIR /challenge

# 从构建阶段拷贝二进制文件和初始化脚本
COPY --from=builder /build/exam-cli /usr/local/bin/exam-cli
COPY init.sh /usr/local/bin/init.sh
COPY /challenge /home/exam
# 确保文件有执行权限
RUN chmod +x /usr/local/bin/exam-cli && \
    chmod +x /usr/local/bin/init.sh

# 配置SSH允许密码登录
RUN sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config && \
    echo "AllowUsers exam" >> /etc/ssh/sshd_config && \
    echo "ForceCommand HOME=/home/exam /usr/local/bin/exam-cli login" >> /etc/ssh/sshd_config

# 创建exam用户（使用默认shell）
RUN useradd -m exam && \
    echo "exam:exam2025" | chpasswd


# 暴露SSH端口
EXPOSE 22

ENTRYPOINT ["bash","/usr/local/bin/init.sh"]
