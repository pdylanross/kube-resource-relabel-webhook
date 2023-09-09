#!/usr/bin/env bash

set -e

REPO=ghcr.io/pdylanross/kube-resource-relabel-webhook

COMMIT_HASH=$(git rev-parse HEAD)
TAG=${1:-$COMMIT_HASH}

docker tag $REPO:$COMMIT_HASH $REPO:$TAG
docker push $REPO:$TAG