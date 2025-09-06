# 使用 Go 1.25 的 Alpine 基础镜像
FROM golang:1.25-alpine

# [关键改动] 更换 Alpine 的软件源为国内镜像 (此处使用清华大学镜像)
# 以解决国内访问官方源不稳定的问题
# [Key Change] Switch Alpine's software repository to a domestic mirror (Tsinghua University mirror)
# to solve the instability of accessing the official source from within China.
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

# 设置环境变量
ENV GOPROXY=https://goproxy.cn,direct
ENV TZ=Asia/Shanghai
ENV GOEXPERIMENT=jsonv2,greenteagc

# 更新软件包列表、安装时区数据并设置时区，合并为一条 RUN 指令以减小镜像层数
# Update package list, install timezone data, and set the timezone in a single RUN command to reduce image layers
RUN apk update --no-cache && \
    apk add --no-cache tzdata && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

# 安装 modd 用于热重载
RUN go install github.com/cortesi/modd/cmd/modd@latest

# 设置工作目录
WORKDIR /go

# 容器启动时执行的默认命令
CMD ["modd"]