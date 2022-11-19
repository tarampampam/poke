# syntax=docker/dockerfile:1.2

# Image page: <https://hub.docker.com/_/golang>
FROM golang:1.19-alpine as builder

# can be passed with any prefix (like `v1.2.3@GITHASH`)
# e.g.: `docker build --build-arg "APP_VERSION=v1.2.3@GITHASH" .`
ARG APP_VERSION="undefined@docker"

# This argument allows to install additional software for local development using docker and avoid it \
# in the production build
ARG DEV_MODE="false"

RUN set -x \
    && if [ "${DEV_MODE}" = "true" ]; then \
      # The following dependencies are needed for `go test` to work
      apk add --no-cache gcc musl-dev \
      # The following tool is used to format the imports in the source code
      && GOBIN=/bin go install golang.org/x/tools/cmd/goimports@latest \
    ;fi

COPY . /src

WORKDIR /src

# arguments to pass on each go tool link invocation
ENV LDFLAGS="-s -w -X github.com/tarampampam/poke/internal/version.version=$APP_VERSION"

RUN set -x \
    && go generate ./... \
    && CGO_ENABLED=0 go build -trimpath -ldflags "$LDFLAGS" -o /tmp/poke ./cmd/poke/ \
    && /tmp/poke --version \
    && /tmp/poke -h

# prepare rootfs for runtime
RUN mkdir -p /tmp/rootfs

WORKDIR /tmp/rootfs

RUN set -x \
    && mkdir -p \
        ./etc \
        ./bin \
    && echo 'appuser:x:10001:10001::/nonexistent:/sbin/nologin' > ./etc/passwd \
    && echo 'appuser:x:10001:' > ./etc/group \
    && mv /tmp/poke ./bin/poke

# use empty filesystem
FROM scratch as runtime

ARG APP_VERSION="undefined@docker"

LABEL \
    # Docs: <https://github.com/opencontainers/image-spec/blob/master/annotations.md>
    org.opencontainers.image.title="poke" \
    org.opencontainers.image.description="" \
    org.opencontainers.image.url="https://github.com/tarampampam/poke" \
    org.opencontainers.image.source="https://github.com/tarampampam/poke" \
    org.opencontainers.image.vendor="tarampampam" \
    org.opencontainers.version="$APP_VERSION" \
    org.opencontainers.image.licenses="MIT"

# Import from builder
COPY --from=builder /tmp/rootfs /

# Use an unprivileged user
USER 10001:10001

ENTRYPOINT ["/bin/poke"]
