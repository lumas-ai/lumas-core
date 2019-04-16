FROM alpine:3.9

ENV LD_LIBRARY_PATH=/opt/ffmpeg/lib:/usr/local/lib64
ENV PKG_CONFIG_PATH=/opt/ffmpeg/lib/pkgconfig:/usr/local/lib64/pkgconfig

COPY --from=lumas/lumas-processor-build-image /usr/local/lib64 /usr/local/lib64
COPY --from=lumas/lumas-processor-build-image /usr/local/lib64/pkgconfig/opencv4.pc /usr/local/lib64/pkgconfig/opencv4.pc
COPY --from=lumas/lumas-processor-build-image /usr/local/include/opencv4/opencv2 /usr/local/include/opencv4/opencv2
COPY --from=lumas/lumas-processor-build-image /opt/ffmpeg /opt/ffmpeg
COPY --from=lumas/lumas-processor-build-image /processor-server /processor-server

RUN apk add --update git pkgconfig \
    openssl lame libogg libvpx libvorbis libass \
    freetype libtheora opus libwebp x264 x264-libs x265 \
    libjpeg-turbo libpng libwebp libwebp-dev tiff \
    libavc1394 jasper-libs openblas openssl libgphoto2 \
    gstreamer gst-plugins-base


ENTRYPOINT ["/processor-server"]
