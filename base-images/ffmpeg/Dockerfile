FROM golang:alpine

ENV FFMPEG_VERSION=4.1 BUILD_PREFIX=/opt/ffmpeg

RUN apk add --update build-base curl nasm tar bzip2 \
  zlib-dev openssl-dev yasm-dev lame-dev libogg-dev \
  x264-dev libvpx-dev libvorbis-dev x265-dev  \
  freetype-dev libass-dev libwebp-dev libtheora-dev \
  opus-dev && \
  DIR=$(mktemp -d) && cd ${DIR} && \
  curl -s http://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz | tar zxvf - -C . && \
  cd ffmpeg-${FFMPEG_VERSION} && \
  ./configure \
  --enable-shared --enable-version3 --enable-gpl --enable-nonfree --enable-small --enable-libmp3lame --enable-libx264 --enable-libx265 --enable-libvpx --enable-libtheora --enable-libvorbis --enable-libopus --enable-libass --enable-libwebp --enable-postproc --enable-avresample --enable-libfreetype --enable-openssl --disable-debug --prefix=${BUILD_PREFIX} && \
  make && \
  make install && \
  rm -rf ${DIR} && \
  apk del build-base curl tar bzip2 x264 openssl nasm && rm -rf /var/cache/apk/* && \
  ln -s /opt/ffmpeg/bin/ffmpeg /usr/bin/ffmpeg
