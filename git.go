package main

import (
	"os"

	gogit "github.com/go-git/go-git/v6"
)

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
