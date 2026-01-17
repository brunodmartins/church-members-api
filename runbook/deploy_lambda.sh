#!/bin/bash

if [ $# -ne 2 ]; then
    echo "Usage: $0 <lambda_name> <image_tag>"
    exit 1
fi

LAMBDA_NAME="$1"
IMAGE_TAG="$2"

ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
REGION=$(aws configure get region)

REPO="church-members-api-container"
IMAGE_URI_TAG="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/${REPO}:${IMAGE_TAG}"

# Try to fetch the image manifest and select the digest for platform amd64/linux.
# Falls back to the repository digest (may be a manifest-list digest) and finally
# to the tag if we can't obtain a digest.
echo "Resolving image digest for ${IMAGE_URI_TAG} (preferring amd64/linux)..."

DIGEST=""

# Prefer using docker's remote manifest inspect (requires docker and access to ECR).
# This gives a manifest-list we can query for the amd64/linux digest.
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
    echo "Using image digest: ${DIGEST}"
else
    IMAGE_URI="${IMAGE_URI_TAG}"
    echo "Could not resolve platform-specific digest; falling back to tag: ${IMAGE_URI_TAG}"
fi

echo "Updating function ${LAMBDA_NAME} with image ${IMAGE_URI}..."

VERSION=$(aws lambda update-function-code \
    --function-name "${LAMBDA_NAME}" \
    --image-uri "${IMAGE_URI}" \
    --publish \
    --query 'Version' \
    --output text)

if [ $? -ne 0 ]; then
    echo "❌ Error updating Lambda"
    exit 1
fi

echo "📦 New version published: $VERSION"
echo "🔄 Waiting for Lambda to be ready..."

# Function to check Lambda status
check_status() {
    aws lambda get-function --function-name "${LAMBDA_NAME}" \
        --query 'Configuration.State' \
        --output text
}

# Function to check last update status
check_last_update() {
    aws lambda get-function --function-name "${LAMBDA_NAME}" \
        --query 'Configuration.LastUpdateStatus' \
        --output text
}

# Wait until Lambda is active and last update is successful
MAX_ATTEMPTS=30
ATTEMPT=1

while [ $ATTEMPT -le $MAX_ATTEMPTS ]; do
    STATUS=$(check_status)
    UPDATE_STATUS=$(check_last_update)

    echo "Current state: $STATUS, Update status: $UPDATE_STATUS"

    if [ "$STATUS" = "Active" ] && [ "$UPDATE_STATUS" = "Successful" ]; then
        echo "✅ Lambda updated and ready!"
        exit 0
    elif [ "$UPDATE_STATUS" = "Failed" ]; then
        echo "❌ Lambda update failed"
        exit 1
    fi

    ATTEMPT=$((ATTEMPT + 1))
    sleep 2
done

echo "❌ Timeout waiting for Lambda update"
exit 1