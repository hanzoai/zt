#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
# set -o xtrace

# this script is executed during the docker build, after the build context has been copied to /docker.build.context

ZITI_BIN_DIR="${1}"

if [ -d /docker.build.context/zt-bin ]; then
  mv /docker.build.context/zt-bin/ "${ZITI_BIN_DIR}"
else
  source /docker.build.context/zt-cli-functions.sh
  getZiti
fi
