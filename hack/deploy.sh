#!/bin/bash

HASH=$(docker build -q . | cut -d':' -f 2 )
IMAGE_PUSH="$IMAGE_PUSH_PREFIX:$HASH"
IMAGE_PULL="$IMAGE_PULL_PREFIX:$HASH"
docker tag "$HASH" "$IMAGE_PUSH"
docker push "$IMAGE_PUSH"
cat kubernetes.yaml | IMAGE="$IMAGE_PULL" envsubst | kubectl apply -f -
