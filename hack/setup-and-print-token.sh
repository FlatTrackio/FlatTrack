#!/bin/bash

LANGUAGE=${LANGUAGE:-English}
TIMEZONE=${TIMEZONE:-Pacific/Auckland}
FLATNAME=${FLATNAME:-Flat}
USER_NAMES=${USER_NAMES:-Flatter}
USER_EMAIL=${USER_EMAIL:-flattrack@example.com}
USER_PASSWORD=${USER_PASSWORD:-P@ssw0rd123!}

token=$(curl \
    -H "Accept: application/json" \
    -X POST \
    http://localhost:8080/api/admin/register \
    -d "{
\"language\": \"$LANGUAGE\",
\"timezone\": \"$TIMEZONE\",
\"flatName\": \"$FLATNAME\",
\"user\": {
    \"names\": \"$USER_NAMES\",
    \"email\": \"$USER_EMAIL\",
    \"password\": \"$USER_PASSWORD\"
  }
}" 2> /dev/null | jq -r .data 2>&1)

if [ ! $? -eq 0 ]; then
    echo "failed to register" > /dev/stderr
fi

set -a
. "$(git rev-parse --show-toplevel)/.env"
echo "go to: $APP_URL/login?authToken=$token"
echo
echo "token: $token"
