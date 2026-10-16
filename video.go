package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	mapset "github.com/deckarep/golang-set/v3"
)

const supportedVideoURLsFile = "data/yt-dlp_sites.txt"

const subDir = "yt-dlp"

func DownloadVideo(url string, outputDir string) error {
	outputDir = filepath.Join(outputDir, subDir)
	cmd := exec.Command("yt-dlp",
		"--write-subs", "--windows-filenames", "-P", outputDir, url)
	return cmd.Run()
}

var rxStem = regexp.MustCompile(`(?:https?//:)?(?:ww[^.]\.)?([^/ :]+)`)

// generateVideoLists reads the table of supported websites and saves it as a
// list that can be quickly parsed by matchesVideo.
func generateVideoLists() {
	cmd := exec.Command("yt-dlp", "--list-extractors")
	r, w := io.Pipe()
	cmd.Stdout = w
	scanner := bufio.NewScanner(r)

	set := mapset.NewSet[string]()
	for scanner.Scan() {
		matches := rxStem.FindStringSubmatch(scanner.Text())
		if len(matches) != 2 {
			fmt.Printf("Could not parse domain %s", scanner.Text())
			continue
		}
		set.Add(matches[1])
	}
	_ = w.Close()
	_ = r.Close()

	f, err := os.OpenFile(supportedVideoURLsFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	wr := bufio.NewWriter(f)
	for entry := range set.Iter() {
		_, err = wr.WriteString(entry + "\n")
		if err != nil {
			log.Panic(err)
		}
	}
	_ = wr.Flush()
}

func matchesVideo(url string) bool {
	f, err := os.OpenFile(supportedVideoURLsFile, os.O_RDONLY, 0666)
	if err != nil {
		generateVideoLists()
		f, err = os.OpenFile(supportedVideoURLsFile, os.O_RDONLY, 0666)
		if err != nil {
			return false
		}
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	scanner := bufio.NewScanner(f)
	url = UrlStem(url)
	for scanner.Scan() {
		if scanner.Text() == url {
			return true
		}
	}
	return false
}
