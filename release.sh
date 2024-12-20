#!/bin/bash

#github action release.sh

set -x
function RETVAL() {
  rt=$1
  if [ $rt != 0 ]; then
    echo $rt
    exit 1
  fi
}

#current directory
cpath=$(pwd)

ver=$(cat version)
echo "当前版本 $ver"

rm -rf artifact-dist
mkdir artifact-dist

function archive() {
  arch=$1
  #echo "Organize deployment files $arch"
  arch_name=${arch//\//-}
  echo $arch_name

  deploy="anylink-$ver-$arch_name"
  docker container rm $deploy
  docker container create --platform $arch --name $deploy paprikkafox/anylink:$ver
  rm -rf anylink-deploy
  docker cp -a $deploy:/app ./anylink-deploy

  ls -lh anylink-deploy

  tar zcf ${deploy}.tar.gz anylink-deploy
  mv ${deploy}.tar.gz artifact-dist/
}

echo "copy二进制文件"

archive "linux/amd64"
archive "linux/arm64"

ls -lh artifact-dist

#Note to run with root privileges
#cd anylink-deploy
#sudo ./anylink --conf="conf/server.toml"
