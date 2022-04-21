// GitHub reference as defined in GitHub Actions (eg. refs/head/master)
variable "GITHUB_REF" {
  default = ""
}

// GoReleaser version
variable "GORELEASER_VERSION" {
  default = "1.8.3"
}

// Go version to build GoReleaser and goreleaser-xx
variable "GO_VERSION" {
  default = "1.18"
}

target "_commons" {
  args = {
    GIT_REF = GITHUB_REF
    GORELEASER_VERSION = GORELEASER_VERSION
    GO_VERSION = GO_VERSION
  }
}

// Special target: https://github.com/docker/metadata-action#bake-definition
target "docker-metadata-action" {
  tags = ["goreleaser-xx:local"]
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["_commons", "docker-metadata-action"]
  target = "release"
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/386",
    "linux/amd64",
    "linux/arm/v5",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
    "linux/ppc64le",
    "linux/riscv64",
    "linux/s390x"
  ]
}

target "vendor-update" {
  inherits = ["_commons"]
  target = "vendor-update"
  output = ["."]
}
