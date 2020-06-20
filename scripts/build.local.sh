#!/bin/bash
# args: [version]

set -eu

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/.."

version="${1:-0.0.0}"

if [ -z "${os_list:-}" ]; then
  os_list="darwin linux windows"
fi

if [ -z "${arch_list:-}" ]; then
  arch_list="amd64"
fi

rm -fr tmp/build
mkdir -p tmp/build

commit=$( git rev-parse HEAD | cut -c-10 )$( git diff-index --quiet HEAD -- || echo "+dirty" )
built=$( date -u +%Y-%m-%dT%H:%M:%S+00:00 )

export CGO_ENABLED=0

cli=gget

for os in $os_list ; do
  for arch in $arch_list ; do
    name=$cli-$version-$os-$arch

    if [ "$os" == "windows" ]; then
      name=$name.exe
    fi

    echo "$name"
    GOOS=$os GOARCH=$arch go build \
      -ldflags "
        -s -w
        -X main.appSemver=$version
        -X main.appCommit=$commit
        -X main.appBuilt=$built
        -X main.appOrigin=github.com/dpb587/gget
      " \
      -o tmp/build/$name \
      .
    
    # TODO 
    # if which upx > /dev/null ; then
    #   upx --ultra-brute tmp/build/$name
    # fi
  done
done
