FROM golang:alpine as build-stage

RUN apk add --update git && \
    go get -d gocv.io/x/gocv && \
    go get github.com/lumas-ai/protos/golang
