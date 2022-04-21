# syntax=docker/dockerfile:1

ARG GO_VERSION=1.18
ARG GORELEASER_XX_BASE=crazymax/goreleaser-xx:edge

FROM --platform=$BUILDPLATFORM ${GORELEASER_XX_BASE} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM crazymax/osxcross:11.3 AS osxcross
FROM --platform=$BUILDPLATFORM crazymax/goxx:${GO_VERSION} AS base
COPY --from=goreleaser-xx / /
ENV CGO_ENABLED=1
RUN goxx-apt-get install --no-install-recommends -y git
WORKDIR /go/src/github.com/crazy-max/goreleaser-xx/demo/cpp

FROM base AS build
ARG TARGETPLATFORM
RUN --mount=type=cache,sharing=private,target=/var/cache/apt \
  --mount=type=cache,sharing=private,target=/var/lib/apt/lists \
  goxx-apt-get install -y binutils gcc g++ pkg-config
RUN --mount=type=bind,source=.,rw \
  --mount=from=osxcross,target=/osxcross,src=/osxcross,rw \
  --mount=type=cache,target=/root/.cache \
  goreleaser-xx --debug \
    --name="cpp-nummer" \
    --dist="/out" \
    --main="." \
    --ldflags="-s -w" \
    --envs="GO111MODULE=auto" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out /

FROM scratch
COPY --from=build /usr/local/bin/cpp-nummer /cpp-nummer
ENTRYPOINT [ "/cpp-nummer" ]
