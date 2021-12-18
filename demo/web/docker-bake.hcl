group "default" {
  targets = ["image-local"]
}

target "image" {}

target "image-local" {
  inherits = ["image"]
  tags = ["myapp:local"]
  output = ["type=docker"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
    "linux/386",
    "linux/ppc64le"
  ]
}

target "artifact" {
  target = "artifact"
  output = ["./dist"]
}

target "artifact-all" {
  inherits = ["artifact"]
  platforms = [
    "linux/amd64",
    "linux/arm/v5",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
    "linux/386",
    "linux/ppc64le",
    "windows/amd64",
    "windows/386",
    "darwin/amd64",
    "freebsd/amd64",
    "freebsd/386"
  ]
}
