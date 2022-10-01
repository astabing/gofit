#!/bin/sh
set -e

IMAGE=astabing/gofit
NETWORK=$(docker network ls -f name=gofit -q)
cd $(dirname $0)

docker run --rm -ti --name gofit \
	--network $NETWORK \
	--env-file ${PWD}/.env_gofit \
	-v ${PWD}/.cache:/app/.cache \
	$IMAGE $@

