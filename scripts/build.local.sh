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

built=$( date -u +%Y-%m-%dT%H:%M:%S+00:00 )
commit=$( git rev-parse HEAD | cut -c-10 )

if [[ $( git clean -dnx | wc -l ) -gt 0 ]] ; then
  commit="${commit}+dirty"

  if [[ "${version}" != "0.0.0" ]]; then
    echo "ERROR: building an official version requires a clean repository"
    echo "WARN: refusing to clean repository"
    git clean -dnx

    exit 1
  fi
fi

mkdir -p tmp/build

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
      " \
      -o tmp/build/$name \
      .
    
    # TODO 
    # if which upx > /dev/null ; then
    #   upx --ultra-brute tmp/build/$name
    # fi
  done
done
