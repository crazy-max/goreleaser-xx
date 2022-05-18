package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"sort"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/goreleaser/goreleaser/pkg/config"
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
	Debug        bool              `kong:"name='debug',env='DEBUG',default='false',help='Enable debug.'"`
	GitRef       string            `kong:"name='git-ref',env='GIT_REF',help='The branch or tag like refs/tags/v1.0.0 (default to your working tree info).'"`
	GoReleaser   string            `kong:"name='goreleaser',env='GORELEASER_PATH',default='goreleaser',help='Set a specific GoReleaser binary to use.'"`
	Config       string            `kong:"name='config',type='path',env='GORELEASER_CONFIG',help='Load GoReleaser configuration from file.'"`
	GoBinary     string            `kong:"name='go-binary',env='GORELEASER_GOBINARY',help='Set a specific go binary to use when building.'"`
	Name         string            `kong:"name='name',env='GORELEASER_NAME',help='Project name.'"`
	Dist         string            `kong:"name='dist',env='GORELEASER_DIST',default='./dist',help='Dist folder where artifact will be stored.'"`
	Artifacts    []string          `kong:"name='artifacts',env='GORELEASER_ARTIFACTS',enum='archive,bin',default='archive',help='Types of artifact to create.'"`
	Main         string            `kong:"name='main',env='GORELEASER_MAIN',default='.',help='Path to main.go file or main package.'"`
	Flags        string            `kong:"name='flags',env='GORELEASER_FLAGS',help='Custom flags templates.'"`
	Asmflags     string            `kong:"name='asmflags',env='GORELEASER_ASMFLAGS',help='Custom asmflags templates.'"`
	Gcflags      string            `kong:"name='gcflags',env='GORELEASER_GCFLAGS',help='Custom gcflags templates.'"`
	Ldflags      string            `kong:"name='ldflags',env='GORELEASER_LDFLAGS',help='Custom ldflags templates.'"`
	Tags         string            `kong:"name='tags',env='GORELEASER_TAGS',help='Custom build tags templates.'"`
	Files        []string          `kong:"name='files',env='GORELEASER_FILES',help='Additional files/template/globs you want to add to the archive.'"`
	Replacements map[string]string `kong:"name='replacements',env='GORELEASER_REPLACEMENTS',help='Replacements for GOOS and GOARCH in the archive/binary name.'"`
	Envs         []string          `kong:"name='envs',env='GORELEASER_ENVS',help='Custom environment variables to be set during the build.'"`
	PreHooks     []string          `kong:"name='pre-hooks',env='GORELEASER_PRE_HOOKS',help='Hooks which will be executed before the build.'"`
	PostHooks    []string          `kong:"name='post-hooks',env='GORELEASER_POST_HOOKS',help='Hooks which will be executed after the build.'"`
	Snapshot     bool              `kong:"name='snapshot',env='GORELEASER_SNAPSHOT',default='false',help='Run in snapshot mode.'"`
	Checksum     bool              `kong:"name='checksum',env='GORELEASER_CHECKSUM',default='true',help='Create checksum.'"`

	// Deprecated flags
	ArtifactType   string   `kong:"hidden,name='artifact-type',env='GORELEASER_ARTIFACTTYPE',help='artifacts'"`
	Hooks          []string `kong:"hidden,name='hooks',env='GORELEASER_HOOKS',help='pre-hooks'"`
	BuildPreHooks  []string `kong:"hidden,name='build-pre-hooks',env='GORELEASER_BUILD_PRE_HOOKS',help='pre-hooks'"`
	BuildPostHooks []string `kong:"hidden,name='build-post-hooks',env='GORELEASER_BUILD_POST_HOOKS',help='post-hooks'"`
}

