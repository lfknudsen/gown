package main

import (
	"os/exec"
	"regexp"
)

var rxMegaURL = regexp.MustCompile(`(?:https?://)?(?:www\.)?mega\.nz/.+`)

func matchesMegaURL(url string) bool {
	return rxMegaURL.MatchString(url)
}

func downloadMega(url string, outputDir string) error {
	cmd := exec.Command("megatools", "--no-progress", "--path", outputDir, url)
	return cmd.Run()
}
