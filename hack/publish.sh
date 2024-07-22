#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

KO_FLAGS="${KO_FLAGS:-}"
# NOTE budget /bin/sh way
if echo "${@:-}" | grep -q '\-\-debug'; then
    set -x
    KO_FLAGS="--verbose $KO_FLAGS"
fi
if echo "${@:-}" | grep -q '\-\-sign'; then
    SIGN=true
fi
if echo "${@:-}" | grep -q '\-\-local'; then
    KO_FLAGS="--local $KO_FLAGS"
fi
if echo "${@:-}" | grep -q '\-\-tarball-test-only'; then
    TEST_TARBALL=true
    KO_FLAGS="--tarball /tmp/flattrack.tar --platform=linux/amd64 --push=false $KO_FLAGS"
fi
if echo "${@:-}" | grep -q '\-\-insecure'; then
    KO_FLAGS="--insecure-registry $KO_FLAGS"
fi

cd "$(git rev-parse --show-toplevel)" || exit 1

KO_DOCKER_REPO="${KO_DOCKER_REPO:-registry.gitlab.com/flattrack/flattrack}"
APP_BUILD_HASH="${APP_BUILD_HASH:-$(git rev-parse HEAD | cut -c -8)}"
APP_BUILD_DATE="$(git show -s --format=%cd --date=format:'%Y.%m.%d.%H%M')"
APP_BUILD_VERSION="${APP_BUILD_VERSION:-0.0.0}"
APP_BUILD_MODE="${APP_BUILD_MODE:-staging}"
IMAGE_DESTINATIONS="latest,${APP_BUILD_HASH}"
if [[ -n "${CI_COMMIT_TAG:-}" ]]; then
  APP_BUILD_VERSION="${CI_COMMIT_TAG:-}"
  APP_BUILD_MODE=production
  IMAGE_DESTINATIONS="$APP_BUILD_VERSION"
fi
echo "Commit made on '${APP_BUILD_DATE:-}'"

export KO_DOCKER_REPO \
    APP_BUILD_HASH \
    APP_BUILD_DATE \
    APP_BUILD_MODE \
    APP_BUILD_VERSION \
    IMAGE_DESTINATIONS

IMAGE="$(ko publish \
    --bare \
    --tags "${IMAGE_DESTINATIONS}" \
    $KO_FLAGS \
    .)"

if [ "${SIGN:-}" = true ]; then
    cosign sign --recursive -y "$IMAGE"
    cosign download sbom "$IMAGE" > /tmp/sbom-spdx.json
    cosign attest -y --recursive --predicate /tmp/sbom-spdx.json "$IMAGE"
fi

if [ "${TEST_TARBALL:-}" = true ]; then
    crane registry serve --address :5001&
    REGPID=$!
    IMG=$(crane --insecure push /tmp/flattrack.tar localhost:5001/ft)

    if ! crane export "$IMG" - | tar -tvf - \
        | grep -q -E '(etc/passwd|usr/share/zoneinfo|etc/ssl/certs|ko-app/flattrack|ko/web/assets/.*\.js|ko/migrations/.*\.sql)'; then
        echo "error: failed to build image correctly" >/dev/stderr
        kill "$REGPID"
        exit 1
    fi
    echo "success: image is built correctly"
    kill "$REGPID"
fi
