#/usr/bin/env bash
docker build . -t lumas/lumas-streamserver-build-image

if [ $? -eq 0 ]; then
  echo ""
  echo "Done building lumas/lumas-streamserver-build-image"
  echo 'Push to Docker Hub with `docker push lumas/lumas-streamserver-build-image`'
fi
