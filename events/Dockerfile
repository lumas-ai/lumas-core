FROM golang:alpine as app-build

RUN apk add --update git pkgconfig g++ gcc libc-dev \
    openssl && \
    go get -u github.com/golang/protobuf/proto && \
    go get -u github.com/golang/protobuf/ptypes/struct && \
    go get -u google.golang.org/grpc
