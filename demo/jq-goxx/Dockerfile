# syntax=docker/dockerfile:1

ARG GO_VERSION=1.17
ARG GORELEASER_XX_BASE=crazymax/goreleaser-xx:edge

FROM --platform=$BUILDPLATFORM ${GORELEASER_XX_BASE} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM crazymax/osxcross:11.3 AS osxcross
FROM --platform=$BUILDPLATFORM crazymax/goxx:${GO_VERSION} AS base
COPY --from=osxcross /osxcross /osxcross
COPY --from=goreleaser-xx / /
ENV GO111MODULE=auto
ENV CGO_ENABLED=1
ENV OSXCROSS_MP_INC=1
RUN goxx-apt-get install --no-install-recommends -y git
WORKDIR /src

FROM base AS build
ARG TARGETPLATFORM
RUN --mount=type=cache,sharing=private,target=/var/cache/apt \
  --mount=type=cache,sharing=private,target=/var/lib/apt/lists \
  goxx-apt-get install -y binutils gcc pkg-config libjq-dev libonig-dev
RUN goxx-macports --static install jq
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache \
  goreleaser-xx --debug \
    --config=".goreleaser.yml" \
    --name="jq-goxx" \
    --dist="/out" \
    --main="." \
    --ldflags="-s -w" \
    --tags="netgo" \
    --files="README.md"

FROM scratch
COPY --from=build /out /
