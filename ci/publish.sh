#!/usr/bin/env bash

set -e

COMMIT_HASH=$(git rev-parse HEAD)

docker push ghcr.io/pdylanross/kube-resource-relabel-webhook:$COMMIT_HASH