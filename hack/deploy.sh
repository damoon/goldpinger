#!/bin/bash

set -euxo pipefail

if [ ! -v "IMAGE_PUSH_PREFIX" ]; then
    IMAGE_PUSH_PREFIX="localhost:5000/goldpinger/goldpinger"
fi
if [ ! -v "IMAGE_PULL_PREFIX" ]; then
    IMAGE_PULL_PREFIX="registry.registry.svc/goldpinger/goldpinger"
fi

docker build .
HASH=$(docker build -q . | cut -d':' -f 2 )
IMAGE_PUSH="$IMAGE_PUSH_PREFIX:$HASH"
IMAGE_PULL="$IMAGE_PULL_PREFIX:$HASH"
docker tag "$HASH" "$IMAGE_PUSH"
docker push "$IMAGE_PUSH"
IMAGE="$IMAGE_PULL" DOLLAR="$" envsubst < kubernetes.yaml | kubectl apply -f -
