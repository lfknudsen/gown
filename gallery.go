package main

import (
	"os/exec"
	"regexp"
)

var rxGallery = regexp.MustCompile(`(?:https?://)?(?:www.)?deviantart.com/.+`)

func DownloadGallery(url string, outputDir string) error {
	cmd := exec.Command("gallery-dl", "-d", outputDir, url)
	return cmd.Run()
}
