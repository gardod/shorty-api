FROM golang:1.14-alpine3.11
ENV CGO_ENABLED 0
EXPOSE 40000

RUN apk add --no-cache git
RUN go get github.com/derekparker/delve/cmd/dlv

COPY . /go/src/github.com/gardod/shorty-api
COPY ./config/config.dev.yaml /etc/shorty/config.yaml
COPY ./internal/driver/postgres/migrations /opt/shorty-api/migrations

ENTRYPOINT ["dlv", "debug", "--headless", "--accept-multiclient", "--continue", \
    "--listen=:40000", "--api-version=2", "--log", \
    "github.com/gardod/shorty-api", "--"]
