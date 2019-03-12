FROM golang:alpine as build-stage

RUN apk add --update git && \
    go get -d gocv.io/x/gocv && \
    go get github.com/lumas-ai/lumas-core/protos/golang && \
    go get github.com/lumas-ai/lumas-core/cmd/lumasctl && \
    go get gopkg.in/urfave/cli.v1 && \
    cd / && go build -o /lumasctl /go/src/github.com/lumas-ai/lumas-core/cmd/lumasctl/main.go

FROM alpine:3.9

COPY --from=build-stage /lumasctl /lumasctl

ENTRYPOINT ["/lumasctl"]
