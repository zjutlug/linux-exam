# Build stage
FROM golang:1.24-bookworm AS builder

WORKDIR /build
COPY exam-cli .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o exam-cli

# Final stage
FROM ubuntu:20.04

RUN apt-get update && \
    apt install -y software-properties-common && \
    add-apt-repository -y ppa:neurobin/ppa && \
    apt install -y bash uuid-runtime shc gcc openssh-server vim nano screen tmux curl pv && \
    mkdir /run/sshd && \
    rm -rf /var/lib/apt/lists/*


RUN useradd -m exam && \
    echo "exam:exam2025" | chpasswd

COPY --from=builder /build/exam-cli /usr/local/bin/exam-cli
COPY init.sh /usr/local/bin/init.sh
COPY /challenge /home/exam

RUN chmod +x /usr/local/bin/exam-cli && \
    chmod +x /usr/local/bin/init.sh && \
    chown -R exam: /home/exam && \
    chown exam: /usr/bin/vim && \
    chown exam: /usr/bin/nano

RUN sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config && \
    echo "AllowUsers exam" >> /etc/ssh/sshd_config && \
    echo "ForceCommand HOME=/home/exam /usr/local/bin/exam-cli login" >> /etc/ssh/sshd_config





EXPOSE 22

ENTRYPOINT ["bash","/usr/local/bin/init.sh"]
