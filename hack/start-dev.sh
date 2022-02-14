#!/bin/bash

cd $(dirname $0)
cd $(git rev-parse --show-toplevel)

CONTAINER_RUNTIME="${1:-podman}"
TMATE_SOCKET="/tmp/tmate-flattrack-${USER}"

(
    until tmate -S "${TMATE_SOCKET}" wait-for tmate-ready; do
        echo "Waiting"
        sleep 1
    done
    tmate -F -v -S "${TMATE_SOCKET}" new-window -d -c "$PWD" -n ft-pg bash
    tmate -S "${TMATE_SOCKET}" send-keys -t ft-pg "${CONTAINER_RUNTIME} run -it --rm --env-file .env --name flattrack-postgres -p 5432:5432 postgres:12.7-alpine" Enter

    tmate -F -v -S "${TMATE_SOCKET}" new-window -d -c "$PWD" -n ft-pg-psql bash
    tmate -S "${TMATE_SOCKET}" send-keys -t ft-pg-psql "until nc -zv localhost 5432; do sleep 1s; done" Enter
    tmate -S "${TMATE_SOCKET}" send-keys -t ft-pg-psql "${CONTAINER_RUNTIME} exec -it flattrack-postgres psql" Enter

    tmate -F -v -S "${TMATE_SOCKET}" new-window -d -c "$PWD" -n ft-minio bash
    tmate -S "${TMATE_SOCKET}" send-keys -t ft-minio "${CONTAINER_RUNTIME} run --env-file .env -p 9000:9000 -p 9001:9001 minio/minio:RELEASE.2022-01-07T01-53-23Z server /data --console-address \":9001\"" Enter

    sleep 2s
    tmate -F -v -S "${TMATE_SOCKET}" new-window -d -c "$PWD" -n ft-backend bash
    tmate -S "${TMATE_SOCKET}" send-keys -t ft-backend "go run ." Enter

    tmate -F -v -S "${TMATE_SOCKET}" new-window -d -c "$PWD/web" -n ft-frontend bash
    tmate -S "${TMATE_SOCKET}" send-keys -t ft-frontend "npm run build" Enter
)&

tmate -S "${TMATE_SOCKET}"