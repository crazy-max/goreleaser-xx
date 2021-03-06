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
	var variant string

	if targetPlatform := os.Getenv("TARGETPLATFORM"); targetPlatform != "" {
		stp := strings.Split(targetPlatform, "/")
		tgt.Os, tgt.Arch = stp[0], stp[1]
		if len(stp) == 3 {
			variant = stp[2]
		}
	}

	if targetOs := os.Getenv("TARGETOS"); targetOs != "" {
		tgt.Os = targetOs
	}
	if targetArch := os.Getenv("TARGETARCH"); targetArch != "" {
		tgt.Arch = targetArch
	}
	if targetVariant := os.Getenv("TARGETVARIANT"); targetVariant != "" {
		variant = targetVariant
	}

	if tgt.Arch == "arm" && len(variant) > 0 {
		switch variant {
		case "v5":
			tgt.Arm = "5"
		case "v6":
			tgt.Arm = "6"
		default:
			tgt.Arm = "7"
		}
	}

	if strings.HasPrefix(tgt.Arch, "mips") {
		if len(variant) > 0 {
			tgt.Mips = variant
		} else {
			tgt.Mips = "hardfloat"
		}
	}

	if tgt.Os == "wasi" {
		tgt.Os = "js"
	}

	return
}
