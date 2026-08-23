#!/bin/bash

cd ../

account_id=$(aws account get-account-information | jq -r ".AccountId")
image_tag=$(cat .version)
region=$(aws configure get region)

aws ecr get-login-password --region $region | docker login --username AWS --password-stdin $account_id.dkr.ecr.$region.amazonaws.com
image=$account_id.dkr.ecr.$region.amazonaws.com/church-members-api-container:$image_tag

if test -z "$(docker images -q $image)"; then
  docker buildx build \
  --platform linux/amd64 \
  --load \
  -t $image .
fi
docker push $image