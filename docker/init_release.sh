#!/bin/sh

set -x

#TODO Use mirroring when packaging locally
if [[ $CN == "yes" ]]; then
  #sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
  sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
  export GOPROXY=https://goproxy.cn
fi


# Docker starts using kernel 4.19 or above
apk add --no-cache ca-certificates bash iproute2 tzdata iptables

# alpine:3.19 compatible with older versions of iptables
apk add --no-cache iptables-legacy
rm /sbin/iptables
ln -s /sbin/iptables-legacy /sbin/iptables


chmod +x /app/docker_entrypoint.sh
mkdir /app/log

#Backup configuration files
cp -r /app/conf /home/conf-bak

tree /app

uname -a
date -Iseconds
