# syntax=docker/dockerfile:1.2

ARG GO_VERSION=1.17
ARG GORELEASER_XX_BASE=crazymax/goreleaser-xx:edge

FROM --platform=$BUILDPLATFORM ${GORELEASER_XX_BASE} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
COPY --from=goreleaser-xx / /
RUN apk add --no-cache git
WORKDIR /src

FROM base AS build
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,target=/src,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --name="xx-hello" \
    --dist="/out" \
    --main="." \
    --ldflags="-s -w" \
    --files="README.md" \
    --pre-hooks="go mod tidy" \
    --pre-hooks="go mod download"

FROM scratch AS artifact
COPY --from=build /out /

FROM scratch
COPY --from=build /usr/local/bin/xx-hello /xx-hello
CMD [ "/xx-hello" ]
