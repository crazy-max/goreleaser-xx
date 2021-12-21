package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/goreleaser/goreleaser/pkg/config"
	"gopkg.in/yaml.v2"
)

func getGRConfig(cli Cli, target Target) (string, string, error) {
	dist, err := os.MkdirTemp(os.TempDir(), "dist")
	if err != nil {
		return "", "", err
	}

	var arFiles []config.File
	for _, f := range cli.Files {
		arFiles = append(arFiles, config.File{
			Source: f,
		})
	}

	var buildPreHooks []config.Hook
	for _, cmd := range cli.BuildPreHooks {
		buildPreHooks = append(buildPreHooks, config.Hook{
			Cmd: cmd,
		})
	}

	var buildPostHooks = []config.Hook{
		{
			Cmd: `cp "{{ .Path }}" /usr/local/bin/{{ .ProjectName }}`,
		},
	}
	for _, cmd := range cli.BuildPostHooks {
		buildPostHooks = append(buildPostHooks, config.Hook{
			Cmd: cmd,
		})
	}

	var flags config.FlagArray
	if len(cli.Flags) > 0 {
		flags = append(flags, cli.Flags)
	}

	var asmflags config.StringArray
	if len(cli.Asmflags) > 0 {
		asmflags = append(asmflags, cli.Asmflags)
	}

	var gcflags config.StringArray
	if len(cli.Gcflags) > 0 {
		gcflags = append(gcflags, cli.Gcflags)
	}

	var ldflags config.StringArray
	if len(cli.Ldflags) > 0 {
		ldflags = append(ldflags, cli.Ldflags)
	}

	var tags config.FlagArray
	if len(cli.Tags) > 0 {
		tags = append(tags, cli.Tags)
	}

	b, err := yaml.Marshal(&config.Project{
		ProjectName: cli.Name,
		Dist:        dist,
		Before: config.Before{
			Hooks: cli.Hooks,
		},
		Builds: []config.Build{
			{
				Main:     cli.Main,
				Flags:    flags,
				Asmflags: asmflags,
				Gcflags:  gcflags,
				Ldflags:  ldflags,
				Tags:     tags,
				Goos:     []string{target.Os},
				Goarch:   []string{target.Arch},
				Goarm:    []string{target.Arm},
				Gomips:   []string{target.Mips},
				Env:      append([]string{"CGO_ENABLED=0"}, cli.Envs...),
				GoBinary: cli.GoBinary,
				Hooks: config.BuildHookConfig{
					Pre:  buildPreHooks,
					Post: buildPostHooks,
				},
			},
		},
		Archives: []config.Archive{
			{
				Replacements: cli.Replacements,
				FormatOverrides: []config.FormatOverride{
					{
						Goos:   "windows",
						Format: "zip",
					},
				},
				Files: arFiles,
			},
		},
		Checksum: config.Checksum{
			Disable: true,
		},
		Release: config.Release{
			Disable: true,
		},
		Changelog: config.Changelog{
			Skip: true,
		},
	})
	if err != nil {
		return "", "", err
	}

	file, err := ioutil.TempFile(os.TempDir(), ".goreleaser.yml")
	if err != nil {
		return "", "", err
	}
	if err := ioutil.WriteFile(file.Name(), b, 0644); err != nil {
		_ = os.Remove(file.Name())
		return "", "", err
	}

	if cli.Debug {
		log.Println("DBG: goreleaser config:")
		scanner := bufio.NewScanner(strings.NewReader(string(b)))
		for scanner.Scan() {
			log.Printf("  %s", scanner.Text())
		}
	}

	return file.Name(), dist, nil
}
