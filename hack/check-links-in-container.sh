#!/bin/sh -x

set -o errexit
set -o nounset
set -o pipefail

cd "$(git rev-parse --show-toplevel)" || exit 0

C_DIR="/builds/$(basename $PWD)"
podman run -it --rm \
  -v "$PWD:$C_DIR:ro" --workdir "$C_DIR" \
  --entrypoint /bin/sh \
  --platform linux/amd64 \
  docker.io/lycheeverse/lychee:sha-a7c11c9-alpine \
  ./hack/check-links.sh
