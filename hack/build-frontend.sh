#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

cd "$(git rev-parse --show-toplevel)/web" || exit 1
rm -rf ../kodata/web

npm run build
