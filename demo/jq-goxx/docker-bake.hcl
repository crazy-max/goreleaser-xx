variable "GORELEASER_XX_BASE" {
  default = "crazymax/goreleaser-xx:edge"
}

target "_commons" {
  args = {
    GORELEASER_XX_BASE = GORELEASER_XX_BASE
  }
}

group "default" {
  targets = ["artifact"]
}

target "artifact" {
  inherits = ["_commons"]
  output = ["./dist"]
}

target "artifact-all" {
  inherits = ["artifact"]
  platforms = [
    "darwin/amd64",
    "darwin/arm64",
    "linux/amd64",
    "linux/arm64",
    "linux/ppc64le",
    "linux/riscv64",
    "linux/s390x"
  ]
}
