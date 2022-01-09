# syntax=docker/dockerfile:1-labs

ARG GO_VERSION=1.17
ARG GORELEASER_XX_BASE=crazymax/goreleaser-xx:edge
ARG JQ_VERSION="jq-1.6"

FROM --platform=$BUILDPLATFORM ${GORELEASER_XX_BASE} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM crazymax/osxcross:11.3 AS osxcross
FROM --platform=$BUILDPLATFORM crazymax/goxx:${GO_VERSION} AS base
COPY --from=osxcross /osxcross /osxcross
COPY --from=goreleaser-xx / /
RUN goxx-apt-get install --no-install-recommends -y git
ENV GO111MODULE=auto
ENV CGO_ENABLED=1
ENV OSXCROSS_MP_INC=1
WORKDIR /src

FROM base AS libjq-linux
RUN apt-get install -y autoconf automake flex libtool
ARG TARGETPLATFORM
RUN --mount=type=cache,sharing=private,target=/var/cache/apt \
  --mount=type=cache,sharing=private,target=/var/lib/apt/lists \
  goxx-apt-get install -y binutils gcc pkg-config
WORKDIR /usr/local/src/jq
ARG JQ_VERSION
RUN <<EOT
set -e
git clone --depth 1 --recurse-submodules --shallow-submodules -b $JQ_VERSION https://github.com/stedolan/jq.git .
HOST_TRIPLE=$(. goxx-env && echo $GOXX_TRIPLE)
BUILD_TRIPLE=$(TARGETPLATFORM= . goxx-env && echo $GOXX_TRIPLE)
autoreconf -fi
CC="$HOST_TRIPLE-gcc" ./configure \
  --prefix=/usr/$HOST_TRIPLE \
  --host=$HOST_TRIPLE \
  --build=$BUILD_TRIPLE \
  --target=$BUILD_TRIPLE \
  --disable-maintainer-mode \
  --disable-docs \
  --enable-all-static \
  --with-oniguruma
make
make -j$(nproc) LDFLAGS=-all-static install DESTDIR="/out"
EOT

FROM scratch AS libjq-dummy
WORKDIR /out

FROM libjq-dummy AS libjq-windows
FROM libjq-dummy AS libjq-darwin
FROM libjq-${TARGETOS} AS libjq

FROM base AS build
COPY --from=libjq /out /
ARG TARGETPLATFORM
RUN --mount=type=cache,sharing=private,target=/var/cache/apt \
  --mount=type=cache,sharing=private,target=/var/lib/apt/lists \
  goxx-apt-get install -y binutils gcc pkg-config
RUN goxx-macports --static install jq
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache <<EOT
set -e
EXTLDFLAGS="-v"
if [ "$(. goxx-env && echo $GOOS)" = "linux" ]; then
  EXTLDFLAGS="$EXTLDFLAGS -static"
  export CGO_CFLAGS="-lm -g -O2"
  export CGO_LDFLAGS="-lm -g -O2"
fi
goreleaser-xx --debug \
  --config=".goreleaser.yml" \
  --name="jq-static-goxx" \
  --dist="/out" \
  --main="." \
  --ldflags="-s -w -extldflags '$EXTLDFLAGS'" \
  --tags="netgo" \
  --files="README.md"
EOT

FROM scratch AS artifact
COPY --from=build /out /

FROM scratch
COPY --from=build /usr/local/bin/jq-static-goxx /jq-static-goxx
ENTRYPOINT [ "/jq-static-goxx" ]
