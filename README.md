[![goreleaser-xx](.github/goreleaser-xx.png)](https://github.com/crazy-max/goreleaser-xx)

[![GitHub release](https://img.shields.io/github/release/crazy-max/goreleaser-xx.svg?style=flat-square)](https://github.com/crazy-max/goreleaser-xx/releases/latest)
[![Build Status](https://img.shields.io/github/workflow/status/crazy-max/goreleaser-xx/build?label=build&logo=github&style=flat-square)](https://github.com/crazy-max/goreleaser-xx/actions?query=workflow%3Abuild)
[![Docker Stars](https://img.shields.io/docker/stars/crazymax/goreleaser-xx.svg?style=flat-square&logo=docker)](https://hub.docker.com/r/crazymax/goreleaser-xx/)
[![Docker Pulls](https://img.shields.io/docker/pulls/crazymax/goreleaser-xx.svg?style=flat-square&logo=docker)](https://hub.docker.com/r/crazymax/goreleaser-xx/)
[![Go Report Card](https://goreportcard.com/badge/github.com/crazy-max/goreleaser-xx)](https://goreportcard.com/report/github.com/crazy-max/goreleaser-xx)

## About

`goreleaser-xx` is a small wrapper around the fantastic [GoReleaser](https://github.com/goreleaser/goreleaser) build
tool to be able to handle a functional [multi-platform scratch Docker image](https://hub.docker.com/r/crazymax/goreleaser-xx/tags?page=1&ordering=last_updated)
to ease the integration and cross compilation in a Dockerfile for your Go projects.

![](.github/goreleaser-xx.gif)
> Building [yasu](https://github.com/crazy-max/yasu) with `goreleaser-xx`

___

* [Features](#features)
* [Projects using goreleaser-xx](#projects-using-goreleaser-xx)
* [Image](#image)
* [CLI](#cli)
* [Usage](#usage)
  * [Minimal](#minimal)
  * [Multi-platform image](#multi-platform-image)
  * [CGo](#cgo)
* [Build](#build)
* [Contributing](#contributing)
* [License](#license)

## Features

* Handle `--platform` in your Dockerfile for multi-platform support
* Build into any architecture
* Translation of [platform ARGs in the global scope](https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope) into Go compiler's target
* Auto generation of `.goreleaser.yml` config based on target architecture
* [Demo projects](demo)

## Projects using goreleaser-xx

* [ddns-route53](https://github.com/crazy-max/ddns-route53)
* [Discord bot](https://github.com/blueprintue/discord-bot)
* [Diun](https://github.com/crazy-max/diun)
* [FTPGrab](https://github.com/crazy-max/ftpgrab)
* [swarm-cronjob](https://github.com/crazy-max/swarm-cronjob)
* [yasu](https://github.com/crazy-max/yasu)
* [Zipper](https://github.com/pratikbalar/zipper)

## Image

| Registry                                                                                                  | Image                                |
|-----------------------------------------------------------------------------------------------------------|--------------------------------------|
| [Docker Hub](https://hub.docker.com/r/crazymax/goreleaser-xx/)                                            | `crazymax/goreleaser-xx`             |
| [GitHub Container Registry](https://github.com/users/crazy-max/packages/container/package/goreleaser-xx)  | `ghcr.io/crazy-max/goreleaser-xx`    |

## CLI

`goreleaser-xx` handles basic [GoReleaser customizations](https://goreleaser.com/customization/) to be able
to generate a minimal `.goreleaser.yml` configuration.

```shell
docker run --rm -t crazymax/goreleaser-xx:latest goreleaser-xx --help
```

| Flag                 | Env var                       | Description   |
|----------------------|-------------------------------|---------------|
| `--debug`            | `DEBUG`                       | Enable debug (default `false`) |
| `--git-ref`          | `GIT_REF`                     | The branch or tag like `refs/tags/v1.0.0` (default to your working tree info) |
| `--goreleaser`       | `GORELEASER_PATH`             | Path to GoReleaser binary (default `/opt/goreleaser-xx/goreleaser`) |
| `--go-binary`        | `GORELEASER_GOBINARY`         | Set a specific go binary to use when building (default `go`) |
| `--name`             | `GORELEASER_NAME`             | Project name |
| `--dist`             | `GORELEASER_DIST`             | Dist folder where artifact will be stored |
| `--artifacts`        | `GORELEASER_ARTIFACTS`        | Types of artifact to create (`archive`, `bin`) (default `archive`) |
| `--main`             | `GORELEASER_MAIN`             | Path to main.go file or main package (default `.`) |
| `--flags`            | `GORELEASER_FLAGS`            | Custom flags templates |
| `--asmflags`         | `GORELEASER_ASMFLAGS`         | Custom asmflags templates |
| `--gcflags`          | `GORELEASER_GCFLAGS`          | Custom gcflags templates |
| `--ldflags`          | `GORELEASER_LDFLAGS`          | Custom ldflags templates |
| `--tags`             | `GORELEASER_TAGS`             | Custom build tags templates |
| `--files`            | `GORELEASER_FILES`            | Additional files/template/globs you want to add to the [archive](https://goreleaser.com/customization/archive/) |
| `--replacements`     | `GORELEASER_REPLACEMENTS`     | Replacements for `GOOS` and `GOARCH` in the archive/binary name |
| `--envs`             | `GORELEASER_ENVS`             | Custom environment variables to be set during the build |
| `--pre-hooks`        | `GORELEASER_PRE_HOOKS`        | [Hooks](https://goreleaser.com/customization/build/#build-hooks) which will be executed before the build |
| `--post-hooks`       | `GORELEASER_POST_HOOKS`       | [Hooks](https://goreleaser.com/customization/build/#build-hooks) which will be executed after the build |
| `--snapshot`         | `GORELEASER_SNAPSHOT`         | Run in [snapshot](https://goreleaser.com/customization/snapshots/) mode |
| `--checksum`         | `GORELEASER_CHECKSUM`         | Create checksum (default `true`) |

## Usage

### Minimal

In the following example we are going to build a simple Go application against `linux/amd64`, `linux/arm64`,
`linux/arm/v7`, `windows/amd64` and `darwin/amd64` platforms using `goreleaser-xx` and
[buildx](https://github.com/docker/buildx).

```Dockerfile
# syntax=docker/dockerfile:1.2

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:alpine AS base
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
    --name="myapp" \
    --dist="/out" \
    --ldflags="-s -w -X 'main.version={{.Version}}'" \
    --files="LICENSE" \
    --files="README.md" \
    --replacements="386=i386" \
    --replacements="amd64=x86_64" \
    --envs="FOO=bar" \
    --envs="BAR=foo"

FROM scratch AS artifact
COPY --from=build /out /
```

Now let's build with buildx:

```shell
docker buildx build \
  --platform "linux/amd64,linux/arm64,linux/arm/v7,windows/amd64,darwin/amd64" \
  --output "type=local,dest=./dist" \
  --target "artifact" \
  --file "./Dockerfile" .
```

Archives created by GoReleaser will be available in `./dist`:

```text
./dist
├── darwin_amd64
│ ├── myapp_v1.0.0-SNAPSHOT-00655a9_darwin_x86_64.tar.gz
│ └── myapp_v1.0.0-SNAPSHOT-00655a9_darwin_x86_64.tar.gz.sha256
├── linux_amd64
│ ├── myapp_v1.0.0-SNAPSHOT-00655a9_linux_x86_64.tar.gz
│ └── myapp_v1.0.0-SNAPSHOT-00655a9_linux_x86_64.tar.gz.sha256
├── linux_arm64
│ ├── myapp_v1.0.0-SNAPSHOT-00655a9_linux_arm64.tar.gz
│ └── myapp_v1.0.0-SNAPSHOT-00655a9_linux_arm64.tar.gz.sha256
├── linux_arm_v7
│ ├── myapp_v1.0.0-SNAPSHOT-00655a9_linux_armv7.tar.gz
│ └── myapp_v1.0.0-SNAPSHOT-00655a9_linux_armv7.tar.gz.sha256
└── windows_amd64
│ ├── myapp_v1.0.0-SNAPSHOT-00655a9_windows_x86_64.tar.gz
│ └── myapp_v1.0.0-SNAPSHOT-00655a9_windows_x86_64.tar.gz.sha256
```

### Multi-platform image

We can enhance the previous example to also create a multi-platform image in addition to the generated artifacts:

```Dockerfile
# syntax=docker/dockerfile:1.2

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:alpine AS base
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
    --name="myapp" \
    --dist="/out" \
    --ldflags="-s -w -X 'main.version={{.Version}}'" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out /

FROM alpine AS image
RUN apk --update --no-cache add ca-certificates openssl shadow \
  && addgroup -g 1000 myapp \
  && adduser -u 1000 -G myapp -s /sbin/nologin -D myapp
COPY --from=build /usr/local/bin/myapp /usr/local/bin/myapp
USER myapp
EXPOSE 8080
ENTRYPOINT [ "myapp" ]
```

As you can see, we have added a new stage called `image`. The artifact of each platform is available with
`goreleaser-xx` in `/usr/local/bin/{name}` (`build` stage) and will be retrieved and included in your `image` stage
through `COPY --from=build` instruction.

Now let's build, tag and push our multi-platform image in our favorite registry with buildx:

```shell
docker buildx build \
  --tag "user/myapp:latest" \
  --platform "linux/amd64,linux/arm64,linux/arm/v7" \
  --target "image" \
  --push \
  --file "./Dockerfile" .
```
> `windows/amd64` and `darwin/amd64` platforms have been removed here
> because `alpine:3.14` does not support them.

### CGo

If you need to use CGo to build your project, you can use `goreleaser-xx` with
[`tonistiigi/xx`](https://github.com/tonistiigi/xx):

```dockerfile
# syntax=docker/dockerfile:1.3-labs

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.1.0 AS xx
FROM --platform=$BUILDPLATFORM golang:alpine AS base
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
    --name="myapp" \
    --dist="/out" \
    --ldflags="-s -w -X 'main.version={{.Version}}'" \
    --envs="CGO_ENABLED=1" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out /

FROM scratch
COPY --from=build /usr/local/bin/myapp /myapp
ENTRYPOINT [ "/myapp" ]
```

## Build

Everything is dockerized and handled by [buildx bake](docker-bake.hcl) for an agnostic usage of this repo:

```shell
git clone https://github.com/crazy-max/goreleaser-xx.git goreleaser-xx
cd goreleaser-xx

# test goreleaser-xx
docker buildx bake test

# build docker image and output to docker with goreleaser-xx:local tag (default)
docker buildx bake

# build multi-platform image
docker buildx bake image-all
```

## Contributing

Want to contribute? Awesome! The most basic way to show your support is to star the project, or to raise issues. If
you want to open a pull request, please read the [contributing guidelines](.github/CONTRIBUTING.md).

You can also support this project by [**becoming a sponsor on GitHub**](https://github.com/sponsors/crazy-max) or by
making a [Paypal donation](https://www.paypal.me/crazyws) to ensure this journey continues indefinitely!

Thanks again for your support, it is much appreciated! :pray:

## License

MIT. See `LICENSE` for more details.
