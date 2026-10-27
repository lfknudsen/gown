package main

import (
	"os/exec"
	"strings"
)

func matchesMagnetURL(url string) bool {
	return strings.HasPrefix(url, "magnet:?")
}

const subDirTorrentFiles = "torrents"

func downloadMagnet(url string, outputDir string) error {
	cmd := exec.Command("node", "magnet-convert.js", url)
	cmd.Dir = outputDir
	err := cmd.Run()
	return err
}
