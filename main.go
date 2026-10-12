package main

import (
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

var outputDir string = "downloads"

var chStatus = make(chan int, 1)

func FmtUrl(url string) string {
	url = strings.TrimSpace(url)
	url = strings.TrimSuffix(url, "/")
	url = regexp.MustCompile("http://").ReplaceAllString(url, "https://")
	return url
}

func Download(target string, autoRename bool, chStatus chan int) {
	var err error = nil
	if strings.Contains(target, "mediafire.com") {
		err = downloadMediafire(target, outputDir, autoRename, chStatus)
	} else if rxGit.MatchString(target) {
		err = DownloadGit(target, outputDir)
	} else if rxGallery.MatchString(target) {
		err = DownloadGallery(target, outputDir)
	} else if rxVideo.MatchString(target) {
		err = DownloadVideo(target, outputDir)
	} else {
		// Fallback
		err = DownloadGallery(target, outputDir)
		if err != nil {
			err = DownloadVideo(target, outputDir)
			if err != nil {
				log.Printf("Unsupported URL:\n\t%s\n", target)
			}
		}
	}
	if err != nil {
		log.Panic(err)
	}
}

func Queue(url string) {

}

func main() {
	autoRename := true
	args := os.Args[1:]
	targets := make([]string, 0, 1)
	for _, arg := range args {
		if arg == "--overwrite" {
			autoRename = false
		} else {
			targets = append(targets, arg)
		}
	}

	wait := sync.WaitGroup{}
	for _, target := range targets {
		wait.Go(func() { Download(target, autoRename, chStatus) })
	}
	wait.Wait()
}
