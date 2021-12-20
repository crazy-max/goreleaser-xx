# syntax=docker/dockerfile:1.2

ARG GO_VERSION=1.17
ARG GORELEASER_XX_IMAGE=crazymax/goreleaser-xx:edge

FROM --platform=$BUILDPLATFORM ${GORELEASER_XX_IMAGE} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
COPY --from=goreleaser-xx / /
RUN apk add --no-cache git
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,target=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored AS build
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,target=/src,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --name="xx-echo" \
    --dist="/out" \
    --main="." \
    --ldflags="-s -w -X 'main.version={{.Version}}' -X 'main.commit={{.Commit}}' -X 'main.date={{.Date}}'" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /

FROM alpine

RUN apk --update --no-cache add \
    ca-certificates \
    libressl \
    shadow \
  && addgroup -g 1000 xx-echo \
  && adduser -u 1000 -G xx-echo -s /sbin/nologin -D xx-echo

COPY --from=build /usr/local/bin/xx-echo /usr/local/bin/xx-echo

USER xx-echo
EXPOSE 8080
ENTRYPOINT [ "xx-echo" ]
