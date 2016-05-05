#!/usr/bin/env bash
set -ETeu -o pipefail

DEFAULT_XENIA_COMMAND_PATH=/go/src/github.com/coralproject/xenia/cmd/xenia

# Create the database collections if environment variable XENIA_CREATE_DATABASE is present.
# XENIA_DATABASE_JSON - the location of the scrdb json.
if [ -z "${XENIA_CREATE_DATABASE:-}" ]; then
    /go/bin/xenia db create -f ${XENIA_DATABASE_JSON:-${DEFAULT_XENIA_COMMAND_PATH}/scrdb/database.json}
fi  

# Update queries if the environment variable XENIA_UPDATE_QUERY is present.
# XENIA_QUERY_JSON - the location of the scrquery directory.
if [ -z "${XENIA_UPDATE_QUERY:-}"]; then
    /go/bin/xenia query upsert -p ${XENIA_QUERY_JSON:-${DEFAULT_XENIA_COMMAND_PATH}/scrquery/}

# Run the Xenia daemon last and always to maintain backwards compatibility with previous Dockerfile version.
/go/bin/xeniad "$@"