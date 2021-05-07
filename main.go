package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
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
	Version      kong.VersionFlag
	Debug        bool     `kong:"name='debug',env='DEBUG',default='false',help='Enable debug.'"`
	GitRef       string   `kong:"name='git-ref',env='GIT_REF',help='The branch or tag like refs/tags/v1.0.0 (default to your working tree info).'"`
	GoReleaser   string   `kong:"name='goreleaser',env='GORELEASER_PATH',default='/opt/goreleaser-xx/goreleaser',help='Path to GoReleaser binary.'"`
	Name         string   `kong:"name='name',env='GORELEASER_NAME',help='Project name.'"`
	Dist         string   `kong:"name='dist',env='GORELEASER_DIST',default='./dist',help='Dist folder where artifact will be stored.'"`
	ArtifactType string   `kong:"name='artifact-type',env='GORELEASER_ARTIFACTTYPE',enum='all,archive,bin',default='archive',help='Which type of artifact to create. Can be all, archive or bin.'"`
	Hooks        []string `kong:"name='hooks',env='GORELEASER_HOOKS',help='Hooks which will be executed before the build is started.'"`
	Main         string   `kong:"name='main',env='GORELEASER_MAIN',default='.',help='Path to main.go file or main package.'"`
	Ldflags      string   `kong:"name='ldflags',env='GORELEASER_LDFLAGS',help='Custom ldflags templates.'"`
	Files        []string `kong:"name='files',env='GORELEASER_FILES',help='Additional files/template/globs you want to add to the archive.'"`
	Envs         []string `kong:"name='envs',env='GORELEASER_ENVS',help='Custom environment variables to be set during the build.'"`
	Snapshot     bool     `kong:"name='snapshot',env='GORELEASER_SNAPSHOT',default='false',help='Run in snapshot mode.'"`
	Checksum     bool     `kong:"name='checksum',env='GORELEASER_CHECKSUM',default='true',help='Create checksum.'"`
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
		if err := os.Setenv("GORELEASER_CURRENT_TAG", strings.TrimLeft(cli.GitRef, "refs/tags/")); err != nil {
			log.Printf("WARN: cannot set GORELEASER_CURRENT_TAG env var: %v", err)
		}
	}
	gitTag, err := getGitTag()
	if err != nil {
		gitTag = "v0.0.0"
	}

	// Git validate
	gitDirty := isGitDirty()
	gitWrongRef := isWrongRef(gitTag)
	if gitDirty || gitWrongRef || cli.Snapshot {
		grFlags = append(grFlags, "--snapshot")
	}

	// Git status
	log.Println("INF: git status:")
	log.Printf("  tag=%s", gitTag)
	log.Printf("  dirty=%t", gitDirty)
	log.Printf("  wrongref=%t", gitWrongRef)

	// Start GoReleaser
	log.Printf("INF: %s %s", cli.GoReleaser, strings.Join(grFlags, " "))
	goreleaser := exec.Command(cli.GoReleaser, grFlags...)
	goreleaser.Stdout = os.Stdout
	goreleaser.Stderr = os.Stderr
	if err := goreleaser.Run(); err != nil {
		log.Fatalf("ERR: goreleaser failed: %v", err)
	}

	// Post build
	distFolder, err := os.Open(cli.Dist)
	if err != nil {
		log.Printf("WARN: cannot open dist foler: %v", err)
	}
	defer distFolder.Close()
	names, err := distFolder.Readdir(-1)
	if err != nil {
		log.Printf("WARN: cannot read dist foler: %v", err)
	}
	for _, name := range names {
		if name.IsDir() {
			if err := os.RemoveAll(path.Join(cli.Dist, name.Name())); err != nil {
				log.Printf("WARN: cannot remove: %v", err)
			}
			continue
		}
		if strings.HasPrefix(name.Name(), "config") {
			if err := os.Remove(path.Join(cli.Dist, name.Name())); err != nil {
				log.Printf("WARN: cannot remove: %v", err)
			}
			continue
		}
		if cli.Checksum {
			checksum(path.Join(cli.Dist, name.Name()))
		}
	}
}

func checksum(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("ERR: failed to open file to compute: %v", err)
	}
	defer file.Close()
	checksumFile := filename + ".sha256"

	h := sha256.New()
	_, err = io.Copy(h, file)
	if err != nil {
		log.Fatalf("ERR: failed to checksum: %v", err)
	}

	sha256file, err := os.OpenFile(checksumFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		log.Fatalf("ERR: failed to open file: %v", err)
	}
	defer sha256file.Close()

	_, err = sha256file.WriteString(hex.EncodeToString(h.Sum(nil)))
	if err != nil {
		log.Fatalf("ERR: failed to write file: %v", err)
	}

	log.Printf("INF: checksum file created in %s", checksumFile)
}
