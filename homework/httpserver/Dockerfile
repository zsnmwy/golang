FROM golang:1.16-alpine AS build
LABEL maintainer="zsnmwy szlszl35622@gmail.com"
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY *.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /httpserver




FROM tini/mini

## ps -ef
#PID   USER     TIME   COMMAND
#    1 root       0:00 /sbin/tini -- /docker-entrypoint.sh /httpserver
#    7 httpserv   0:00 /httpserver
#   70 root       0:00 sh
#   77 root       0:00 ps -ef
#
# tini 避免僵尸进程， su-exec 非root执行程序

ENV USERNAME httpserver
ENV USERID 1000

COPY --from=build /httpserver /
COPY docker-entrypoint.sh /
RUN chmod a+x /docker-entrypoint.sh
RUN chmod a+x /httpserver

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories  \
    && apk upgrade --update \
    && apk add -U --no-cache tzdata curl tini\
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata \
    && rm -rf \
    /var/cache/apk/*

ENTRYPOINT ["/sbin/tini", "--", "/docker-entrypoint.sh"]
CMD ["/httpserver"]

HEALTHCHECK --start-period=5s --timeout=3s \
  CMD curl -f http://localhost/healthz || exit 1
