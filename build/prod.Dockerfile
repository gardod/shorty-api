FROM golang:1.13.5-alpine3.11 AS build-env
WORKDIR /go/src/github.com/gardod/shorty-api
COPY . /go/src/github.com/gardod/shorty-api
RUN go build -o /server

FROM alpine:3.11
WORKDIR /opt/shorty-api
COPY --from=build-env /server .
ENTRYPOINT ["./server"]
