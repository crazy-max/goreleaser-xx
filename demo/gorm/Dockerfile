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
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,target=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored AS build
ARG TARGETPLATFORM
RUN xx-apk add --no-cache \
    gcc \
    musl-dev
# XX_CC_PREFER_STATIC_LINKER prefers ld to lld in ppc64le and 386.
ENV XX_CC_PREFER_STATIC_LINKER=1
RUN --mount=type=bind,source=.,target=/src,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --go-binary="xx-go" \
    --name="xx-gorm" \
    --dist="/out" \
    --main="." \
    --ldflags="-s -w -linkmode external -extldflags -static" \
    --envs="CGO_ENABLED=1" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /

FROM scratch
COPY --from=build /usr/local/bin/xx-gorm /xx-gorm
ENTRYPOINT [ "/xx-gorm" ]
