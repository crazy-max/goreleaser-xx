package main

import (
	"os"
	"strings"
)

// Target holds Go compiler's target platform
type Target struct {
	Os   string
	Arch string
	Arm  string
	Mips string
}

func getTarget() (tgt Target) {
	if targetPlatform := os.Getenv("TARGETPLATFORM"); targetPlatform != "" {
		stp := strings.Split(targetPlatform, "/")
		goos, goarch := stp[0], stp[1]
		if goos != "" && goarch != "" {
			tgt.Os = goos
			tgt.Arch = goarch
			if goarch == "arm" {
				switch stp[2] {
				case "v5":
					tgt.Arm = "5"
				case "v6":
					tgt.Arm = "6"
				default:
					tgt.Arm = "7"
				}
			}
		}
	}

	if targetOs := os.Getenv("TARGETOS"); targetOs != "" {
		tgt.Os = targetOs
	}

	if targetArch := os.Getenv("TARGETARCH"); targetArch != "" {
		tgt.Arch = targetArch
	}

	if tgt.Arch == "arm" {
		if targetVariant := os.Getenv("TARGETVARIANT"); targetVariant != "" {
			if tgt.Arch == "arm" {
				switch targetVariant {
				case "v5":
					tgt.Arm = "5"
				case "v6":
					tgt.Arm = "6"
				default:
					tgt.Arm = "7"
				}
			}
		} else {
			tgt.Arm = "7"
		}
	}

	if strings.HasPrefix(tgt.Arch, "mips") {
		if targetVariant := os.Getenv("TARGETVARIANT"); targetVariant != "" {
			tgt.Mips = targetVariant
		} else {
			tgt.Mips = "hardfloat"
		}
	}

	if tgt.Os == "wasi" {
		tgt.Os = "js"
	}

	return
}
