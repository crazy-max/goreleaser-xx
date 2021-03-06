package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func git(args ...string) (string, error) {
	var extraArgs = []string{
		"-c", "log.showSignature=false",
	}

	args = append(extraArgs, args...)
	var cmd = exec.Command("git", args...)

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", errors.New(strings.TrimSuffix(stderr.String(), "\n"))
	}

	return strings.ReplaceAll(strings.Split(stdout.String(), "\n")[0], "'", ""), nil
}

func getGitTag() (string, error) {
	var tag string
	var err error
	for _, fn := range []func() (string, error){
		func() (string, error) {
			return os.Getenv("GORELEASER_CURRENT_TAG"), nil
		},
		func() (string, error) {
			return git("tag", "--points-at", "HEAD", "--sort", "-version:creatordate")
		},
		func() (string, error) {
			return git("describe", "--tags", "--abbrev=0")
		},
	} {
		tag, err = fn()
		if tag != "" || err != nil {
			return tag, err
		}
	}

	return tag, err
}

func isWrongRef(tag string) bool {
	if _, err := git("describe", "--exact-match", "--tags", "--match", tag); err != nil {
		return true
	}
	return false
}

func isGitDirty() bool {
	out, err := git("status", "--porcelain")
	return strings.TrimSpace(out) != "" || err != nil
}
