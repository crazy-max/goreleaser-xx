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

func getGRConfig(cli Cli, target Target) (string, error) {
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

	var ldflags config.StringArray
	if len(cli.Ldflags) > 0 {
		ldflags = append(ldflags, cli.Ldflags)
	}

	b, err := yaml.Marshal(&config.Project{
		ProjectName: cli.Name,
		Dist:        cli.Dist,
		Before: config.Before{
			Hooks: cli.Hooks,
		},
		Builds: []config.Build{
			{
				Main:    cli.Main,
				Flags:   flags,
				Ldflags: ldflags,
				Goos:    []string{target.Os},
				Goarch:  []string{target.Arch},
				Goarm:   []string{target.Arm},
				Gomips:  []string{target.Mips},
				Env:     append([]string{"CGO_ENABLED=0"}, cli.Envs...),
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
		return "", err
	}

	file, err := ioutil.TempFile(os.TempDir(), ".goreleaser.yml")
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(file.Name(), b, 0644); err != nil {
		os.Remove(file.Name())
		return "", err
	}

	if cli.Debug {
		log.Println("DBG: goreleaser config:")
		scanner := bufio.NewScanner(strings.NewReader(string(b)))
		for scanner.Scan() {
			log.Printf("  %s", scanner.Text())
		}
	}

	return file.Name(), nil
}
