// GitHub reference as defined in GitHub Actions (eg. refs/head/master)
variable "GITHUB_REF" {
  default = ""
}

// GoReleaser version
variable "GORELEASER_VERSION" {
  default = "0.162.0"
}

// Go version to build GoReleaser and goreleaser-xx
variable "GO_VERSION" {
  default = "1.16"
}

target "args" {
  args = {
    GIT_REF = GITHUB_REF
    GORELEASER_VERSION = GORELEASER_VERSION
    GO_VERSION = GO_VERSION
  }
}

// Special target: https://github.com/crazy-max/ghaction-docker-meta#bake-definition
target "ghaction-docker-meta" {
  tags = ["goreleaser-xx:local"]
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["args", "ghaction-docker-meta"]
  target = "release"
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
    "linux/386"
  ]
}

target "test" {
  inherits = ["args", "ghaction-docker-meta"]
  target = "test"
}
