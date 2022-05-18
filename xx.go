package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// Target holds Go compiler's target platform
type Target struct {
	Pair string

	Os     string
	Arch   string
	Arm    string
	Mips   string
	Mips64 string
	Amd64  string
}

// Compilers holds C and C++ compilers
type Compilers struct {
	Triple string

	Ar        string
	Cc        string
	Cxx       string
	PkgConfig string

	CgoCflags   string
	CgoCppflags string
	CgoCxxflags string
	CgoFflags   string
	CgoLdflags  string
}

func getTarget() (target Target) {
	var variant string

	if targetPlatform := os.Getenv("TARGETPLATFORM"); targetPlatform != "" {
		stp := strings.Split(targetPlatform, "/")
		target.Os, target.Arch = stp[0], stp[1]
		if len(stp) == 3 {
			variant = stp[2]
		}
	}

	if targetOs := os.Getenv("TARGETOS"); targetOs != "" {
		target.Os = targetOs
	}
	if targetArch := os.Getenv("TARGETARCH"); targetArch != "" {
		target.Arch = targetArch
	}
	if targetVariant := os.Getenv("TARGETVARIANT"); targetVariant != "" {
		variant = targetVariant
	}

	target.Pair = fmt.Sprintf("%s-%s%s", target.Os, target.Arch, variant)

	if target.Arch == "arm" && len(variant) > 0 {
		switch variant {
		case "v5":
			target.Arm = "5"
		case "v6":
			target.Arm = "6"
		default:
			target.Arm = "7"
		}
	}

	if strings.HasPrefix(target.Arch, "mips64") && len(variant) > 0 {
		target.Mips64 = variant
	} else if strings.HasPrefix(target.Arch, "mips") && len(variant) > 0 {
		target.Mips = variant
	}

	if target.Arch == "amd64" && len(variant) > 0 {
		switch variant {
		case "v4":
			target.Amd64 = "v4"
		case "v3":
			target.Amd64 = "v3"
		case "v2":
			target.Amd64 = "v2"
		default:
			target.Amd64 = "v1"
		}
	}

	if target.Os == "wasi" {
		target.Os = "js"
	}

	return
}

func formatTarget(target Target) string {
	if target.Os == "" {
		return "unknown"
	}
	return path.Join(target.Os, target.Arch, target.Arm+target.Mips)
}

func getCompilers(target Target) (cp Compilers) {
	switch target.Arch {
	case "386":
		cp.Triple = "i686-linux-gnu"
		if target.Os == "windows" {
			cp.Triple = "i686-w64-mingw32"
		}
	case "amd64":
		if target.Os != "darwin" {
			cp.Triple = "x86_64-linux-gnu"
			if target.Os == "windows" {
				cp.Triple = "x86_64-w64-mingw32"
			}
		}
	case "arm":
		switch target.Arm {
		case "5":
			cp.Triple = "arm-linux-gnueabi"
			if target.Os == "windows" {
				cp.Triple = "armv5-w64-mingw32"
			}
			cp.CgoCflags = "-march=armv5t"
			cp.CgoCxxflags = "-march=armv5t"
		case "6":
			cp.Triple = "arm-linux-gnueabi"
			if target.Os == "windows" {
				cp.Triple = "armv6-w64-mingw32"
			}
			cp.CgoCflags = "-march=armv6"
			cp.CgoCxxflags = "-march=armv6"
		case "7":
			cp.Triple = "arm-linux-gnueabihf"
			if target.Os == "windows" {
				cp.Triple = "armv7-w64-mingw32"
			}
			cp.CgoCflags = "-march=armv7-a"
			cp.CgoCxxflags = "-march=armv7-a"
		default:
			cp.Triple = "arm-linux-gnueabihf"
		}
	case "arm64":
		if target.Os != "darwin" {
			cp.Triple = "aarch64-linux-gnu"
			if target.Os == "windows" {
				cp.Triple = "aarch64-w64-mingw32"
			}
		}
	case "mips":
		cp.Triple = "mips-linux-gnu"
	case "mipsle":
		cp.Triple = "mipsel-linux-gnu"
	case "mips64":
		cp.Triple = "mips64-linux-gnuabi64"
	case "mips64le":
		cp.Triple = "mips64el-linux-gnuabi64"
	case "ppc64le":
		cp.Triple = "powerpc64le-linux-gnu"
	case "riscv64":
		cp.Triple = "riscv64-linux-gnu"
	case "s390x":
		cp.Triple = "s390x-linux-gnu"
	}

	if target.Os == "darwin" {
		cp.Cc = "o64-clang"
		cp.Cxx = "o64-clang++"
	} else {
		cp.Ar = cp.Triple + "-ar"
		cp.Cc = cp.Triple + "-gcc"
		cp.Cxx = cp.Triple + "-g++"
		cp.PkgConfig = cp.Triple + "-pkg-config"
	}
	if target.Os == "windows" {
		cp.CgoCflags = "-D_WIN32_WINNT=0x0400"
		cp.CgoCxxflags = "-D_WIN32_WINNT=0x0400"
	}

	if v := os.Getenv("AR"); v != "" {
		cp.Ar = v
	}
	if v := os.Getenv("CC"); v != "" {
		cp.Cc = v
	}
	if v := os.Getenv("CXX"); v != "" {
		cp.Cxx = v
	}
	if v := os.Getenv("PKG_CONFIG"); v != "" {
		cp.PkgConfig = v
	}
	if v := os.Getenv("CGO_CFLAGS"); v != "" {
		cp.CgoCflags = v
	}
	if v := os.Getenv("CGO_CPPFLAGS"); v != "" {
		cp.CgoCppflags = v
	}
	if v := os.Getenv("CGO_CXXFLAGS"); v != "" {
		cp.CgoCxxflags = v
	}
	if v := os.Getenv("CGO_FFLAGS"); v != "" {
		cp.CgoFflags = v
	}
	if v := os.Getenv("CGO_LDFLAGS"); v != "" {
		cp.CgoLdflags = v
	}

	return
}
