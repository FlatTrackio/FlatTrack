#!/bin/bash

cd $(dirname $0)
cd $(git rev-parse --show-toplevel)

rm -r ./kodata
mkdir -p ./kodata/{web/dist,}
cp -r ./migrations ./kodata
cp -r ./web/dist ./kodata/web

ko publish \
  --jobs 100 \
  --bare \
  --platform all \
  --local \
  .
