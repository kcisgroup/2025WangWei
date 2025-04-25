#!/bin/bash

# 先停docker-compose.yaml
docker-compose -f ./docker-compose-cli.yaml down -v

# 查询链码容器
contain_ids=$(docker ps -a |awk '($2 ~ /dev-*/){print $1}')
# 先停掉链码容器
docker stop $contain_ids
# 然后删除链码容器
docker rm $contain_ids

# 删除链码镜像
images_ids=$(docker images |awk '($1 ~ /dev-*/){print $3}')
docker rmi $images_ids

# 删除生成的文件
rm -rf accountchannel.block channel-artifacts/ crypto-config/ masterchannel.block
