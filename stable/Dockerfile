FROM golang:1.10

WORKDIR /app
COPY ./src src
ENV GOPATH "/app"
RUN go build -o bin/main src/main.go && \
    rm -r src
ENTRYPOINT [ "bin/main" ]
