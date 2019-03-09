FROM lumas/lumas-base-opencv as app-build

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


FROM alpine:3.9

ENV LD_LIBRARY_PATH=/opt/ffmpeg/lib:/usr/local/lib64
ENV PKG_CONFIG_PATH=/opt/ffmpeg/lib/pkgconfig:/usr/local/lib64/pkgconfig

COPY --from=app-build /usr/local/lib64 /usr/local/lib64
COPY --from=app-build /usr/local/lib64/pkgconfig/opencv4.pc /usr/local/lib64/pkgconfig/opencv4.pc
COPY --from=app-build /usr/local/include/opencv4/opencv2 /usr/local/include/opencv4/opencv2
COPY --from=app-build /opt/ffmpeg /opt/ffmpeg
COPY --from=app-build /processor-server /processor-server

RUN apk add --update git pkgconfig \
    openssl lame libogg libvpx libvorbis libass \
    freetype libtheora opus libwebp x264 x264-libs x265 \
    libjpeg-turbo libpng libwebp libwebp-dev tiff \
    libavc1394 jasper-libs openblas openssl libgphoto2 \
    gstreamer gst-plugins-base


ENTRYPOINT ["/processor-server"]
