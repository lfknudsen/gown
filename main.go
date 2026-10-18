package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func getExecutableDir() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(ex)
}

var root = getExecutableDir()
var outputDir string = filepath.Join(root, "downloads")
var queueDir string = filepath.Join(root, "queue")

var chStatus = make(chan int, 1)

func FmtUrl(url string) string {
	url = strings.TrimSpace(url)
	url = strings.TrimSuffix(url, "/")
	url = regexp.MustCompile("http://").ReplaceAllString(url, "https://")
	return url
}

func Download(target string, overwrite bool, chStatus chan int) {
	var err error = nil
	if strings.Contains(target, "mediafire.com") {
		err = downloadMediafire(target, outputDir, overwrite, chStatus)
	} else if rxGit.MatchString(target) {
		err = DownloadGit(target, outputDir)
	} else if matchesGallery(target) {
		err = DownloadGallery(target, outputDir)
	} else if matchesVideo(target) {
		err = DownloadVideo(target, outputDir)
	} else if matchesMegaURL(target) {
		err = downloadMega(target, outputDir)
	} else {
		// Fallback
		err = downloadFile(target, outputDir, overwrite)
	}
	if err != nil {
		log.Println(err)
		// If nothing else worked, add the URL to a queue file.
		log.Printf("Unsupported URL:\n\t%s\n", target)
		Queue(target)
	}
}

// initialiseDirectories ensures that the given directories exist.
func initialiseDirectories(dir ...string) {
	for _, d := range dir {
		err := os.MkdirAll(d, os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	overwrite := false
	args := os.Args[1:]
	targets := make([]string, 0, 1)
	for _, arg := range args {
		if arg == "--overwrite" {
			overwrite = true
		} else if arg == "--generate-lists" {
			generateGalleryList()
		} else {
			targets = append(targets, arg)
		}
	}

	initialiseDirectories(outputDir, queueDir)

	if len(targets) > 0 {
		wait := sync.WaitGroup{}
		for _, target := range targets {
			wait.Go(func() { Download(target, overwrite, chStatus) })
		}
		wait.Wait()
	}
}
