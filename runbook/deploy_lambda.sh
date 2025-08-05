#!/bin/bash

if [ $# -ne 2 ]; then
    echo "Usage: $0 <lambda_name> <image_tag>"
    exit 1
fi

LAMBDA_NAME="$1"
IMAGE_TAG="$2"

ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
REGION=$(aws configure get region)

IMAGE_URI="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/church-members-api-container:${IMAGE_TAG}"

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
