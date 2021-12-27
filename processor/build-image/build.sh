#/usr/bin/env bash
docker build . -t lumas/lumas-processor-build-image

if [ $? -eq 0 ]; then
  echo ""
  echo "Done building lumas/lumas-processor-build-image"
  echo 'Push to Docker Hub with `docker push lumas/lumas-processor-build-image`'
fi
