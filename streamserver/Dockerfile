FROM golang:alpine

RUN apk add --update git && \
    go get -u github.com/golang/protobuf/proto && \
    go get -u github.com/golang/protobuf/ptypes/struct && \
    go get -u google.golang.org/grpc
