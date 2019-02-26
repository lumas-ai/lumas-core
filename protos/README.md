# Compile new protobufs

```
docker run -it -v `pwd`:/protos lumas/protos
```

# Build a new image

```
docker build . -t lumas/protos
```
