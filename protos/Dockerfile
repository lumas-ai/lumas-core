FROM golang:alpine

RUN apk update && apk add protobuf git && \
    go get google.golang.org/grpc && \
    go get github.com/golang/protobuf/protoc-gen-go && \
    git clone https://github.com/protocolbuffers/protobuf /google-protos && \
    ln -s /protos/google /google-protos/src/google

COPY build.sh /

CMD /build.sh
