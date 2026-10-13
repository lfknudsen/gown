package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Microsoft/go-winio/pkg/guid"
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
	} else if matchesGallery(target) {
		err = DownloadGallery(target, outputDir)
	} else if rxVideo.MatchString(target) {
		err = DownloadVideo(target, outputDir)
	} else {
		// Fallback
		err = DownloadGallery(target, outputDir)
		if err != nil {
			err = DownloadVideo(target, outputDir)
		}
	}
	if err != nil {
		log.Println(err)
		// If nothing else worked, append the URL to a file.
		log.Printf("Unsupported URL:\n\t%s\n", target)
		filename := time.Now().UTC().Format(time.RFC3339)
		filename = strings.Replace(filename, ":", "_", -1)
		uid, err := guid.NewV4()
		if err == nil {
			filename += "-" + uid.String()
		}
		filename = fmt.Sprintf("queue/%s.txt", filename)
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.WriteString(target + "\n")
		if err != nil {
			log.Println(err)
		}
		_ = f.Close()
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
		} else if arg == "--generate-lists" {
			generateLists()
		} else {
			targets = append(targets, arg)
		}
	}

	if len(targets) > 0 {
		wait := sync.WaitGroup{}
		for _, target := range targets {
			wait.Go(func() { Download(target, autoRename, chStatus) })
		}
		wait.Wait()
	}
}
