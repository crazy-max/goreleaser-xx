package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type GoReleaserConfig struct {
	Path    string
	Project config.Project
}

func getConfig(cli Cli, target Target, compilers Compilers) (grc GoReleaserConfig, _ error) {
	var cfg config.Project
	var err error

	if len(cli.Config) > 0 {
		b, err := os.ReadFile(cli.Config)
		if err != nil {
			return grc, err
		}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return grc, err
		}
	}

	if len(cfg.Dist) > 0 {
		log.Printf("WARN: dist specified in your config file is overrided")
	}
	cfg.Dist, err = os.MkdirTemp(os.TempDir(), "dist")
	if err != nil {
		return grc, err
	}
	if len(cli.Name) > 0 {
		cfg.ProjectName = cli.Name
	}

	var build config.Build
	if len(cfg.Builds) == 1 {
		build = cfg.Builds[0]
	} else if len(cfg.Builds) > 1 {
		return grc, errors.New("multiple builds found, please specify one")
	}

	cgoEnabled := false
	if ok := hasCgoEnabled(os.Environ()); ok {
		cgoEnabled = true
	}
	if ok := hasCgoEnabled(cli.Envs); ok {
		cgoEnabled = true
	}
	if ok := hasCgoEnabled(build.Env); ok {
		cgoEnabled = true
	}
	if cgoEnabled {
		if len(compilers.Ar) > 0 {
			if _, err = exec.LookPath(compilers.Ar); err == nil {
				cfg.Env = append(cfg.Env, "AR="+compilers.Ar)
			}
		}
		if len(compilers.Cc) > 0 {
			if _, err = exec.LookPath(compilers.Cc); err == nil {
				cfg.Env = append(cfg.Env, "CC="+compilers.Cc)
				if len(compilers.CgoCflags) > 0 {
					cfg.Env = append(cfg.Env, "CGO_CFLAGS="+compilers.CgoCflags)
				}
			}
		}
		if len(compilers.Cxx) > 0 {
			if _, err = exec.LookPath(compilers.Cxx); err == nil {
				cfg.Env = append(cfg.Env, "CXX="+compilers.Cxx)
				if len(compilers.CgoCxxflags) > 0 {
					cfg.Env = append(cfg.Env, "CGO_CXXFLAGS="+compilers.CgoCxxflags)
				}
			}
		}
	}
	if len(cli.Envs) > 0 {
		cfg.Env = append(cfg.Env, cli.Envs...)
	}
	if len(build.Env) > 0 {
		cfg.Env = append(cfg.Env, build.Env...)
		build.Env = nil
	}

	if len(build.Goos) > 0 {
		log.Printf("WARN: goos specified in your config file is overrided")
	}
	if len(build.Goarch) > 0 {
		log.Printf("WARN: goarch specified in your config file is overrided")
	}
	if len(build.Goarm) > 0 {
		log.Printf("WARN: goarm specified in your config file is overrided")
	}
	if len(build.Gomips) > 0 {
		log.Printf("WARN: gomips specified in your config file is overrided")
	}
	build.Goos = []string{target.Os}
	build.Goarch = []string{target.Arch}
	build.Goarm = []string{target.Arm}
	build.Gomips = []string{target.Mips}

	if len(cli.Main) > 0 {
		build.Main = cli.Main
	}
	if len(cli.Flags) > 0 {
		build.Flags = append(build.Flags, cli.Flags)
	}
	if len(cli.Asmflags) > 0 {
		build.Asmflags = append(build.Asmflags, cli.Asmflags)
	}
	if len(cli.Gcflags) > 0 {
		build.Gcflags = append(build.Gcflags, cli.Gcflags)
	}
	if len(cli.Ldflags) > 0 {
		build.Ldflags = append(build.Ldflags, cli.Ldflags)
	}
	if len(cli.Tags) > 0 {
		build.Tags = append(build.Tags, cli.Tags)
	}
	if len(cli.GoBinary) > 0 {
		build.GoBinary = cli.GoBinary
	}

	if len(cfg.Before.Hooks) > 0 {
		for _, cmd := range cfg.Before.Hooks {
			build.Hooks.Pre = append(build.Hooks.Pre, config.Hook{
				Cmd: cmd,
			})
		}
		cfg.Before.Hooks = nil
	}
	for _, cmd := range cli.PreHooks {
		build.Hooks.Pre = append(build.Hooks.Pre, config.Hook{
			Cmd: cmd,
		})
	}

	build.Hooks.Post = append(build.Hooks.Post, config.Hook{
		Cmd: `cp "{{ .Path }}" "/usr/local/bin/{{ .ProjectName }}{{ .Ext }}"`,
	})
	for _, cmd := range cli.PostHooks {
		build.Hooks.Post = append(build.Hooks.Post, config.Hook{
			Cmd: cmd,
		})
	}

	var archive config.Archive
	if len(cfg.Archives) == 1 {
		archive = cfg.Archives[0]
	} else if len(cfg.Archives) > 1 {
		return grc, errors.New("multiple archives found, please specify one")
	}

	for _, f := range cli.Files {
		archive.Files = append(archive.Files, config.File{
			Source: f,
		})
	}
	if archive.Replacements == nil {
		archive.Replacements = make(map[string]string)
	}
	for k, v := range cli.Replacements {
		archive.Replacements[k] = v
	}

	winOverride := false
	for _, o := range archive.FormatOverrides {
		if o.Goos == "windows" {
			winOverride = true
			break
		}
	}
	if !winOverride {
		archive.FormatOverrides = append(archive.FormatOverrides, config.FormatOverride{
			Goos:   "windows",
			Format: "zip",
		})
	}

	cfg.Builds = []config.Build{build}
	cfg.Archives = []config.Archive{archive}

	if !reflect.DeepEqual(cfg.Checksum, config.Checksum{}) {
		log.Printf("WARN: checksum section specified in your config file is disabled")
	}
	cfg.Checksum = config.Checksum{Disable: true}

	if !reflect.DeepEqual(cfg.Release, config.Release{}) {
		log.Printf("WARN: release section specified in your config file disabled")
	}
	cfg.Release = config.Release{Disable: true}

	if !reflect.DeepEqual(cfg.Changelog, config.Changelog{}) {
		log.Printf("WARN: changelog section specified in your config file skipped")
	}
	cfg.Changelog = config.Changelog{Skip: true}

	if len(cfg.NFPMs) > 0 && !allowNfpms(target) {
		log.Printf("WARN: nfpms section specified in your config file disabled for %s", formatTarget(target))
		cfg.NFPMs = nil
	}

	if len(cfg.Brews) > 0 && !allowBrews(target) {
		log.Printf("WARN: brews section specified in your config file disabled for %s", formatTarget(target))
		cfg.Brews = nil
	}

	if len(cfg.Snapcrafts) > 0 && !allowSnaps(target) {
		log.Printf("WARN: snapcrafts section specified in your config file disabled for %s", formatTarget(target))
		cfg.Snapcrafts = nil
	}

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return grc, err
	}

	f, err := os.CreateTemp(os.TempDir(), ".goreleaser.yml")
	if err != nil {
		return grc, err
	}
	defer f.Close()
	if err := os.WriteFile(f.Name(), b, 0644); err != nil {
		defer os.Remove(f.Name())
		return grc, err
	}

	if cli.Debug {
		log.Println("DBG: goreleaser config:")
		scanner := bufio.NewScanner(strings.NewReader(string(b)))
		for scanner.Scan() {
			log.Printf("  %s", scanner.Text())
		}
	}

	return GoReleaserConfig{
		Path:    f.Name(),
		Project: cfg,
	}, nil
}

func allowNfpms(target Target) bool {
	return target.Os == "linux"
}

func allowBrews(target Target) bool {
	if target.Os != "linux" && target.Os != "darwin" {
		return false
	}
	for _, a := range []string{"arm64", "arm", "amd64"} {
		if target.Arch == a {
			return true
		}
	}
	return false
}

func allowSnaps(target Target) bool {
	if target.Os != "linux" {
		return false
	}
	for _, a := range []string{"s390x", "ppc64le", "arm64", "arm", "amd64", "386"} {
		if target.Arch == a {
			return true
		}
	}
	return false
}

func hasCgoEnabled(envs []string) bool {
	for _, e := range envs {
		if e == "CGO_ENABLED=1" {
			return true
		}
	}
	return false
}
