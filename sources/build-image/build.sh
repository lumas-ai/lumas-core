#/usr/bin/env bash
docker build . -t lumas/lumas-sources-build-image

if [ $? -eq 0 ]; then
  echo ""
  echo "Done building lumas/lumas-sources-build-image"
  echo 'Push to Docker Hub with `docker push lumas/lumas-sources-build-image`'
fi
