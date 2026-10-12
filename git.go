package main

import (
	"github.com/go-git/go-git/v6"
	"os"
	"log"
)

func DownloadGit(url string, outputDir string) bool {
	opts := git.CloneOptions {
		URL: url,
		Progress: os.Stdout
	}
	_, err := git.PlainClone(outputDir, &opts)
	if err != null {
		log.Println(err)
		return false
	}
	return true
}
