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
	b, err := yaml.Marshal(&config.Project{
		ProjectName: cli.Name,
		Dist:        cli.Dist,
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
				Env: []string{
					"CGO_ENABLED=0",
				},
				Hooks: config.HookConfig{
					Post: []config.BuildHook{
						{
							Cmd: `cp "{{ .Path }}" /usr/local/bin/` + cli.Name,
						},
					},
				},
			},
		},
		Archives: []config.Archive{
			{
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
			},
		},
		Release: config.Release{
			Disable: true,
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
