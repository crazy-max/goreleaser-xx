package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alecthomas/kong"
)

var (
	cli     Cli
	version = "dev"
	name    = "goreleaser-xx"
	desc    = "Cross compilation helper for GoReleaser"
	url     = "https://github.com/crazy-max/goreleaser-xx"
)

// Cli holds command line args, flags and cmds
type Cli struct {
	Version    kong.VersionFlag
	Debug      bool     `kong:"name='debug',env='DEBUG',default='false',help='Enable debug.'"`
	GitRef     string   `kong:"name='git-ref',env='GIT_REF',help='The branch or tag like refs/tags/v1.0.0 (default to your working tree info).'"`
	GoReleaser string   `kong:"name='goreleaser',env='GORELEASER_PATH',default='/opt/goreleaser-xx/goreleaser',help='Path to GoReleaser binary.'"`
	Name       string   `kong:"name='name',env='GORELEASER_NAME',help='Project name.'"`
	Dist       string   `kong:"name='dist',env='GORELEASER_DIST',default='./dist',help='Dist folder where artifact will be stored.'"`
	Hooks      []string `kong:"name='hooks',env='GORELEASER_HOOKS',help='Hooks which will be executed before the build is started.'"`
	Main       string   `kong:"name='main',env='GORELEASER_MAIN',default='.',help='Path to main.go file or main package.'"`
	Ldflags    string   `kong:"name='ldflags',env='GORELEASER_LDFLAGS',help='Custom ldflags templates.'"`
	Files      []string `kong:"name='files',env='GORELEASER_FILES',help='Additional files/template/globs you want to add to the archive.'"`
	Snapshot   bool     `kong:"name='snapshot',env='GORELEASER_SNAPSHOT',default='false',help='Run in snapshot mode.'"`
}

func main() {
	var err error
	var grFlags []string

	// Parse command line
	_ = kong.Parse(&cli,
		kong.Name(name),
		kong.Description(fmt.Sprintf("%s. More info: %s", desc, url)),
		kong.UsageOnError(),
		kong.Vars{
			"version": fmt.Sprintf("%s/%s", name, version),
		},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	// Init
	log.SetFlags(0)
	log.Printf("INF: starting %s", fmt.Sprintf("%s/%s", name, version))

	// Handle os signals
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-channel
		log.Printf("WARN: caught signal %v", sig)
		os.Exit(0)
	}()

	if cli.Debug {
		log.Println("DBG: environment:")
		for _, e := range os.Environ() {
			log.Printf("  %s", e)
		}
	}

	// Compiler target
	target := getTarget()
	if cli.Debug {
		log.Printf("DBG: target: %+v", target)
	}

	// GoReleaser config
	grConfig, err := getGRConfig(cli, target)
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
	defer os.Remove(grConfig)
	grFlags = append(grFlags, "release", "--config", grConfig)

	// Git tag
	if strings.HasPrefix(cli.GitRef, "refs/tags/v") {
		if err := os.Setenv("GORELEASER_CURRENT_TAG", strings.TrimLeft(cli.GitRef, "refs/tags/v")); err != nil {
			log.Printf("WARN: cannot set GORELEASER_CURRENT_TAG env var: %v", err)
		}
	}
	gitTag, err := getGitTag()
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
	log.Printf("INF: git tag found: %s", gitTag)

	// Validate
	gitDirty := isGitDirty()
	log.Printf("INF: git dirty: %t", gitDirty)
	gitWrongRef := isWrongRef(gitTag)
	log.Printf("INF: git wrong ref: %t", gitWrongRef)
	if gitDirty || gitWrongRef || cli.Snapshot {
		grFlags = append(grFlags, "--snapshot")
	}

	// Display status
	log.Println("INF: git status:")
	log.Printf("  tag=%s", gitTag)
	log.Printf("  dirty=%t", gitDirty)
	log.Printf("  wrongref=%t", gitWrongRef)

	// Start GoReleaser
	log.Printf("INF: %s %s", cli.GoReleaser, strings.Join(grFlags, " "))
	if err = syscall.Exec(cli.GoReleaser, grFlags, os.Environ()); err != nil {
		log.Fatalf("ERR: goreleaser failed: %v", err)
	}
}
