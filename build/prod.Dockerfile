FROM golang:1.14-alpine3.11 AS build-env
WORKDIR /go/src/github.com/gardod/shorty-api
COPY . /go/src/github.com/gardod/shorty-api
RUN go build -o /server

FROM alpine:3.11
WORKDIR /opt/shorty-api
COPY ./internal/driver/postgres/migrations ./migrations
COPY --from=build-env /server .
ENTRYPOINT ["./server"]
