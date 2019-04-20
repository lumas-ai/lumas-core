FROM golang:alpine

RUN apk add --update git && \
    go get -u github.com/golang/protobuf/proto && \
    go get -u github.com/golang/protobuf/ptypes/struct && \
    go get -u github.com/gomodule/redigo/redis && \
    go get -u github.com/nitishm/go-rejson && \
    go get -u github.com/go-redis/redis && \
    go get -u google.golang.org/grpc
