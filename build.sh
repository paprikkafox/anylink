#!/bin/bash

#current directory
cpath=$(pwd)

ver=$(cat version)
echo $ver

#Front-end compilation only needs to be executed once
#bash ./build_web.sh

bash build_docker.sh

deploy="anylink-deploy-$ver"
docker container rm $deploy
docker container create --name $deploy paprikkafox/anylink:$ver
rm -rf anylink-deploy anylink-deploy.tar.gz
docker cp -a $deploy:/app ./anylink-deploy
tar zcf ${deploy}.tar.gz anylink-deploy


./anylink-deploy/anylink -v


echo "anylink Compilation completed, directory: anylink-deploy"
ls -lh anylink-deploy


