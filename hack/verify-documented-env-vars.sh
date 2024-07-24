#!/bin/sh

set -o errexit
set -o nounset

FAIL=false
for VAR in $(grep 'GetEnvOrDefault("' internal/common/common.go | sed 's/.*("\(.*\)",.*/\1/g' | xargs) ; do
    if ! grep -q -E "\`$VAR\`" ./docs/configuration.md; then
        echo "$VAR not found"
        FAIL=true
    fi
done

if [ "$FAIL" = true ]; then
    exit 1
fi
