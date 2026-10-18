package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/deckarep/golang-set/v3"
)

// supportedURLsFile is the path to the file which contains the list of URLs supported by gallery-dl
const supportedURLsFile = "gallery-dl.txt"

// DownloadGallery runs the `gallery-dl` command with the given url and output directory.
func DownloadGallery(url string, outputDir string) error {
	cmd := exec.Command("gallery-dl", "-d", outputDir, url)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// matchesGallery returns true if the given URL supported by `gallery-dl`.
func matchesGallery(url string) bool {
	f, err := os.OpenFile(supportedURLsFile, os.O_RDONLY, 0666)
	if err != nil {
		generateGalleryList()
		f, err = os.OpenFile(supportedURLsFile, os.O_RDONLY, 0666)
		if err != nil {
			return false
		}
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	scanner := bufio.NewScanner(f)
	url = UrlRoot(url)
	for scanner.Scan() {
		if scanner.Text() == url {
			return true
		}
	}
	return false
}

// generateGalleryList reads the table of supported websites and saves it as a
// list that can be quickly parsed by matchesGallery.
func generateGalleryList() {
	res, err := http.Get("https://codeberg.org/mikf/gallery-dl/raw/branch/master/docs/supportedsites.md")
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	reader := bufio.NewReader(res.Body)
	content, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(supportedURLsFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Panic(err)
		}
	}(f)

	set := mapset.NewSet[string]()

	var rxURL = regexp.MustCompile(`https://(?:www\.)?([^/"]+)`)
	found := rxURL.FindAllSubmatch(content, -1)
	for _, match := range found {
		set.Add(string(match[1]))
	}
	wr := bufio.NewWriter(f)
	for entry := range set.Iter() {
		_, err = wr.WriteString(entry + "\n")
		if err != nil {
			log.Panic(err)
		}
	}
	_ = wr.Flush()
}
