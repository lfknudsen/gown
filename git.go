package main

import (
	"log"
	"os"

	"github.com/go-git/go-git/v6"
)

func DownloadGit(url string, outputDir string) bool {
	opts := git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}
	r, err := git.PlainClone(outputDir, &opts)
	if err != nil {
		log.Println(err)
		return false
	}
	defer func() {
		err := r.Close()
		if err != nil {
			return
		}
	}()
	return true
}
