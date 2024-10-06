FROM golang:1.23.2-alpine3.20

RUN echo "https://uk.alpinelinux.org/alpine/v3.20/main/" > /etc/apk/repositories && \
    apk add --no-cache git

RUN go install github.com/go-delve/delve/cmd/dlv@v1.20.1