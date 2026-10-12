package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	gogit "github.com/go-git/go-git/v6"
)

var rxGit = regexp.MustCompile(`(?:https?://)?(?:git(?:lab|hub)\.com/[^/]+/([^/]+)(?:.git)?|.+/([^/]+)\.git)`)

func DownloadPublicGit(url string, outputDir string) error {
	opts := gogit.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}
	r, err := gogit.PlainClone(outputDir, &opts)
	if err != nil {
		return err
	}
	return r.Close()
}

func DownloadGit(url string, outputDir string) error {
	subs := rxGit.FindStringSubmatch(url)
	if len(subs) <= 1 {
		return fmt.Errorf("invalid git url: %s", url)
	}
	name := UrlStem(subs[1])
	dir := filepath.Join(outputDir, name)
	cmd := exec.Command("git", "clone", url, dir)
	return cmd.Run()
}
