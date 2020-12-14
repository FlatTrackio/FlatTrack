#!/bin/bash

SCRIPT_PATH="$(dirname $(realpath $0))"

# wait until API is ready
echo "waiting for FlatTrack API to be ready"
until curl -H "Accept: application/json" http://localhost:8080/api 2>&1 > /dev/null; do
    echo "FlatTrack API not ready yet"
    sleep 1
done
echo "FlatTrack API is now ready"

export APP_DB_MIGRATIONS_PATH="$SCRIPT_PATH/../migrations"
go run "$SCRIPT_PATH"/cmd/insertcollaborators/insertcollaborators.go
