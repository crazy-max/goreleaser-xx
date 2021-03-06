# syntax=docker/dockerfile:1.2
ARG GORELEASER_VERSION=0.159.0
ARG GO_VERSION=1.16

FROM --platform=${BUILDPLATFORM:-linux/amd64} tonistiigi/xx:golang AS xgo
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:${GO_VERSION}-alpine AS base
COPY --from=xgo / /
RUN apk --update --no-cache add \
    bash \
    build-base \
    git \
  && rm -rf /tmp/* /var/cache/apk/*
ENV CGO_ENABLED=0
WORKDIR /src

FROM base AS goreleaser
ARG GORELEASER_VERSION
RUN git clone --branch v${GORELEASER_VERSION} https://github.com/goreleaser/goreleaser .
ARG TARGETPLATFORM
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod tidy \
  && go build -v -ldflags "-w -s -X 'main.version=${GORELEASER_VERSION}' -X main.builtBy=goreleaser-xx" \
  && ./goreleaser --version

FROM base AS gomod
ARG TARGETPLATFORM
RUN --mount=type=bind,target=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM gomod AS build
ARG GIT_REF
ARG TARGETPLATFORM
RUN --mount=type=bind,target=/src,rw \
  --mount=target=/go/pkg/mod,type=cache \
  case "$GIT_REF" in \
    refs/tags/v*) gitTag="${GIT_REF#refs/tags/v}" ;; \
    *) gitTag="0.0.0" ;; \
  esac \
  && go build -v -ldflags "-w -s -X 'main.version=${gitTag}'" -o /usr/local/bin/goreleaser-xx \
  && goreleaser-xx --help

FROM scratch AS release
LABEL maintainer="CrazyMax"
COPY --from=goreleaser /src/goreleaser /opt/goreleaser-xx/goreleaser
COPY --from=build /usr/local/bin/goreleaser-xx /usr/bin/goreleaser-xx

FROM golang:${GO_VERSION}-alpine AS test
RUN apk --no-cache add git
COPY --from=release / /
WORKDIR /src
ARG GIT_REF
ARG GORELEASER_VERSION
RUN git clone --branch v${GORELEASER_VERSION} https://github.com/goreleaser/goreleaser .
ARG TARGETPLATFORM
ENV GORELEASER_LDFLAGS="-w -s -X 'main.version=${GORELEASER_VERSION}'"
RUN goreleaser-xx --debug \
  --name goreleaser \
  --dist /out \
  --files "README.md,LICENSE.md" \
  --main .
