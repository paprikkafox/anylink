#!/bin/bash

ver=$(cat version)
echo $ver

echo "docker tag latest $ver"

docker pull --platform=linux/amd64 paprikkafox/anylink:$ver

docker tag paprikkafox/anylink:$ver registry.cn-hangzhou.aliyuncs.com/paprikkafox/anylink:latest
docker push registry.cn-hangzhou.aliyuncs.com/paprikkafox/anylink:latest

docker tag paprikkafox/anylink:$ver registry.cn-hangzhou.aliyuncs.com/paprikkafox/anylink:$ver
docker push registry.cn-hangzhou.aliyuncs.com/paprikkafox/anylink:$ver

docker rmi paprikkafox/anylink:$ver

#arm64
docker pull --platform=linux/arm64 paprikkafox/anylink:$ver

docker tag paprikkafox/anylink:$ver registry.cn-hangzhou.aliyuncs.com/paprikkafox/anylink:arm64v8-latest
docker push registry.cn-hangzhou.aliyuncs.com/paprikkafox/anylink:arm64v8-latest

docker tag paprikkafox/anylink:$ver registry.cn-hangzhou.aliyuncs.com/paprikkafox/anylink:arm64v8-$ver
docker push registry.cn-hangzhou.aliyuncs.com/paprikkafox/anylink:arm64v8-$ver

docker rmi paprikkafox/anylink:$ver
