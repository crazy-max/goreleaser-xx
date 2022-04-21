# syntax=docker/dockerfile:1

ARG GORELEASER_VERSION="1.8.3"
ARG XX_VERSION="1.1.0"
ARG GO_VERSION="1.18"

FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
ENV CGO_ENABLED=0
COPY --from=xx / /
RUN apk --update --no-cache add file git
WORKDIR /src

FROM base AS goreleaser
ARG GORELEASER_VERSION
WORKDIR /goreleaser
RUN git clone --branch v${GORELEASER_VERSION} https://github.com/goreleaser/goreleaser .
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download
ARG TARGETPLATFORM
RUN --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  xx-go build -v -ldflags "-w -s -X 'main.version=${GORELEASER_VERSION}' -X main.builtBy=goreleaser-xx" \
  && xx-verify --static ./goreleaser

FROM base AS vendored
RUN --mount=type=bind,target=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download && \
  mkdir /out && cp go.mod go.sum /out

FROM scratch AS vendor-update
COPY --from=vendored /out /

FROM vendored AS build
ARG GIT_REF
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=target=/go/pkg/mod,type=cache \
  case "$GIT_REF" in \
    refs/tags/v*) gitTag="${GIT_REF#refs/tags/}" ;; \
    *) gitTag=$(git describe --match 'v[0-9]*' --dirty='.m' --always --tags) ;; \
  esac \
  && xx-go build -v -ldflags "-w -s -X 'main.version=${gitTag}'" -o /usr/local/bin/goreleaser-xx \
  && xx-verify --static /usr/local/bin/goreleaser-xx

FROM scratch AS release
COPY --from=goreleaser /goreleaser/goreleaser /usr/local/bin/goreleaser
COPY --from=build /usr/local/bin/goreleaser-xx /usr/local/bin/goreleaser-xx
