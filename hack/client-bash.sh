#!/bin/bash -x

if [ -z "$FT_TOKEN" ]; then
    echo "error: \$FT_TOKEN must be set to authenticate" > /dev/stderr
fi

VERB="$1"
PATH="$2"
DATA="$3"

/usr/bin/curl \
    -X "$VERB" \
    -H "Accept: application/json" \
    -H "Authorization: bearer $FT_TOKEN" \
    "http://localhost:8080/$PATH" \
    -d "$DATA" \
    | /usr/bin/jq .
