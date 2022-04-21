# syntax=docker/dockerfile:1

ARG GO_VERSION=1.18
ARG GORELEASER_XX_BASE=crazymax/goreleaser-xx:edge

FROM --platform=$BUILDPLATFORM ${GORELEASER_XX_BASE} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
ENV CGO_ENABLED=0
COPY --from=goreleaser-xx / /
RUN apk add --no-cache git
WORKDIR /src

FROM base AS build
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --config=".goreleaser.yml" \
    --name="with-config" \
    --dist="/out" \
    --main="." \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out /

FROM scratch
COPY --from=build /usr/local/bin/with-config /with-config
CMD [ "/with-config" ]
