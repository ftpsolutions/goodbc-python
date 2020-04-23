#!/usr/bin/env bash

set -e

if [[ -z "${MSSQL_PASSWORD}" ]]; then
  echo "error: need to have MSSQL_PASSWORD env var set"

  exit 1
fi

cd /workspace/dist/go*-python*/

MSSQL_PASSWORD="${MSSQL_PASSWORD}" python manual_test.py
