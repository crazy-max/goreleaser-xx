# syntax=docker/dockerfile:1.3

ARG GORELEASER_VERSION
ARG GO_VERSION

FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.0.0 AS xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
COPY --from=xx / /
RUN apk --update --no-cache add build-base git
ARG GORELEASER_VERSION
WORKDIR /goreleaser
RUN git clone --branch v${GORELEASER_VERSION} https://github.com/goreleaser/goreleaser .
ENV CGO_ENABLED=0
WORKDIR /src

FROM base AS goreleaser
ARG TARGETPLATFORM
WORKDIR /goreleaser
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download \
  && xx-go build -v -ldflags "-w -s -X 'main.version=${GORELEASER_VERSION}' -X main.builtBy=goreleaser-xx" \
  && ./goreleaser --version

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
RUN --mount=type=bind,target=/src,rw \
  --mount=target=/go/pkg/mod,type=cache \
  case "$GIT_REF" in \
    refs/tags/v*) gitTag="${GIT_REF#refs/tags/v}" ;; \
    *) gitTag="0.0.0" ;; \
  esac \
  && xx-go build -v -ldflags "-w -s -X 'main.version=${gitTag}'" -o /usr/local/bin/goreleaser-xx \
  && goreleaser-xx --help \
  && goreleaser-xx --version

FROM scratch AS release
COPY --from=goreleaser /goreleaser/goreleaser /opt/goreleaser-xx/goreleaser
COPY --from=build /usr/local/bin/goreleaser-xx /usr/bin/goreleaser-xx

FROM --platform=$BUILDPLATFORM golang:1.17-alpine AS test
RUN apk --no-cache add git
WORKDIR /src
ARG GIT_REF
RUN git clone --branch v2.7.0 https://github.com/crazy-max/ddns-route53 .
COPY --from=release / /
ARG TARGETPLATFORM
RUN goreleaser-xx --debug \
    --name="ddns-route53" \
    --dist="/dist" \
    --hooks="go mod tidy" \
    --hooks="go mod download" \
    --main="./cmd/main.go" \
    --ldflags="-s -w -X 'main.version={{.Version}}'" \
    --files="CHANGELOG.md" \
    --files="LICENSE" \
    --files="README.md" \
    --replacements="386=i386" \
    --replacements="amd64=x86_64" \
  && ls -al /dist/

FROM scratch AS test-artifact
COPY --from=test /dist /

FROM alpine AS test-image
COPY --from=test /usr/local/bin/ddns-route53 /usr/local/bin/ddns-route53
RUN ddns-route53 --version
