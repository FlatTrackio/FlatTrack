#!/bin/bash

function usage() {
    echo "Command help: $(basename $0) 1.0.0 1.1.0"
}

function prepare-command() {
    echo sed -i "s/${CURRENT_VERSION}/${NEW_VERSION}/g"
}

SEARCH_IN=(
    docs/
    pkg/
    build/
    cmd/
    deployments/
    docs/
    README.org
    LICENSE
    main.go
    migrations/
    templates/
    test
    web/package.json
    fly.toml
)

cd "$(git rev-parse --show-toplevel)"

CURRENT_VERSION="${1}"
NEW_VERSION="${2}"

if [ -z "${CURRENT_VERSION}" ] || [ -z "${NEW_VERSION}" ]; then
    usage
    exit 0
fi

find ${SEARCH_IN[*]} -exec $(prepare-command) {} \;
