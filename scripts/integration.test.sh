#!/bin/bash

set -eu

mkdir -p tmp/integrationtest

go build -o tmp/integrationtest/gget .

rm -fr tmp/integrationtest/workdir
mkdir tmp/integrationtest/workdir
cd tmp/integrationtest/workdir

../gget github.com/gohugoio/hugo@v0.63.1 --exclude='*extended*' 'hugo_*_Linux-ARM.deb'

diff <( shasum * ) - <<EOF
34ee738fc56c3eb479a8ff71ef88e6bfd6ad0bde  hugo_0.63.1_Linux-ARM.deb
EOF

rm *

../gget --exec github.com/stedolan/jq@jq-1.6 my-custom-name=jq-osx-amd64

diff <( shasum * ) - <<EOF
8673400d1886ed051b40fe8dee09d89237936502  my-custom-name
EOF

rm *

../gget --type=blob github.com/stedolan/jq@jq-1.5-branch README.md

diff <( shasum * ) - <<EOF
cded31e0fbf9b7dbf9e6ffa9132201ce1d0b0f2d  README.md
EOF

rm *

../gget --stdout github.com/buildpacks/pack@v0.8.1 '*macos*' | tar -xzf-

diff <( shasum * ) - <<EOF
1fe75bead2f16823f0bdb182f666afc2176cb6a5  pack
EOF

rm *

../gget --help > /dev/null

cd ../../

rm -fr tmp/integrationtest
