#!/bin/sh -x

set -o errexit
set -o nounset
set -o pipefail

cd "$(git rev-parse --show-toplevel)" || exit 0

C_DIR="/builds/$(basename $PWD)"
podman run --rm --network=host \
    -v "$PWD:$C_DIR" --workdir "$C_DIR" \
    docker.io/python:3.8-buster@sha256:04c3f641c2254c229fd2f704c5199ff4bea57d26c1c29008ae3a4afddde98709 \
      sh -c "
git config --global --add safe.directory $C_DIR ;
./hack/build-docs.sh ${*:-}
"
