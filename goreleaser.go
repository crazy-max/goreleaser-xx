package main

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func getGRConfig(cli Cli, target Target) (string, string, error) {
	var cfg config.Project
	var err error

	if len(cli.Config) > 0 {
		b, err := os.ReadFile(cli.Config)
		if err != nil {
			return "", "", err
		}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return "", "", err
		}
	}

	cfg.Dist, err = os.MkdirTemp(os.TempDir(), "dist")
	if err != nil {
		return "", "", err
	}
	if len(cli.Name) > 0 {
		cfg.ProjectName = cli.Name
	}
	if len(cli.Envs) > 0 {
		cfg.Env = append(cfg.Env, cli.Envs...)
	}

	var build config.Build
	if len(cfg.Builds) == 1 {
		build = cfg.Builds[0]
	} else if len(cfg.Builds) > 1 {
		return "", "", errors.New("multiple builds found, please specify one")
	}

	if len(build.Goos) > 0 {
		log.Printf("WARN: goos specified in your config file is overriden")
	}
	if len(build.Goarch) > 0 {
		log.Printf("WARN: goarch specified in your config file is overriden")
	}
	if len(build.Goarm) > 0 {
		log.Printf("WARN: goarm specified in your config file is overriden")
	}
	if len(build.Gomips) > 0 {
		log.Printf("WARN: gomips specified in your config file is overriden")
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
		Cmd: `cp "{{ .Path }}" /usr/local/bin/{{ .ProjectName }}`,
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
		return "", "", errors.New("multiple archives found, please specify one")
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

	if len(cfg.NFPMs) > 0 && target.Os != "linux" {
		log.Printf("WARN: nfpms section specified in your config file disabled")
		cfg.NFPMs = nil
	}

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return "", "", err
	}

	f, err := os.CreateTemp(os.TempDir(), ".goreleaser.yml")
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	if err := os.WriteFile(f.Name(), b, 0644); err != nil {
		defer os.Remove(f.Name())
		return "", "", err
	}

	if cli.Debug {
		log.Println("DBG: goreleaser config:")
		scanner := bufio.NewScanner(strings.NewReader(string(b)))
		for scanner.Scan() {
			log.Printf("  %s", scanner.Text())
		}
	}

	return f.Name(), cfg.Dist, nil
}
