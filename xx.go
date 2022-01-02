package main

import (
	"os"
	"path"
	"strings"
)

// Target holds Go compiler's target platform
type Target struct {
	Os   string
	Arch string
	Arm  string
	Mips string
}

// Compilers holds C and C++ compilers
type Compilers struct {
	Ar          string
	Cc          string
	Cxx         string
	CgoCflags   string
	CgoCxxflags string
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

func formatTarget(target Target) string {
	if target.Os == "" {
		return "unknown"
	}
	return path.Join(target.Os, target.Arch, target.Arm+target.Mips)
}

func getCompilers(t Target) (cp Compilers) {
	switch t.Arch {
	case "386":
		cp.Ar = "i686-linux-gnu-ar"
		cp.Cc = "i686-linux-gnu-gcc"
		cp.Cxx = "i686-linux-gnu-g++"
		if t.Os == "windows" {
			cp.Ar = "i686-w64-mingw32-ar"
			cp.Cc = "i686-w64-mingw32-gcc"
			cp.Cxx = "i686-w64-mingw32-g++"
			cp.CgoCflags = "-D_WIN32_WINNT=0x0400"
			cp.CgoCxxflags = "-D_WIN32_WINNT=0x0400"
		}
	case "amd64":
		if t.Os == "darwin" {
			cp.Cc = "o64-clang"
			cp.Cxx = "o64-clang++"
		} else if t.Os == "windows" {
			cp.Ar = "x86_64-w64-mingw32-ar"
			cp.Cc = "x86_64-w64-mingw32-gcc"
			cp.Cxx = "x86_64-w64-mingw32-g++"
			cp.CgoCflags = "-D_WIN32_WINNT=0x0400"
			cp.CgoCxxflags = "-D_WIN32_WINNT=0x0400"
		} else {
			cp.Ar = "x86_64-linux-gnu-ar"
			cp.Cc = "x86_64-linux-gnu-gcc"
			cp.Cxx = "x86_64-linux-gnu-g++"
		}
	case "arm":
		switch t.Arm {
		case "5":
			cp.Ar = "arm-linux-gnueabi-ar"
			cp.Cc = "arm-linux-gnueabi-gcc"
			cp.Cxx = "arm-linux-gnueabi-g++"
			cp.CgoCflags = "-march=armv5t"
			cp.CgoCxxflags = "-march=armv5t"
		case "6":
			cp.Ar = "arm-linux-gnueabi-ar"
			cp.Cc = "arm-linux-gnueabi-gcc"
			cp.Cxx = "arm-linux-gnueabi-g++"
			cp.CgoCflags = "-march=armv6"
			cp.CgoCxxflags = "-march=armv6"
		case "7":
			cp.Ar = "arm-linux-gnueabihf-ar"
			cp.Cc = "arm-linux-gnueabihf-gcc"
			cp.Cxx = "arm-linux-gnueabihf-g++"
			cp.CgoCflags = "-march=armv7-a"
			cp.CgoCxxflags = "-march=armv7-a"
		default:
			cp.Ar = "arm-linux-gnueabihf-ar"
			cp.Cc = "arm-linux-gnueabihf-gcc"
			cp.Cxx = "arm-linux-gnueabihf-g++"
		}
	case "arm64":
		if t.Os == "darwin" {
			cp.Cc = "o64-clang"
			cp.Cxx = "o64-clang++"
		} else {
			cp.Ar = "aarch64-linux-gnu-ar"
			cp.Cc = "aarch64-linux-gnu-gcc"
			cp.Cxx = "aarch64-linux-gnu-g++"
		}
	case "mips":
		cp.Ar = "mips-linux-gnu-ar"
		cp.Cc = "mips-linux-gnu-gcc"
		cp.Cxx = "mips-linux-gnu-g++"
	case "mipsle":
		cp.Ar = "mipsel-linux-gnu-ar"
		cp.Cc = "mipsel-linux-gnu-gcc"
		cp.Cxx = "mipsel-linux-gnu-g++"
	case "mips64":
		cp.Ar = "mips64-linux-gnuabi64-ar"
		cp.Cc = "mips64-linux-gnuabi64-gcc"
		cp.Cxx = "mips64-linux-gnuabi64-g++"
	case "mips64le":
		cp.Ar = "mips64el-linux-gnuabi64-ar"
		cp.Cc = "mips64el-linux-gnuabi64-gcc"
		cp.Cxx = "mips64el-linux-gnuabi64-g++"
	case "ppc64le":
		cp.Ar = "powerpc64le-linux-gnu-ar"
		cp.Cc = "powerpc64le-linux-gnu-gcc"
		cp.Cxx = "powerpc64le-linux-gnu-g++"
	case "riscv64":
		cp.Ar = "riscv64-linux-gnu-ar"
		cp.Cc = "riscv64-linux-gnu-gcc"
		cp.Cxx = "riscv64-linux-gnu-g++"
	case "s390x":
		cp.Ar = "s390x-linux-gnu-ar"
		cp.Cc = "s390x-linux-gnu-gcc"
		cp.Cxx = "s390x-linux-gnu-g++"
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
	if v := os.Getenv("CGO_CFLAGS"); v != "" {
		cp.CgoCflags = v
	}
	if v := os.Getenv("CGO_CXXFLAGS"); v != "" {
		cp.CgoCxxflags = v
	}

	return
}
