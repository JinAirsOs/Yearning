FROM repositories.nioint.com/adas/ado/golang:1.18.3-buster AS build

WORKDIR /

ENV GIN_MODE=release
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

COPY . .
RUN go build -v -o /opt/Yearning main.go

FROM alpine:3.15

LABEL maintainer="panzhendong <zhendong.pan@nio.com>"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update  \
    && apk add --no-cache ca-certificates bash tree tzdata libc6-compat dumb-init \
    && cp -rf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

COPY --from=build /opt/Yearning /opt/Yearning
COPY conf.toml.template /opt/conf.toml

WORKDIR /opt/

EXPOSE 8000

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/opt/Yearning", "run"]
