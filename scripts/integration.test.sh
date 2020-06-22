#!/bin/bash

set -eu

mkdir -p tmp/integrationtest

go build -o tmp/integrationtest/gget .

rm -fr tmp/integrationtest/workdir
mkdir tmp/integrationtest/workdir
cd tmp/integrationtest/workdir

../gget github.com/dpb587/gget --ref-version=0.2.x --dump-info=info.txt '*linux*'

diff <( shasum * ) - <<EOF
734d4ef1448dd9892852ae370933e7629fe528d5  gget-0.2.0-linux-amd64
92a43b6a3eb807c26f0c8a76ff0fec96621de742  info.txt
EOF

rm *

../gget github.com/dpb587/gget --no-download

[[ "$( ls -l . )" == "total 0" ]]

../gget github.com/dpb587/gget --ref-version=0.2.x --dump-info=info.txt '*linux*'

diff <( shasum * ) - <<EOF
734d4ef1448dd9892852ae370933e7629fe528d5  gget-0.2.0-linux-amd64
92a43b6a3eb807c26f0c8a76ff0fec96621de742  info.txt
EOF

rm *

../gget github.com/gohugoio/hugo@v0.63.1 --exclude='*extended*' 'hugo_*_Linux-ARM.deb'

diff <( shasum * ) - <<EOF
34ee738fc56c3eb479a8ff71ef88e6bfd6ad0bde  hugo_0.63.1_Linux-ARM.deb
EOF

rm *

../gget --executable github.com/stedolan/jq@jq-1.6 my-custom-name=jq-osx-amd64

diff <( shasum * ) - <<EOF
8673400d1886ed051b40fe8dee09d89237936502  my-custom-name
EOF
[ -e my-custom-name ]

rm *

../gget --type=blob github.com/stedolan/jq@jq-1.5-branch README.md

diff <( shasum * ) - <<EOF
cded31e0fbf9b7dbf9e6ffa9132201ce1d0b0f2d  README.md
EOF

rm *

../gget --type=blob github.com/stedolan/jq@a17dd3248a README.md

diff <( shasum * ) - <<EOF
1c336249ffa502059d99ac700579c90382b0462b  README.md
EOF

rm *

../gget --stdout github.com/buildpacks/pack@v0.8.1 '*macos*' | tar -xzf-

diff <( shasum * ) - <<EOF
1fe75bead2f16823f0bdb182f666afc2176cb6a5  pack
EOF

rm *

../gget gitlab.com/gitlab-org/gitlab@v12.10.0-ee 'gitlab-*-released'

grep -q 'GitLab 12.10 released with Requirements Management, Autoscaling CI on AWS Fargate' gitlab-*-released

rm *

../gget gitlab.com/gitlab-org/gitlab-runner --ref-version=11.x 'gitlab-runner_amd64.deb'

diff <( shasum * ) - <<EOF
8b5f4e982e692331571fc9cdf055f9f73a74b09d  gitlab-runner_amd64.deb
EOF

rm *

../gget --version > /dev/null

../gget --version=0.0.0 > /dev/null

../gget --help > /dev/null

cd ../../

rm -fr tmp/integrationtest

echo Success
