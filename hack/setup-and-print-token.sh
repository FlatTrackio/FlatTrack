#!/bin/bash

SCRIPT_PATH=$(dirname $(realpath $0))
export LANGUAGE="${LANGUAGE:-English}"
export TIMEZONE="${TIMEZONE:-Pacific/Auckland}"
export FLATNAME="${FLATNAME:-Flat}"
export USER_NAMES="${USER_NAMES:-Flatter}"
export USER_EMAIL="${USER_EMAIL:-flattrack@example.com}"
export USER_PASSWORD="${USER_PASSWORD:-P@ssw0rd123!}"

token=$(curl \
    -H "Accept: application/json" \
    -X POST \
    http://localhost:8080/api/admin/register \
    -d "$(envsubst < $SCRIPT_PATH/defaultSetup.yaml | yq e -j -)" 2> /dev/null | jq -r .data 2>&1)

if [ ! $? -eq 0 ]; then
    echo "failed to register" > /dev/stderr
fi

set -a
. "$(git rev-parse --show-toplevel)/.env"
echo "go to: $APP_URL/login?authToken=$token"
echo
echo "token: $token"
