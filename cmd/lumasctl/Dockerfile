FROM golang:alpine as build-stage

COPY ./ ./

RUN apk add --update git && \
    go get -d gocv.io/x/gocv && \
    go get github.com/lumas-ai/lumas-core/protos/golang && \
    GO111MODULE=on go get github.com/urfave/cli/v2 && \
    GO111MODULE=on go get gopkg.in/yaml.v2 && \
    go get github.com/mitchellh/go-homedir && \
    CGO_ENABLED=0 GOARCH=amd64 GOOS=linux GO11MODULE=on go build -o /lumasctl ./main.go

FROM scratch AS final

COPY --from=build-stage /lumasctl /lumasctl

ENTRYPOINT ["/lumasctl"]
