package main

import (
	"os"
	"regexp"

	gogit "github.com/go-git/go-git/v6"
)

var rxGit = regexp.MustCompile(`(https?://)?((.+\.git)|(git(lab|hub)\.com/[^/]+/[^/]+(.git)?))`)

func DownloadGit(url string, outputDir string) error {
	opts := gogit.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}
	r, err := gogit.PlainClone(outputDir, &opts)
	if err != nil {
		return err
	}
	defer func() {
		err := r.Close()
		if err != nil {
			return
		}
	}()
	return nil
}
