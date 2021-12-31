# syntax=docker/dockerfile:1.3-labs

ARG GO_VERSION=1.17
ARG GORELEASER_XX_BASE=crazymax/goreleaser-xx:edge

FROM --platform=$BUILDPLATFORM ${GORELEASER_XX_BASE} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM crazymax/goxx:${GO_VERSION} AS base
ENV CGO_ENABLED=1
COPY --from=goreleaser-xx / /
WORKDIR /go/src/github.com/crazy-max/goreleaser-xx/demo/cpp

FROM base AS build
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache \
  goreleaser-xx --debug \
    --name="cpp-goxx" \
    --dist="/out" \
    --artifacts="bin" \
    --main="." \
    --ldflags="-s -w" \
    --envs="GO111MODULE=auto" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out /

FROM scratch
COPY --from=build /usr/local/bin/cpp-goxx /cpp-goxx
ENTRYPOINT [ "/cpp-goxx" ]
