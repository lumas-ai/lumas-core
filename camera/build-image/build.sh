#/usr/bin/env bash
docker build . -t lumas/lumas-camera-build-image

if [ $? -eq 0 ]; then
  echo ""
  echo "Done building lumas/lumas-camera-build-image"
  echo 'Push to Docker Hub with `docker push lumas/lumas-camera-build-image`'
fi
