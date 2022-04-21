# syntax=docker/dockerfile:1

ARG GO_VERSION=1.18
ARG GORELEASER_XX_BASE=crazymax/goreleaser-xx:edge

FROM --platform=$BUILDPLATFORM ${GORELEASER_XX_BASE} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
ENV CGO_ENABLED=0
COPY --from=goreleaser-xx / /
RUN apk add --no-cache git
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,target=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored AS build
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --name="echo-webapp" \
    --dist="/out" \
    --main="." \
    --ldflags="-s -w -X 'main.version={{.Version}}' -X 'main.commit={{.Commit}}' -X 'main.date={{.Date}}'" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out /

FROM alpine
RUN apk --update --no-cache add ca-certificates openssl
COPY --from=build /usr/local/bin/echo-webapp /usr/local/bin/echo-webapp
EXPOSE 8080
ENTRYPOINT [ "echo-webapp" ]
