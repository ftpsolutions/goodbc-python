#!/usr/bin/env bash

set -e

if [[ -z "${MSSQL_PASSWORD}" ]]; then
  echo "error: need to have MSSQL_PASSWORD env var set"

  exit 1
fi

./test.sh ./_manual_test.sh
