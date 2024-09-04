#!/usr/bin/env bash

set -e
set -u
set -o pipefail

#if [ -n "${PARAMETER_STORE:-}" ]; then
#  export GESTION_DEPENDENCIAS_MID_PGUSER="$(aws ssm get-parameter --name /${PARAMETER_STORE}gestion_dependencias_mid/db/username --output text --query Parameter.Value)"
#  export GESTION_DEPENDENCIAS_MID_PGPASS="$(aws ssm get-parameter --with-decryption --name /${PARAMETER_STORE}/gestion_dependencias_mid/db/password --output text --query Parameter.Value)"

exec ./main "$@"
