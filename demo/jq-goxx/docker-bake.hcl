variable "GOXX_BASE" {
  default = "crazymax/goxx:latest"
}

target "_commons" {
  args = {
    GOXX_BASE = GOXX_BASE
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
