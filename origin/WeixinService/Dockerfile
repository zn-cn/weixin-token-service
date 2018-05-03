FROM molscar/go1.10_echo:v1
LABEL maintainer="yun_tofar@qq.com"
LABEL version="1.0"
LABEL description="weixin service web"

WORKDIR /app
COPY ./src src
ENV GOPATH "/app"
RUN go build -o bin/main src/main.go && \
    rm -r src
ENTRYPOINT [ "bin/main" ]
