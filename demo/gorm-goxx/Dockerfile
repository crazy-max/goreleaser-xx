# syntax=docker/dockerfile:1.3-labs

ARG GO_VERSION=1.17
ARG GORELEASER_XX_BASE=crazymax/goreleaser-xx:edge

FROM --platform=$BUILDPLATFORM ${GORELEASER_XX_BASE} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM crazymax/goxx:${GO_VERSION} AS base
ENV CGO_ENABLED=1
COPY --from=goreleaser-xx / /
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
    --name="gorm-goxx" \
    --dist="/out" \
    --main="." \
    --ldflags="-s -w" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out /

FROM scratch
COPY --from=build /usr/local/bin/gorm-goxx /gorm-goxx
ENTRYPOINT [ "/gorm-goxx" ]