package main

import (
	"os/exec"
	"regexp"
)

var rxVideo = regexp.MustCompile(`(?:https?://)?(?:www.)?youtube.com/.+`)

func DownloadVideo(url string, outputDir string) error {
	cmd := exec.Command("yt-dlp",
		"--write-subs", "--windows-filenames", "-P", outputDir, url)
	return cmd.Run()
}
