#!/bin/bash

set -euo pipefail

export BASE_URL=https://gget.io/
export GA_TRACKING_ID=UA-37464314-3

suffix="$( go run . --version -vv | grep ^runtime | tr -d ';' | awk '{ print $7 "-" $5 }' )"
go run . github.com/dpb587/gget --no-download --export=json > website/static/latest.json
go run . github.com/dpb587/gget --executable website/static/gget="*-${suffix}"
./website/static/gget --help > website/static/latest-help.txt
rm website/static/gget

cd website/

yarn build
yarn start &
nuxt=$!

sleep 5
rm -fr dist
cp -rp static dist
cp -rp .nuxt/dist/client dist/_nuxt

curl http://localhost:3000/ > dist/index.html

kill "${nuxt}"

cd dist

touch .nojekyll
echo -n gget.io > CNAME

git init .
git add .
git commit -m 'regenerate'
git branch -m gh-pages
git remote add origin git@github.com:dpb587/gget.git
git config branch.gh-pages.remote origin
git config branch.gh-pages.merge refs/heads/gh-pages
