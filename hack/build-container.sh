#!/bin/bash

function requiredEnv {
  name="${1}"
  value="${!1}"
  if [ -z "${value}" ]; then
    echo "error: '${name}' is required to be set" > /dev/stderr
    exit 1
  fi
}

function ko-publish-development {
  ko publish \
    --jobs 100 \
    --bare \
    --platform all \
    --local \
    .
}

function ko-publish-production {
  ko publish \
    --push \
    --jobs 100 \
    --bare \
    --platform all \
    --tags "${IMAGE_DESTINATIONS}" \
    .
}

cd $(dirname $0)
cd $(git rev-parse --show-toplevel)

rm -r ./kodata
mkdir -p ./kodata/{web/dist,}
cp -r ./migrations ./kodata
cp -r ./web/dist ./kodata/web

requiredEnvs=(APP_BUILD_HASH APP_BUILD_DATE APP_BUILD_VERSION APP_BUILD_MODE)
for env in ${requiredEnvs[*]}; do
  requiredEnv "${env}"
done
if [ ! "${1}" = "production" ]; then
  ko-publish-development
  exit $?
fi

requiredEnvs=(KO_DOCKER_REPO IMAGE_DESTINATIONS)
for env in ${requiredEnvs[*]}; do
  requiredEnv "${env}"
done
ko-publish-production
