#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# cd to project root
DIRNAME=$(dirname "$0")
cd "$DIRNAME/../.."  # the images can be built from anywhere, but the Dockerfile defaults assume the project root is the build context root

# TODO:
# - add an --arm64 flag to build the zt binary and container images for arm64
# - add a --namespace option to specify the registry hostname and namespace/org to tag each image, e.g., --namespace 127.0.0.1:5000/localtest
# - add a --push flag to push the images to the registry instead of loading them into the build context
# - add a --no-console flag to omit the console from the controller image
# - add a --console-version option to override the Dependabot-managed version of the console in the controller image
# - add a standalone console image w/ Angular or Node server?

# define a version based on the most recent tag
: "${ZITI_VERSION:=$(git describe --tags --always)}"

: build the go build env
docker buildx build \
    --tag=zt-go-builder \
    --build-arg uid=$UID \
    --load \
    ./dist/docker-images/cross-build/

: build the zt binary for amd64 in ARTIFACTS_DIR/TARGETARCH/TARGETOS/zt
docker run \
    --rm \
    --user "$UID" \
    --name=zt-go-builder \
    --volume=$PWD:/mnt/zt \
    --volume=${GOCACHE:-${HOME}/.cache/go-build}:/.cache/go-build \
    --env=GOCACHE=/.cache/go-build \
    zt-go-builder amd64

: build the cli image with binary from ARTIFACTS_DIR/TARGETARCH/TARGETOS/zt
docker buildx build \
  --platform=linux/amd64 \
  --tag "zt-cli:${ZITI_VERSION}" \
  --file ./dist/docker-images/zt-cli/Dockerfile \
  --load \
  $PWD 

docker build \
  --build-arg ZITI_CLI_IMAGE="zt-cli" \
  --build-arg ZITI_CLI_TAG="${ZITI_VERSION}" \
  --platform=linux/amd64 \
  --tag "zt-controller:${ZITI_VERSION}" \
  --file ./dist/docker-images/zt-controller/Dockerfile \
  --load \
  $PWD

docker build \
  --build-arg ZITI_CLI_IMAGE="zt-cli" \
  --build-arg ZITI_CLI_TAG="${ZITI_VERSION}" \
  --platform=linux/amd64 \
  --tag "zt-router:${ZITI_VERSION}" \
  --file ./dist/docker-images/zt-router/Dockerfile \
  --load \
  $PWD

docker build \
  --build-arg ZITI_CLI_IMAGE="zt-cli" \
  --build-arg ZITI_CLI_TAG="${ZITI_VERSION}" \
  --platform=linux/amd64 \
  --tag "zt-tunnel:${ZITI_VERSION}" \
  --file ./dist/docker-images/zt-tunnel/Dockerfile \
  --load \
  $PWD
