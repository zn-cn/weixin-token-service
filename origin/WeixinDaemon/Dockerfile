FROM golang:1.10
LABEL maintainer="yun_tofar@qq.com"
LABEL version="1.0"
LABEL description="weixin service daemon"

WORKDIR /app
COPY ./src src
ENV GOPATH "/app"
RUN go get -d ./src/... && \
    go build -o bin/main src/main.go && \
    rm -r src
ENTRYPOINT [ "bin/main" ]
