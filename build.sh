#!/bin/sh

IMAGE_NAME="$1"
IMAGE_TAG="$2"

docker build -t ${IMAGE_NAME}:${IMAGE_TAG} -f build/Dockerfile .
