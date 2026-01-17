#!/bin/bash

account_id=$1
image_tag=$2
region=$3

cd ../
aws ecr get-login-password --region $region | docker login --username AWS --password-stdin $account_id.dkr.ecr.$region.amazonaws.com
image=$account_id.dkr.ecr.$region.amazonaws.com/church-members-api-container:$image_tag

if test -z "$(docker images -q $image)"; then
  docker buildx build \
  --platform linux/amd64 \
  --load \
  -t $image .
fi
docker push $image