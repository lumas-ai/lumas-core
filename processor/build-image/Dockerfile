FROM lumas/lumas-base-opencv

ENV LD_LIBRARY_PATH=/opt/ffmpeg/lib:/usr/local/lib64
ENV PKG_CONFIG_PATH=/opt/ffmpeg/lib/pkgconfig:/usr/local/lib64/pkgconfig
ENV CGO_CPPFLAGS -I/usr/local/include
ENV CGO_CXXFLAGS "--std=c++1z"
ENV CGO_LDFLAGS "-L/usr/local/lib -lopencv_core -lopencv_face -lopencv_videoio -lopencv_imgproc -lopencv_highgui -lopencv_imgcodecs -lopencv_objdetect -lopencv_features2d -lopencv_video -lopencv_dnn -lopencv_xfeatures2d -lopencv_plot -lopencv_tracking"

RUN apk add --update git pkgconfig g++ gcc libc-dev \
    openssl lame libogg libvpx libvorbis libass \
    freetype libtheora opus libwebp x264 x264-libs x265 && \
    go get -u github.com/3d0c/gmf && \
    go get -u github.com/golang/protobuf/proto && \
    go get -u github.com/golang/protobuf/ptypes/struct && \
    go get -u google.golang.org/grpc  && \
    go get -u github.com/lumas-ai/lumas-core/processor && \
    go get -u -d gocv.io/x/gocv && \
    cd / && go build /go/src/github.com/lumas-ai/lumas-core/processor/cmd/server/processor-server.go
