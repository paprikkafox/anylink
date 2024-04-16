#!/bin/bash

ver=$(cat version)
echo $ver

# docker login -u paprikkafox

# Generation time 2024-01-30T21:41:27+08:00
# date -Iseconds

#bash ./build_web.sh

# docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 本地不生成镜像
docker build -t paprikkafox/anylink:latest --no-cache --progress=plain \
  --build-arg CN="yes" --build-arg appVer=$ver --build-arg commitId=$(git rev-parse HEAD) \
  -f docker/Dockerfile .

echo "docker tag latest $ver"
docker tag paprikkafox/anylink:latest paprikkafox/anylink:$ver

exit 0

docker push paprikkafox/anylink:latest
