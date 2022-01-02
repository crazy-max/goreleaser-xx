# syntax=docker/dockerfile:1

ARG GO_VERSION=1.17
ARG GORELEASER_XX_BASE=crazymax/goreleaser-xx:edge
ARG XX_VERSION=master

FROM --platform=$BUILDPLATFORM ${GORELEASER_XX_BASE} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bullseye AS base
ENV CGO_ENABLED=1
COPY --from=goreleaser-xx / /
COPY --from=xx / /
RUN apt-get update \
  && apt-get install --no-install-recommends -y \
    clang \
    file \
    git \
    lld \
    llvm \
    pkg-config
WORKDIR /go/src/github.com/crazy-max/goreleaser-xx/demo/c

FROM base AS build
ARG TARGETPLATFORM
RUN xx-apt-get install -y \
    binutils \
    gcc \
    g++ \
    libc6-dev
# XX_CC_PREFER_STATIC_LINKER prefers ld to lld in ppc64le and 386.
ENV XX_CC_PREFER_STATIC_LINKER=1
RUN --mount=type=bind,source=.,rw \
  --mount=from=dockercore/golang-cross:xx-sdk-extras,target=/xx-sdk,src=/xx-sdk \
  --mount=type=cache,target=/root/.cache \
  goreleaser-xx --debug \
    --go-binary="xx-go" \
    --name="c-xx-debian" \
    --dist="/out" \
    --artifacts="archive" \
    --artifacts="bin" \
    --main="." \
    --ldflags="-s -w" \
    --envs="GO111MODULE=auto" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out /

FROM scratch
COPY --from=build /usr/local/bin/c-xx-debian /c-xx-debian
ENTRYPOINT [ "/c-xx-debian" ]
