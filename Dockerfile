# Build stage
FROM golang:1.27.3 AS builder
#ENV GOPROXY=https://proxy.golang.com.cn

WORKDIR /app

COPY . .

RUN bash ./scripts/build.sh

# Production stage
FROM amd64/alpine:latest

# Set System TimeZone
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
COPY --from=builder /app/bin/main .

ENV WORKDIR /app
CMD ./main start
