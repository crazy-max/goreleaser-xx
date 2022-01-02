# syntax=docker/dockerfile:1

ARG GO_VERSION=1.17
ARG GORELEASER_XX_BASE=crazymax/goreleaser-xx:edge
ARG XX_VERSION=1.1.0

FROM --platform=$BUILDPLATFORM ${GORELEASER_XX_BASE} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
ENV CGO_ENABLED=1
COPY --from=goreleaser-xx / /
COPY --from=xx / /
RUN apk add --no-cache \
    clang \
    git \
    file \
    lld \
    llvm \
    pkgconfig
WORKDIR /go/src/github.com/crazy-max/goreleaser-xx/demo/cpp

FROM base AS build
ARG TARGETPLATFORM
RUN xx-apk add --no-cache \
    gcc \
    g++ \
    linux-headers \
    musl-dev
# XX_CC_PREFER_STATIC_LINKER prefers ld to lld in ppc64le and 386.
ENV XX_CC_PREFER_STATIC_LINKER=1
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache \
  goreleaser-xx --debug \
    --go-binary="xx-go" \
    --name="cpp-xx-alpine" \
    --dist="/out" \
    --artifacts="bin" \
    --main="." \
    --ldflags="-s -w -extldflags -static" \
    --envs="GO111MODULE=auto" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out /

FROM scratch
COPY --from=build /usr/local/bin/cpp-xx-alpine /cpp-xx-alpine
ENTRYPOINT [ "/cpp-xx-alpine" ]