func main() {
	var err error

	// Parse command line
	ctx := kong.Parse(&cli,
		kong.Name(name),
		kong.Description(fmt.Sprintf("%s. More info: %s", desc, url)),
		kong.UsageOnError(),
		kong.Vars{
			"version": version,
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
	signal.Notify(channel, os.Interrupt, SIGTERM)
	go func() {
		sig := <-channel
		log.Printf("WARN: caught signal %v", sig)
		os.Exit(0)
	}()

	if cli.Debug {
		log.Println("DBG: environment:")
		printEnv()
	}

	// Warn on deprecated flag usage and assign to new flag
	for _, f := range ctx.Flags() {
		if f.Hidden && f.Set {
			for _, fa := range ctx.Flags() {
				if fa.Name != f.Help {
					continue
				}
				if fa.Set {
					continue
				}
				switch f.Name {
				case "artifact-type":
					cli.Artifacts = []string{f.Target.String()}
				case "hooks", "build-pre-hooks":
					cli.PreHooks = f.Target.Interface().([]string)
				case "build-post-hooks":
					cli.PostHooks = f.Target.Interface().([]string)
				}
				break
			}
			log.Printf("WARN: --%s is deprecated and will be removed in a future release. Use --%s instead.", f.Name, f.Help)
		}
	}

	// Seek GoReleaser bin
	goreleaserPath, err := exec.LookPath(cli.GoReleaser)
	if err != nil {
		log.Fatalf("ERR: %s not found in PATH", cli.GoReleaser)
	}

	// Get target
	target := getTarget()
	if cli.Debug {
		log.Println("DBG: target:")
		log.Printf("  GOOS=%s", target.Os)
		log.Printf("  GOARCH=%s", target.Arch)
		log.Printf("  GOARM=%s", target.Arm)
		log.Printf("  GOMIPS=%s", target.Mips)
		log.Printf("  GOMIPS64=%s", target.Mips64)
		log.Printf("  GOAMD64=%s", target.Amd64)
	}

	// Get compilers
	compilers := getCompilers(target)

	// GoReleaser config
	cfg, err := getConfig(cli, target, compilers)
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
	defer func() {
		_ = os.Remove(cfg.Path)
		_ = os.RemoveAll(cfg.Project.Dist)
	}()

	if cli.Debug {
		log.Println("DBG: go env:")
		printGoEnv(cfg.Project, target)
	}

	// Set GoReleaser flags
	goreleaserFlags := []string{"release", "--config", cfg.Path}

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
		goreleaserFlags = append(goreleaserFlags, "--snapshot")
	}

	// Git status
	log.Println("INF: git status:")
	log.Printf("  tag=%s", gitTag)
	log.Printf("  dirty=%t", gitDirty)
	log.Printf("  wrongref=%t", gitWrongRef)

	// Start GoReleaser
	log.Printf("INF: %s %s", goreleaserPath, strings.Join(goreleaserFlags, " "))
	goreleaser := exec.Command(cli.GoReleaser, goreleaserFlags...)
	goreleaser.Stdout = os.Stdout
	goreleaser.Stderr = os.Stderr
	if err := goreleaser.Run(); err != nil {
		log.Fatalf("ERR: goreleaser failed: %v", err)
	}

	// Create dist folder
	if err := os.Mkdir(cli.Dist, 0755); err != nil {
		log.Fatal(err)
	}

	// Post build
	fdist, err := os.Open(cfg.Project.Dist)
	if err != nil {
		log.Printf("WARN: cannot open dist folder: %v", err)
	}
	defer fdist.Close()
	fis, err := fdist.Readdir(-1)
	if err != nil {
		log.Printf("WARN: cannot read dist folder: %v", err)
	}
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		archiveExt := ""
		for _, ext := range []string{".tar.gz", ".tar.xz", ".tar", ".gz", ".zip"} {
			if strings.HasSuffix(fi.Name(), ext) {
				archiveExt = ext
				break
			}
		}
		if len(archiveExt) > 0 {
			for _, atf := range cli.Artifacts {
				var atfPath string
				switch atf {
				case "bin":
					binName := cfg.Project.ProjectName
					atfName := strings.TrimSuffix(fi.Name(), archiveExt)
					if target.Os == "windows" {
						atfName += ".exe"
						binName += ".exe"
					}
					atfPath = path.Join(cli.Dist, atfName)
					if err := copyFile(path.Join("/usr/local/bin", binName), atfPath); err != nil {
						log.Fatalf("ERR: cannot copy binary: %v", err)
					}
				case "archive":
					atfPath = path.Join(cli.Dist, fi.Name())
					if err := copyFile(path.Join(fdist.Name(), fi.Name()), atfPath); err != nil {
						log.Fatalf("ERR: cannot copy archive: %v", err)
					}
				default:
					log.Fatalf("ERR: unknown artifact type: %s", atf)
				}
				log.Printf("INF: successfully created: %s", atfPath)
				if cli.Checksum {
					checksum(atfPath)
				}
			}
		} else {
			atfPath := path.Join(cli.Dist, fi.Name())
			if err := copyFile(path.Join(fdist.Name(), fi.Name()), atfPath); err != nil {
				log.Fatalf("ERR: cannot copy from dist: %v", err)
			}
			log.Printf("INF: successfully created: %s", atfPath)
		}
	}
}

func copyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
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
		log.Fatalf("ERR: failed to create checksum: %v", err)
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

	log.Printf("INF: checksum file created: %s", checksumFile)
}

func printEnv() {
	envs := os.Environ()
	sort.Strings(envs)
	for _, e := range envs {
		log.Printf("  %s", e)
	}
}

func printGoEnv(config config.Project, target Target) {
	bin := "go"
	if cli.GoBinary != "" {
		bin = cli.GoBinary
	}
	cmd := exec.Command(bin, "env")
	cmd.Env = append(os.Environ(), append([]string{
		fmt.Sprintf("GOOS=%s", target.Os),
		fmt.Sprintf("GOARCH=%s", target.Arch),
		fmt.Sprintf("GOARM=%s", target.Arm),
		fmt.Sprintf("GOMIPS=%s", target.Mips),
		fmt.Sprintf("GOMIPS64=%s", target.Mips64),
		fmt.Sprintf("GOAMD64=%s", target.Amd64),
	}, config.Env...)...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
	defer cmd.Wait()
	s := bufio.NewScanner(stdout)
	for s.Scan() {
		log.Printf("  %s", s.Text())
	}
}
