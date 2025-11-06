#!/bin/bash

# 设置错误时立即退出
set -e

CONTAINER_ID_FILE="/etc/container_id"
LOG_PREFIX="[Container Init]"

# 日志函数
log() {
    echo "${LOG_PREFIX} $1"
}

# 错误处理函数
handle_error() {
    log "Error on line $1: $2"
    exit 1
}

# 系统信息输出函数
print_system_info() {
    log "=================== 系统信息 ==================="
    log "操作系统: $(cat /etc/os-release | grep PRETTY_NAME | cut -d'"' -f2)"
    log "内核版本: $(uname -r)"
    log "架构: $(uname -m)"
    log "主机名: $(hostname)"
    log "容器ID: $(cat $CONTAINER_ID_FILE)"
    log "=============================================="
}

# 环境变量输出函数
print_env_vars() {
    log "=================== 系统环境变量 ==================="
#    env | while read -r line; do
#        log "$line"
#    done
    while read -r line; do
        log "$line"
    done < /etc/environment
    log "=============================================="
}

# 设置错误处理
trap 'handle_error ${LINENO} "$BASH_COMMAND"' ERR

# 传递docker环境变量到系统环境变量
env > /etc/environment

# 生成 container_id
if [ ! -f "$CONTAINER_ID_FILE" ]; then
    log "生成容器ID..."
    uuidgen > "$CONTAINER_ID_FILE"
    chown root:root "$CONTAINER_ID_FILE"
    chmod 444 "$CONTAINER_ID_FILE"
    log "容器ID已生成: $(cat $CONTAINER_ID_FILE)"
fi

# 输出系统信息
print_system_info

# 输出环境变量信息
print_env_vars

# 启动SSH服务
log "启动SSH服务..."
log "SSH配置信息:"
log "- 允许用户: exam"

/usr/sbin/sshd -D &
SSH_PID=$!

# 设置信号处理
trap 'kill $SSH_PID; exit 0' SIGTERM SIGINT

# 等待SSH服务
log "容器初始化完成，等待SSH连接..."
wait $SSH_PID
