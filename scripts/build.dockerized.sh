#!/bin/bash

set -eu

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/.."

docker build -t dpb587/gget:build build/docker/build
docker run --rm \
  --volume=$PWD:/root \
  --workdir=/root \
  dpb587/gget:build \
  ./scripts/build.sh
