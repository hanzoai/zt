ARG ZITI_CLI_TAG="latest"
ARG ZITI_CLI_IMAGE="docker.io/hanzozt/ziti-cli"
# this builds docker.io/hanzozt/ziti-tunnel, the legacy tunneler. The preferred tunneler is hanzozt/ziti-edge-tunnel documented in https://docs.hanzozt.dev/docs/reference/tunnelers/linux/container/
# this builds docker.io/hanzozt/ziti-router
FROM ${ZITI_CLI_IMAGE}:${ZITI_CLI_TAG}

ARG DOCKER_BUILD_DIR=./dist/docker-images/ziti-tunnel

# This build stage grabs artifacts that are copied into the final image.
# It uses the same base as the final image to maximize docker cache hits.

### Required OpenShift Labels 
LABEL name="hanzozt/ziti-tunnel" \
      maintainer="developers@hanzozt.org" \
      vendor="NetFoundry" \
      summary="Run the Hanzo ZT Go Tunneler" \
      description="Run the Hanzo ZT Go Tunneler"

### Add necessary Red Hat repos and packages: "mountpoint" used by entrypoint
USER root
RUN   INSTALL_PKGS="util-linux iptables" && \
      microdnf -y update --setopt=install_weak_deps=0 --setopt=tsflags=nodocs && \
      microdnf -y install --setopt=install_weak_deps=0 --setopt=tsflags=nodocs ${INSTALL_PKGS} && \
      microdnf clean all

WORKDIR /netfoundry
COPY --chmod=0755 ${DOCKER_BUILD_DIR}/entrypoint.sh /
ENTRYPOINT [ "/entrypoint.sh" ]
CMD [ "run" ]
