#!/bin/bash

set -eu

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/.."

docker build -t dpb587/gget:build build/docker/build
docker run --rm \
  --volume="${PWD}":/root \
  --workdir=/root \
  --user=$( id -u "${USER}" ):$( id -g "${USER}" ) \
  --env=GOCACHE=/tmp/.cache/go-build \
  dpb587/gget:build \
  ./scripts/build.local.sh "$@"
