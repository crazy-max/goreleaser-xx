[![goreleaser-xx](.github/goreleaser-xx.png)](https://github.com/crazy-max/goreleaser-xx)

[![GitHub release](https://img.shields.io/github/release/crazy-max/goreleaser-xx.svg?style=flat-square)](https://github.com/crazy-max/goreleaser-xx/releases/latest)
[![Build Status](https://img.shields.io/github/workflow/status/crazy-max/goreleaser-xx/build?label=build&logo=github&style=flat-square)](https://github.com/crazy-max/goreleaser-xx/actions?query=workflow%3Abuild)
[![Docker Stars](https://img.shields.io/docker/stars/crazymax/goreleaser-xx.svg?style=flat-square&logo=docker)](https://hub.docker.com/r/crazymax/goreleaser-xx/)
[![Docker Pulls](https://img.shields.io/docker/pulls/crazymax/goreleaser-xx.svg?style=flat-square&logo=docker)](https://hub.docker.com/r/crazymax/goreleaser-xx/)
[![Go Report Card](https://goreportcard.com/badge/github.com/crazy-max/goreleaser-xx)](https://goreportcard.com/report/github.com/crazy-max/goreleaser-xx)

## About

`goreleaser-xx` is a small CLI wrapper for [GoReleaser](https://github.com/goreleaser/goreleaser)
and available as a [lightweight and multi-platform scratch Docker image](https://hub.docker.com/r/crazymax/goreleaser-xx/tags?page=1&ordering=last_updated)
to ease the integration and cross compilation in a Dockerfile for your Go
projects.

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
  * [With `.goreleaser.yml`](#with-goreleaseryml)
  * [CGO](#cgo)
    * [`crazy-max/goxx`](#crazy-maxgoxx)
    * [`tonistiigi/xx`](#tonistiigixx)
* [Notes](#notes)
  * [`CGO_ENABLED`](#cgo_enabled)
* [Build](#build)
* [Contributing](#contributing)
* [License](#license)

## Features

* Handle `--platform` in your Dockerfile for multi-platform support
* Build into any architecture
* Handle C and C++ compilers ([CGO](#cgo))
* Translation of [platform ARGs in the global scope](https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope) into Go compiler's target
* Auto generation of `.goreleaser.yml` config based on target architecture
* [Demo projects](demo)

## Projects using goreleaser-xx

* [ddns-route53](https://github.com/crazy-max/ddns-route53): Dynamic DNS for Amazon Route 53
* [Diun](https://github.com/crazy-max/diun): Docker Image Update Notifier
* [Distribution](https://github.com/distribution/distribution): The toolkit to pack, ship, store, and deliver container content (formerly docker/registry)
* [faq](https://github.com/jzelinskie/faq): Format Agnostic jQ
* [FTPGrab](https://github.com/crazy-max/ftpgrab): Grab your files periodically from a remote FTP or SFTP server easily
* [swarm-cronjob](https://github.com/crazy-max/swarm-cronjob): Create jobs on a time-based schedule on Docker Swarm
* [tarrer](https://github.com/pratikbalar/tarrer): Dumb br, bz2, zip, gz, lz4, sz, xz, zstd extractor
* [yasu](https://github.com/crazy-max/yasu): Yet Another Switch User

## Image

| Registry                                                                                                  | Image                                |
|-----------------------------------------------------------------------------------------------------------|--------------------------------------|
| [Docker Hub](https://hub.docker.com/r/crazymax/goreleaser-xx/)                                            | `crazymax/goreleaser-xx`             |
| [GitHub Container Registry](https://github.com/users/crazy-max/packages/container/package/goreleaser-xx)  | `ghcr.io/crazy-max/goreleaser-xx`    |

Following platforms for this image are available:

```
$ docker run --rm mplatform/mquery crazymax/goreleaser-xx:latest
Image: crazymax/goreleaser-xx:latest (digest: sha256:c65c481c014abab6d307b190ddf1fcb229a44b6c1845d2f2a53bd06dc0437cd7)
 * Manifest List: Yes (Image type: application/vnd.docker.distribution.manifest.list.v2+json)
 * Supported platforms:
   - linux/386
   - linux/amd64
   - linux/arm/v5
   - linux/arm/v6
   - linux/arm/v7
   - linux/arm64
   - linux/ppc64le
   - linux/riscv64
   - linux/s390x
```

## CLI

```shell
docker run --rm -t crazymax/goreleaser-xx:latest goreleaser-xx --help
```

| Flag                 | Env var                       | Description                                                                                                     |
|----------------------|-------------------------------|-----------------------------------------------------------------------------------------------------------------|
| `--debug`            | `DEBUG`                       | Enable debug (default `false`)                                                                                  |
| `--git-ref`          | `GIT_REF`                     | The branch or tag like `refs/tags/v1.0.0` (default to your working tree info)                                   |
| `--goreleaser`       | `GORELEASER_PATH`             | Set a specific GoReleaser binary to use (default `goreleaser`)                                                  |
| `--config`           | `GORELEASER_CONFIG`           | Load GoReleaser configuration from file                                                                         |
| `--go-binary`        | `GORELEASER_GOBINARY`         | Set a specific go binary to use when building (default `go`)                                                    |
| `--name`             | `GORELEASER_NAME`             | Project name                                                                                                    |
| `--dist`             | `GORELEASER_DIST`             | Dist folder where artifact will be stored                                                                       |
| `--artifacts`        | `GORELEASER_ARTIFACTS`        | Types of artifact to create (`archive`, `bin`) (default `archive`)                                              |
| `--main`             | `GORELEASER_MAIN`             | Path to main.go file or main package (default `.`)                                                              |
| `--flags`            | `GORELEASER_FLAGS`            | Custom flags templates                                                                                          |
| `--asmflags`         | `GORELEASER_ASMFLAGS`         | Custom asmflags templates                                                                                       |
| `--gcflags`          | `GORELEASER_GCFLAGS`          | Custom gcflags templates                                                                                        |
| `--ldflags`          | `GORELEASER_LDFLAGS`          | Custom ldflags templates                                                                                        |
| `--tags`             | `GORELEASER_TAGS`             | Custom build tags templates                                                                                     |
| `--files`            | `GORELEASER_FILES`            | Additional files/template/globs you want to add to the [archive](https://goreleaser.com/customization/archive/) |
| `--replacements`     | `GORELEASER_REPLACEMENTS`     | Replacements for `GOOS` and `GOARCH` in the archive/binary name                                                 |
| `--envs`             | `GORELEASER_ENVS`             | Custom environment variables to be set during the build                                                         |
| `--pre-hooks`        | `GORELEASER_PRE_HOOKS`        | [Hooks](https://goreleaser.com/customization/build/#build-hooks) which will be executed before the build        |
| `--post-hooks`       | `GORELEASER_POST_HOOKS`       | [Hooks](https://goreleaser.com/customization/build/#build-hooks) which will be executed after the build         |
| `--snapshot`         | `GORELEASER_SNAPSHOT`         | Run in [snapshot](https://goreleaser.com/customization/snapshots/) mode                                         |
| `--checksum`         | `GORELEASER_CHECKSUM`         | Create checksum (default `true`)                                                                                |

## Usage

In order to use it, we will use the `docker buildx` command in the following
examples. [Buildx](https://github.com/docker/buildx) is a Docker component that
enables many powerful build features. All builds executed via buildx run with
[Moby BuildKit](https://github.com/moby/buildkit) builder engine.

### Minimal

Here is a minimal Dockerfile to build a Go project using `goreleaser-xx`:

```Dockerfile
# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:1.17-alpine AS base
ENV CGO_ENABLED=0
COPY --from=goreleaser-xx / /
RUN apk add --no-cache git
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored AS build
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --name="myapp" \
    --dist="/out" \
    --ldflags="-s -w -X 'main.version={{.Version}}'" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /
```

* `FROM --platform=$BUILDPLATFORM ...` command will pull an image that will
  always match the native platform of your machine (e.g., `linux/amd64`). 
  `BUILDPLATFORM` is part of the [ARGs in the global scope](https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope).
* `ARG TARGETPLATFORM` is also an ARG in the global scope that will be set
  to the platform of the target that will default to your current platform or
  can be defined via the [`--platform` flag](https://docs.docker.com/engine/reference/commandline/buildx_build/#platform)
  of buildx so `goreleaser-xx` will be able to automatically build against
  the right platform.

> More details about multi-platform builds in this [blog post](https://medium.com/@tonistiigi/faster-multi-platform-builds-dockerfile-cross-compilation-guide-part-1-ec087c719eaf).

As you can see [`goreleaser-xx` CLI](#cli) handles basic [GoReleaser build customizations](https://goreleaser.com/customization/build/)
with flags to be able to generate a temp and dynamic `.goreleaser.yml` configuration,
but you can also include your own [GoReleaser YAML config](#with-goreleaseryml).

Let's run a simple build against the `artifact` target in our Dockerfile:

```shell
# build and output content of the artifact stage that contains the archive in ./dist
docker buildx build \
  --output "./dist" \
  --target "artifact" .

$ tree ./dist
./dist
├── myapp_v1.0.0-SNAPSHOT-00655a9_linux_amd64.tar.gz
└── myapp_v1.0.0-SNAPSHOT-00655a9_linux_amd64.tar.gz.sha256
```

Here `linux/amd64` arch is used because it's my current platform. If we want
to handle more platforms, we need to create a builder instance as building
multi-platform is currently only supported when using BuildKit with the
[`docker-container` or `kubernetes` drivers](https://docs.docker.com/engine/reference/commandline/buildx_create/#driver).

```shell
# create a builder instance
$ docker buildx create --name "mybuilder" --use

# now build for other platforms
$ docker buildx build \
  --platform "linux/amd64,linux/arm64,linux/arm/v7,windows/amd64,darwin/amd64" \
  --output "./dist" \
  --target "artifact" .

$ tree ./dist
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

We can enhance the previous example to also create a multi-platform image in
addition to the generated artifacts:

```Dockerfile
# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:1.17-alpine AS base
ENV CGO_ENABLED=0
COPY --from=goreleaser-xx / /
RUN apk add --no-cache git
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored AS build
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --name="myapp" \
    --dist="/out" \
    --ldflags="-s -w -X 'main.version={{.Version}}'" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /

FROM alpine AS image
RUN apk --update --no-cache add ca-certificates openssl
COPY --from=build /usr/local/bin/myapp /usr/local/bin/myapp
EXPOSE 8080
ENTRYPOINT [ "myapp" ]
```

As you can see, we have added a new stage called `image`. The artifact of each
platform is available with `goreleaser-xx` in `/usr/local/bin/{{ .ProjectName }}{{ .Ext }}`
(`build` stage) and will be included in your `image` stage via `COPY --from=build`
command.

Now let's build, tag and push our multi-platform image with buildx:

```shell
docker buildx build \
  --tag "user/myapp:latest" \
  --platform "linux/amd64,linux/arm64,linux/arm/v7" \
  --target "image" \
  --push .
```

> `windows/amd64` and `darwin/amd64` platforms have been removed here
> because `alpine:3.14` does not support them.

### With `.goreleaser.yml`

You can also use a `.goreleaser.yml` to configure your build:

```yaml
env:
  - GO111MODULE=auto

gomod:
  proxy: true

builds:
  - mod_timestamp: '{{ .CommitTimestamp }}'
    flags:
      - -trimpath
    ldflags:
      - -s -w

nfpms:
  - file_name_template: '{{ .ConventionalFileName }}'
    homepage:  https://github.com/user/hello
    description: Hello world
    maintainer: Hello <hello@world.com>
    license: MIT
    vendor: HelloWorld
    formats:
      - apk
      - deb
      - rpm
```

```Dockerfile
# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:1.17-alpine AS base
ENV CGO_ENABLED=0
COPY --from=goreleaser-xx / /
RUN apk add --no-cache git
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored AS build
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --config=".goreleaser.yml" \
    --name="hello" \
    --dist="/out" \
    --main="." \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /
```

### CGO

Here are some examples to use CGO to build your project with `goreleaser-xx`:

#### `crazy-max/goxx`

> https://github.com/crazy-max/goxx

```dockerfile
# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM crazymax/osxcross:11.3 AS osxcross
FROM --platform=$BUILDPLATFORM crazymax/goxx:1.17 AS base
COPY --from=osxcross /osxcross /osxcross
COPY --from=goreleaser-xx / /
ENV CGO_ENABLED=1
RUN goxx-apt-get install --no-install-recommends -y git
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored AS build
ARG TARGETPLATFORM
RUN --mount=type=cache,sharing=private,target=/var/cache/apt \
  --mount=type=cache,sharing=private,target=/var/lib/apt/lists \
  goxx-apt-get install -y binutils gcc g++ pkg-config
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --name="myapp" \
    --dist="/out" \
    --ldflags="-s -w -X 'main.version={{.Version}}'" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /
```

#### `tonistiigi/xx`

> https://github.com/tonistiigi/xx

```dockerfile
# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.1.0 AS xx
FROM --platform=$BUILDPLATFORM golang:1.17-alpine AS base
ENV CGO_ENABLED=1
COPY --from=goreleaser-xx / /
COPY --from=xx / /
RUN apk add --no-cache clang git file lld llvm pkgconfig
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,source=.,target=/src,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored AS build
ARG TARGETPLATFORM
RUN xx-apk add --no-cache gcc musl-dev
# XX_CC_PREFER_STATIC_LINKER prefers ld to lld in ppc64le and 386.
ENV XX_CC_PREFER_STATIC_LINKER=1
RUN --mount=type=bind,source=.,rw \
  --mount=from=crazymax/osxcross:11.3,src=/osxsdk,target=/xx-sdk \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --go-binary="xx-go" \
    --name="myapp" \
    --dist="/out" \
    --ldflags="-s -w -X 'main.version={{.Version}}'" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /
```

## Notes

### `CGO_ENABLED`

By default, CGO is enabled in Go when compiling for native architecture and
disabled when cross-compiling. It's therefore recommended to always set
`CGO_ENABLED=0` or `CGO_ENABLED=1` when cross-compiling depending on whether
you need to use CGO or not.

## Build

Everything is dockerized and handled by [buildx bake](docker-bake.hcl) for an
agnostic usage of this repo:

```shell
git clone https://github.com/crazy-max/goreleaser-xx.git goreleaser-xx
cd goreleaser-xx

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
