FROM ubuntu:rolling AS fetch-zt-bins

# optional arg to specify which version to fetch when local "zt-bin" directory is not present
ARG ZITI_VERSION_OVERRIDE
ARG DEBIAN_FRONTEND=noninteractive
ARG GITHUB_REPO_OWNER=hanzozt
ARG GITHUB_REPO_NAME=zt
ENV GITHUB_REPO_OWNER=${GITHUB_REPO_OWNER}
ENV GITHUB_REPO_NAME=${GITHUB_REPO_NAME}

RUN apt-get update \
    && apt-get --yes install \
        jq \
        curl \
    && apt-get --yes autoremove \
    && apt-get clean autoclean \
    && rm -fr /var/lib/apt/lists/{apt,dpkg,cache,log} /tmp/* /var/tmp/*
    # to use a locally built `zt` cli \
    #   create a ./zt-bin directory adjacent to this Dockerfile \
    #   and include any desired binaries to include in the image
    # to use the latest released version of zt \
    #   just build the Dockerfile \
    # to use a specific version of zt, specify ZITI_VERSION_OVERRIDE then  \
    #   build the Dockerfile
COPY . /docker.build.context
RUN bash /docker.build.context/fetch-zt-bins.sh /zt-bin

FROM ubuntu:rolling
ARG RUNUSER=zt
ARG RUNGROUP=zt
ARG RUNUSERID=2171

RUN apt-get update \
    && apt-get --yes install \
        jq \
        curl \
        netcat-openbsd \
        vim \
        inetutils-ping \
        net-tools \
        lsof \
    && apt-get --yes autoremove \
    && apt-get clean autoclean \
    && rm -fr /var/lib/apt/lists/{apt,dpkg,cache,log} /tmp/* /var/tmp/*

RUN useradd -u $RUNUSERID -ms /bin/bash $RUNUSER && \
    mkdir /persistent && chown $RUNUSER:$RUNGROUP /persistent && \
    mkdir /persistent/scripts && chown $RUNUSER:$RUNGROUP /persistent/scripts

USER $RUNUSER

WORKDIR /persistent

ENV ZITI_HOME=/persistent
ENV ZITI_SHARED="${ZITI_HOME}/shared"
ENV ZITI_ENV_FILE="${ZITI_HOME}/zt.env"
ENV ZITI_NETWORK=zt
# Don't put these paths on shared volume in docker-compose deployment (https://github.com/hanzozt/zt/issues/912)
ENV ZITI_BIN_DIR=/var/hanzozt/zt-bin
ENV ZITI_BIN_ROOT="${ZITI_BIN_DIR}"
ENV ZITI_SCRIPTS=/var/hanzozt/scripts

# copy the zt binaries to a directory already on the path
COPY --chown=zt --from=fetch-zt-bins /zt-bin "${ZITI_BIN_DIR}"
COPY --chown=zt zt-cli-functions.sh "${ZITI_SCRIPTS}/"
COPY --chown=zt run-controller.sh "${ZITI_SCRIPTS}/"
COPY --chown=zt run-router.sh "${ZITI_SCRIPTS}/"
COPY --chown=zt run-with-zt-cli.sh "${ZITI_SCRIPTS}/"
COPY --chown=zt access-control.sh "${ZITI_SCRIPTS}/"
COPY --chown=zt run-router-external.sh "${ZITI_SCRIPTS}/"

RUN /bin/bash -c " \
    source ${ZITI_SCRIPTS}/zt-cli-functions.sh; \
    persistEnvironmentValues; \
    echo source ${ZITI_HOME}/zt.env \
        >> ~/.bashrc \
    "
