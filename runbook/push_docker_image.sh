#!/bin/bash

account_id=$1
image_tag=$2
region=$3

cd ../../../
aws ecr get-login-password --region $region | docker login --username AWS --password-stdin $account_id.dkr.ecr.$region.amazonaws.com

if test -z "$(docker images -q $image_tag)"; then
  docker build -t $image_tag .
fi
docker push $image_tag