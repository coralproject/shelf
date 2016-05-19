#!/usr/bin/env bash
set -ETeu -o pipefail

DEFAULT_XENIA_COMMAND_PATH=/go/src/github.com/coralproject/xenia/cmd/xenia
PATH=/go/bin:${PATH}

# Create the database collections if environment variable XENIA_CREATE_DATABASE is present.
# XENIA_DATABASE_SCHEMA_PATH - the location of the scrdb json.
if [ ! -z ${XENIA_CREATE_DATABASE:-} ]; then
    xenia db create -f ${XENIA_DATABASE_SCHEMA_PATH:-${DEFAULT_XENIA_COMMAND_PATH}/scrdb/database.json}
fi  

# Update queries if the environment variable XENIA_UPDATE_QUERY is present.
# XENIA_QUERY_UPSERT_PATH - the location of the scrquery directory or individual query file.
if [ ! -z ${XENIA_UPDATE_QUERY:-} ]; then
    xenia query upsert -p ${XENIA_QUERY_UPSERT_PATH:-${DEFAULT_XENIA_COMMAND_PATH}/scrquery}
fi

# Run the Xenia daemon last and always to maintain backwards compatibility with previous Dockerfile version.
xeniad "$@"
