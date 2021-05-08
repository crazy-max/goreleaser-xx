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
	var archives []config.Archive
	buildHooks := []config.BuildHook{
		{
			Cmd: `cp "{{ .Path }}" /usr/local/bin/{{ .Name }}`,
		},
	}

	if cli.ArtifactType == "all" || cli.ArtifactType == "bin" {
		archives = append(archives, config.Archive{
			ID:     "binary",
			Format: "binary",
		})
		buildHooks = append(buildHooks, config.BuildHook{
			Cmd: `cp "{{ .Path }}" ` + cli.Dist + `/{{ .ProjectName }}_{{ .Version }}_{{ .Target }}{{ .Ext }}`,
		})
	}
	if cli.ArtifactType == "all" || cli.ArtifactType == "archive" {
		archives = append(archives, config.Archive{
			ID: "archive",
			Replacements: map[string]string{
				"386":   "i386",
				"amd64": "x86_64",
			},
			FormatOverrides: []config.FormatOverride{
				{
					Goos:   "windows",
					Format: "zip",
				},
			},
			Files: cli.Files,
		})
	}

	b, err := yaml.Marshal(&config.Project{
		ProjectName: cli.Name,
		Dist:        cli.Dist,
		Before: config.Before{
			Hooks: cli.Hooks,
		},
		Builds: []config.Build{
			{
				Main: cli.Main,
				Ldflags: []string{
					cli.Ldflags,
				},
				Goos:   []string{target.Os},
				Goarch: []string{target.Arch},
				Goarm:  []string{target.Arm},
				Gomips: []string{target.Mips},
				Env:    append([]string{"CGO_ENABLED=0"}, cli.Envs...),
				Hooks: config.HookConfig{
					Post: buildHooks,
				},
			},
		},
		Archives: archives,
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
