#!/usr/bin/env bash
set -euo pipefail

NEW_INSTALL=false
cd $(realpath $(dirname $0))
FOLDER=$(git rev-parse --show-toplevel 2> /dev/null)
IS_IN_GIT_REPO=$?
FOLDER_BIN=$FOLDER/bin
FOLDER_DIST=$FOLDER/web/dist
if [ ! $IS_IN_GIT_REPO -eq 0 ]; then
    FOLDER=$PWD
    FOLDER_BIN=$FOLDER
    FOLDER_DIST=$FOLDER/dist
fi
POSTGRES_FOLDER=$PREFIX/var/lib/postgresql
cd $FOLDER

if [ ! -d $POSTGRES_FOLDER ]; then
    NEW_INSTALL=true
    mkdir -p $POSTGRES_FOLDER
    initdb $POSTGRES_FOLDER
fi

if ! ps aux | grep -v "grep" | grep -q "postgres: checkpointer"; then
    pg_ctl -D $POSTGRES_FOLDER start
fi

if [ "$NEW_INSTALL" = "true" ]; then
    createuser --superuser $(whoami)
    psql postgres -c "ALTER USER $(whoami) WITH PASSWORD 'postgres';"
    createdb postgres
fi
psql postgres -c 'select 0;'

if [ ! -f "$FOLDER/.env" ]; then
    cat <<EOF > "$FOLDER/.env"
APP_PORT=127.0.0.1:8080
APP_DB_USERNAME=$(whoami)
APP_DB_PASSWORD=postgres
APP_DB_DATABASE=postgres
APP_DB_HOST=localhost
APP_METRICS_ENABLED=false
APP_HEALTH_ENABLED=false
APP_DB_MIGRATIONS_PATH=${FOLDER}/migrations
APP_DIST_FOLDER=${FOLDER_DIST}
EOF
fi

export APP_ENV_FILE=$FOLDER/.env
$FOLDER_BIN/flattrack
