package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var rx *regexp.Regexp = regexp.MustCompile(
	`(?i)\s*(?:(?:http|https)://)?(?:www.)?mediafire.com/(?:view|file)/(([^/]+)/([^/]+))/?(?:file)?/?\s*`)

var rxLink *regexp.Regexp = regexp.MustCompile(
	`aria-label="Download file"\s*href="([^"]+)"`)

func downloadFromQuickKey(quickKey string, outputDir string, autoRename bool, chStatus chan int) {
	downloadMediafire("https://www.mediafire.com/?"+quickKey, outputDir, autoRename, chStatus)
}

func downloadMediafire(url string, outputDir string, autoRename bool, chStatus chan int) bool {
	matches := rx.FindAllStringSubmatch(url, -1)
	if len(matches) == 0 || len(matches[0]) < 4 {
		log.Printf("Could not parse URL:\n%s\n", url)
		return false
	}

	url = "https://www.mediafire.com/file/" + matches[0][1]

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		return false
	}

	temp, err := os.CreateTemp(os.TempDir(), "*.html")
	defer func(temp *os.File) {
		_ = temp.Close()
	}(temp)

	sz, err := io.Copy(temp, resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}
	_ = temp.Close()

	tempFile := temp.Name()
	html, err := os.ReadFile(tempFile)
	if err != nil {
		log.Println(err)
		return false
	}
	htmlStr := string(html)
	match := rxLink.FindStringSubmatch(htmlStr)
	if match == nil || len(match) < 2 {
		log.Printf("Could not find regexp in HTML file %s", tempFile)
		return false
	}
	directLink := match[1]
	log.Printf("Direct link: \"%s\"", directLink)

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Println("Error creating output directory:\n", err)
		return false
	}

	filename := matches[0][3]
	out := filepath.Join(outputDir, filename)
	log.Printf("Filename: %s", out)

	f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		log.Printf("Error opening file %s for writing: %s", out, err)
		if autoRename {
			i := 0
			withoutExt, _ := strings.CutSuffix(filename, filepath.Ext(filename))
			for err != nil {
				i++
				newFilename := fmt.Sprintf("%s (%d)%s", withoutExt, i, filepath.Ext(filename))
				out = filepath.Join(filepath.Dir(out), newFilename)
				f, err = os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
			}
		} else {
			f, err = os.Create(out)
			if err != nil {
				log.Println(err)
				return false
			}
		}
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(f)

	resp2, err := http.Get(directLink)
	if err != nil {
		log.Println(err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp2.Body)

	log.Printf("Filename: %s", out)
	wr := bufio.NewWriter(f)
	sz, err = io.Copy(wr, resp2.Body)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Printf("Downloaded %d bytes to\n\t%s\n", sz, out)
}
