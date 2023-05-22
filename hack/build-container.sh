#!/bin/bash

function requiredEnv {
  name="${1}"
  value="${!1}"
  if [ -z "${value}" ]; then
    echo "error: '${name}' is required to be set" >/dev/stderr
    exit 1
  fi
}
function ko-publish {
  ko publish \
    --jobs 100 \
    --bare \
    .
}

cd $(dirname $0)
cd $(git rev-parse --show-toplevel)

mkdir -p ./kodata/web
(
  cd web
  npm run build
)
cp -r ./web/dist ./kodata/web

export \
  APP_BUILD_HASH="$(git rev-parse HEAD)-dirty" \
  APP_BUILD_DATE="$(date +%Y.%m.%d.%H%M)" \
  APP_BUILD_VERSION="0.0.0" \
  APP_BUILD_MODE="development" \
  KO_DOCKER_REPO="${KO_DOCKER_REPO:-registry.gitlab.com/flattrack/flattrack}"
ko-publish
