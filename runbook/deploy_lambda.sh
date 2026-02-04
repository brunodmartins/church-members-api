#!/bin/bash

# Improved logging and error messages
set -u
set -o pipefail

usage() {
    echo "Usage: $0 <lambda_name> <image_tag>"
    exit 1
}

if [ $# -ne 2 ]; then
    usage
fi

LAMBDA_NAME="$1"
IMAGE_TAG="$2"

ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
REGION=$(aws configure get region)

REPO="church-members-api-container"
IMAGE_URI_TAG="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/${REPO}:${IMAGE_TAG}"

timestamp() { date -u +"%Y-%m-%dT%H:%M:%SZ"; }
log() { level="$1"; shift; printf "%s [%s] %s\n" "$(timestamp)" "$level" "$*"; }
info() { log "INFO" "$@"; }
warn() { log "WARN" "$@"; }
error() { log "ERROR" "$@" >&2; }

debug() {
    if [ "${DEBUG:-}" = "1" ]; then
        log "DEBUG" "$@"
    fi
}

info "Resolving image digest for ${IMAGE_URI_TAG} (preferring amd64/linux)..."

DIGEST=""

# Prefer using docker's remote manifest inspect (requires docker and access to ECR).
if command -v docker >/dev/null 2>&1 && command -v jq >/dev/null 2>&1; then
    TMP_MANIFEST=$(mktemp)
    if docker manifest inspect "${IMAGE_URI_TAG}" > "$TMP_MANIFEST" 2>/dev/null; then
        if jq -e 'has("manifests")' "$TMP_MANIFEST" >/dev/null 2>&1; then
            DIGEST=$(jq -r '.manifests[] | select(.platform.architecture=="amd64" and .platform.os=="linux") | .digest' "$TMP_MANIFEST" | head -n1)
            if [ -z "$DIGEST" ] || [ "$DIGEST" = "null" ]; then
                DIGEST=$(jq -r '.manifests[0].digest' "$TMP_MANIFEST" 2>/dev/null || true)
            fi
        else
            DIGEST=$(jq -r '.config.digest // .digest' "$TMP_MANIFEST" 2>/dev/null || true)
        fi
    else
        warn "docker manifest inspect failed for ${IMAGE_URI_TAG}; falling back to AWS lookup"
    fi
    rm -f "$TMP_MANIFEST"
fi

# Fallback to describe-images (AWS) to get the image digest (may be manifest-list digest)
if [ -z "$DIGEST" ]; then
    DIGEST=$(aws ecr describe-images \
        --repository-name "${REPO}" \
        --image-ids imageTag="${IMAGE_TAG}" \
        --query 'imageDetails[0].imageDigest' --output text 2>/dev/null || true)
fi

if [ -n "$DIGEST" ] && [ "$DIGEST" != "None" ]; then
    IMAGE_URI="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/${REPO}@${DIGEST}"
    info "Using image digest: ${DIGEST}"
else
    IMAGE_URI="${IMAGE_URI_TAG}"
    warn "Could not resolve platform-specific digest; falling back to tag: ${IMAGE_URI_TAG}"
fi

info "Updating function ${LAMBDA_NAME} with image ${IMAGE_URI}..."

VERSION=$(aws lambda update-function-code \
    --function-name "${LAMBDA_NAME}" \
    --image-uri "${IMAGE_URI}" \
    --publish \
    --query 'Version' \
    --output text) || { error "Error updating Lambda: aws lambda update-function-code failed"; exit 1; }

info "📦 New version published: $VERSION"
info "🔄 Waiting for Lambda to be ready..."

# Function to check Lambda status
check_status() {
    aws lambda get-function --function-name "${LAMBDA_NAME}" \
        --query 'Configuration.State' --output text 2>/dev/null || echo "Unknown"
}

# Function to check last update status
check_last_update() {
    aws lambda get-function --function-name "${LAMBDA_NAME}" \
        --query 'Configuration.LastUpdateStatus' --output text 2>/dev/null || echo "Unknown"
}

# Wait until Lambda is active and last update is successful
MAX_ATTEMPTS=30
ATTEMPT=1

while [ $ATTEMPT -le $MAX_ATTEMPTS ]; do
    STATUS=$(check_status)
    UPDATE_STATUS=$(check_last_update)

    info "Current state: ${STATUS}, Update status: ${UPDATE_STATUS} (attempt ${ATTEMPT}/${MAX_ATTEMPTS})"
    debug "(raw status: ${STATUS}), (raw update: ${UPDATE_STATUS})"

    if [ "${STATUS}" = "Active" ] && [ "${UPDATE_STATUS}" = "Successful" ]; then
        info "✅ Lambda updated and ready!"
        exit 0
    elif [ "${UPDATE_STATUS}" = "Failed" ]; then
        error "❌ Lambda update failed"
        exit 1
    fi

    ATTEMPT=$((ATTEMPT + 1))
    sleep 2
done

error "❌ Timeout waiting for Lambda update after ${MAX_ATTEMPTS} attempts"
exit 1