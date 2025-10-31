#!/bin/bash

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
    KO_FLAGS="--insecure-registry --platform=linux/amd64 $KO_FLAGS"
    crane registry serve --address :5001 &
    REGPID=$!
    until curl -qsSL localhost:5001/v2; do
        sleep 2s
    done
    export KO_DOCKER_REPO=localhost:5001/ft
fi
if echo "${@:-}" | grep -q '\-\-insecure'; then
    KO_FLAGS="--insecure-registry $KO_FLAGS"
fi

cd "$(git rev-parse --show-toplevel)" || exit 1

LOCAL_SBOM_PATH="./tmp/sbom"
KO_DOCKER_REPO="${KO_DOCKER_REPO:-registry.gitlab.com/flattrack/flattrack}"
APP_BUILD_HASH="${APP_BUILD_HASH:-$(git rev-parse HEAD | cut -c -8)}"
APP_BUILD_DATE="$(git show -s --format=%cd --date=format:'%Y.%m.%d.%H%M')"
APP_BUILD_VERSION="${APP_BUILD_VERSION:-0.0.0}"
APP_BUILD_MODE="${APP_BUILD_MODE:-staging}"
IMAGE_DESTINATIONS="latest,${APP_BUILD_HASH}"
IMAGE_VERSION=latest
if [[ -n "${CI_COMMIT_TAG:-}" ]]; then
    APP_BUILD_VERSION="${CI_COMMIT_TAG:-}"
    APP_BUILD_MODE=production
    IMAGE_DESTINATIONS="$APP_BUILD_VERSION"
    IMAGE_VERSION="$APP_BUILD_VERSION"
fi
_IMAGE_LABELS="org.opencontainers.image.authors='FlatTrack https://flattrack.io'
org.opencontainers.image.created=$(git show -s --format=%cd --date=format:'%Y-%m-%dT%H:%M:%SZ')
org.opencontainers.image.version=$IMAGE_VERSION
org.opencontainers.image.source=https://gitlab.com/flattrack/flattrack
org.opencontainers.image.title=FlatTrack
org.opencontainers.image.url=https://flattrack.io
org.opencontainers.image.vendor=FlatTrack
org.opencontainers.image.description=Collaboration software for flats and living spaces
io.artifacthub.package.readme-url=https://gitlab.com/flattrack/flattrack/-/raw/main/README.md?ref_type=heads
io.artifacthub.package.license=AGPL-3.0
io.artifacthub.package.alternative-locations=index.docker.io/flattrack/flattrack"

IMAGE_LABELS="$(printf "$_IMAGE_LABELS" | tr '\n' ',')"

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
    --sbom-dir "$LOCAL_SBOM_PATH" \
    --image-label "$IMAGE_LABELS" \
    --image-annotation "$IMAGE_LABELS" \
    $KO_FLAGS \
    .)"

if [ "${SIGN:-}" = true ]; then
    syft . --output spdx-json=tmp/sbom/flattrack.spdx.json

    cosign sign --recursive -y "$IMAGE"
    # ko generates an SBOM for the container manifests
    cosign attest -y --recursive --predicate "$LOCAL_SBOM_PATH/flattrack-index.spdx.json" "$IMAGE"

    DIGESTS="$(crane manifest "$IMAGE" | jq -r '.manifests[].digest')"
    for DIGEST in $DIGESTS; do
        cosign attest -y --recursive --predicate "$LOCAL_SBOM_PATH/flattrack.spdx.json" "$KO_DOCKER_REPO@$DIGEST"
    done
fi

if [ "${TEST_TARBALL:-}" = true ]; then
    CORRECT=true
    crane config "$IMAGE" | jq
    crane manifest "$IMAGE" | jq
    RESULT="$(crane export "$IMAGE" - | tar -tvf - |
        grep -E '(etc/passwd|usr/share/zoneinfo|etc/ssl/certs|ko-app/flattrack|ko/web/assets/.*\.js|ko/migrations/.*\.sql|ko/sbom.spdx.json)')"
    CODE=$?
    if [ ! $CODE -eq 0 ]; then
        CORRECT=false
    fi
    TEMPDIR="$(mktemp -d)"
    crane export "$IMAGE" - | tar -xf - --exclude 'dev/*' -C "$TEMPDIR"
    RESULT="$(file "$TEMPDIR/ko-app/flattrack" | grep -q -E 'ELF.*statically linked')"
    CODE=$?
    if [ ! $CODE -eq 0 ]; then
        CORRECT=false
    fi
    echo "$RESULT"
    if [ "$CORRECT" = false ]; then
        echo "error: failed to build image correctly" >/dev/stderr
        echo "$RESULT" >/dev/stderr
        kill "$REGPID"
        exit 1
    fi
    echo "success: image is built correctly"
    kill "$REGPID"
fi

IMAGE_DEST_REFS=""
for DEST in $(echo "$IMAGE_DESTINATIONS" | tr ',' ' '); do
    IMAGE_DEST_REFS="\n- $KO_DOCKER_REPO:$DEST$IMAGE_DEST_REFS"
done

echo -e "Published image to: \n- $IMAGE $IMAGE_DEST_REFS"
if [ -n "${IMAGE_RESULT_FILE:-}" ]; then
    printf "%s" "$IMAGE" >"$IMAGE_RESULT_FILE"
fi
