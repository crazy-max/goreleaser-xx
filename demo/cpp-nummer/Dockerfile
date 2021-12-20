# syntax=docker/dockerfile:1.3-labs

ARG GO_VERSION=1.17
ARG GORELEASER_XX_IMAGE=crazymax/goreleaser-xx:edge
ARG XX_VERSION=1.1.0

FROM --platform=$BUILDPLATFORM ${GORELEASER_XX_IMAGE} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
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
    --name="xx-cpp-nummer" \
    --dist="/out" \
    --main="." \
    --ldflags="-s -w -linkmode external -extldflags -static" \
    --envs="GO111MODULE=auto" \
    --envs="CGO_ENABLED=1" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /

FROM scratch
COPY --from=build /usr/local/bin/xx-cpp-nummer /xx-cpp-nummer
ENTRYPOINT [ "/xx-cpp-nummer" ]