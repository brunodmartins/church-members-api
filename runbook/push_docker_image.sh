#!/bin/bash

set -Eeuo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
REPO_ROOT=$(cd -- "${SCRIPT_DIR}/.." && pwd)

info() { printf '[INFO] %s\n' "$*"; }
error() { printf '[ERROR] %s\n' "$*" >&2; }

trap 'error "Failed at line ${BASH_LINENO[0]}: ${BASH_COMMAND}"' ERR

cd "${REPO_ROOT}"

account_id=$(aws sts get-caller-identity --query Account --output text)
image_tag=$(<.version)
region=$(aws configure get region)
repository="church-members-api-container"
registry="${account_id}.dkr.ecr.${region}.amazonaws.com"
image="${registry}/${repository}:${image_tag}"

info "Checking ECR for ${image}"

set +e
ecr_result=$(aws ecr describe-images \
  --repository-name "${repository}" \
  --image-ids "imageTag=${image_tag}" \
  --query 'imageDetails[0].imageDigest' \
  --output text 2>&1)
ecr_exit_code=$?
set -e

if [[ ${ecr_exit_code} -eq 0 && -n "${ecr_result}" && "${ecr_result}" != "None" ]]; then
  info "Version ${image_tag} already exists in ECR (${ecr_result}); nothing to build or push"
  exit 0
fi

if [[ ${ecr_exit_code} -ne 0 && "${ecr_result}" != *"ImageNotFoundException"* ]]; then
  error "Could not check ECR: ${ecr_result}"
  exit "${ecr_exit_code}"
fi

info "Version ${image_tag} is not in ECR; authenticating and preparing the image"
aws ecr get-login-password --region "${region}" | docker login --username AWS --password-stdin "${registry}"

if [[ -z "$(docker images -q "${image}" 2>/dev/null)" ]]; then
  info "Building ${image}"
  docker buildx build \
    --platform linux/amd64 \
    --load \
    --tag "${image}" .
else
  info "Using existing local image ${image}"
fi

info "Pushing ${image}"
docker push "${image}"
info "Image ${image} pushed successfully"