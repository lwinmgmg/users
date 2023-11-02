FROM golang:1.21-bookworm AS compiler

ARG OUTPUT_DIR=bin
ARG OUTPUT_BINARY_NAME=user-go
ARG WORKDIR_NAME=/build

WORKDIR ${WORKDIR_NAME}

COPY . ${WORKDIR_NAME}/.

RUN pwd

RUN go get && go mod vendor

RUN go build -o "${OUTPUT_DIR}/${OUTPUT_BINARY_NAME}"

FROM debian:bookworm

ARG OUTPUT_DIR=bin
ARG OUTPUT_BINARY_NAME=user-go
ARG WORKDIR_NAME=/build

ENV USER_GOLANG_VERSION=1.21
ENV USER_USER=user
ENV USER_USER_HOME_DIR="/home/${USER_USER}"
ENV USER_USER_UID=2023
ENV USER_WORKDIR="/etc/user/config"
ENV USER_BINARY_DIR="/usr/local/user/bin"
ENV PATH=${PATH}:${USER_BINARY_DIR}
ENV SETTING_PATH="${USER_WORKDIR}/env.yaml"
ENV GIN_MODE=release

RUN groupadd --gid ${USER_USER_UID} ${USER_USER} \
    && useradd --uid ${USER_USER_UID} --gid ${USER_USER_UID} \
    --home-dir ${USER_USER_HOME_DIR} --shell /bin/bash ${USER_USER}

RUN apt update
RUN apt install -y ca-certificates

WORKDIR ${USER_WORKDIR}

# Copy compiled binary
COPY --from=compiler --chown=${USER_USER}:${USER_USER} "${WORKDIR_NAME}/${OUTPUT_DIR}/${OUTPUT_BINARY_NAME}" "${USER_BINARY_DIR}/user-go"

# To save vendor package if they go unsupported
COPY --from=compiler --chown=${USER_USER}:${USER_USER} "${WORKDIR_NAME}/vendor" "${USER_USER_HOME_DIR}/vendor"

EXPOSE 8069
EXPOSE 8079

CMD [ "user-go" ]
