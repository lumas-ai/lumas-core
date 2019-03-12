FROM golang:alpine as build-stage

RUN apk add --update git && \
    go get -d gocv.io/x/gocv && \
    go get github.com/lumas-ai/protos/golang && \
    go get gopkg.in/urfave/cli.v1 && \
    cd / && go build /go/src/github.com/lumas-ai/lumas-core/cmd/lumasctl/lumasctl.go
