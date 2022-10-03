# 指定基础镜像的版本，确保每次构建都是幂等的
FROM golang:latest AS builder

WORKDIR /app

# Copy go.mod and go.sum first, because of caching reasons.
# COPY go.mod go.sum ./
# RUN go mod download

COPY . ./
# Compile project
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# 使用体积更小的基础镜像
# FROM alpine:3.15 AS production
# Golang 项目推荐 scratch 镜像进一步减小体积
FROM scratch AS final

# 不要使用 root 权限运行应用
# RUN adduser -D -u 10000 florin
# USER florin

# 设置时区
# 在使用 Docker 容器时，系统默认的时区就是 UTC 时间（0 时区），和我们实际需要的北京时间相差八个小时
# ENV LANG=en_US.UTF-8 LANGUAGE=en_US:en LC_ALL=en_US.UTF-8 TZ=Asia/Shanghai
# RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=builder /app/main .

# 默认暴露 80 端口
EXPOSE 8080

CMD ["./main"]
