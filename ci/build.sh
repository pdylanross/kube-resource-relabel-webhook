#!/usr/bin/env bash

set -e

COMMIT_HASH=$(git rev-parse HEAD)
VERSION=$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*')
BUILD_TIMESTAMP=$(date --rfc-3339=seconds)

docker build . \
  --build-arg "VERSION=$VERSION" \
  --build-arg "COMMIT_HASH=$COMMIT_HASH" \
  --build-arg "BUILD_TIMESTAMP=$BUILD_TIMESTAMP" \
  -t relabel:dev

docker tag relabel:dev ghcr.io/pdylanross/kube-resource-relabel-webhook:$COMMIT_HASH