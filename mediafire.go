package main

import (
	"bufio"
	"errors"
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

func downloadFromQuickKey(quickKey string, outputDir string, autoRename bool, chStatus chan int) error {
	return downloadMediafire("https://www.mediafire.com/?"+quickKey, outputDir, autoRename, chStatus)
}

func downloadMediafire(url string, outputDir string, autoRename bool, chStatus chan int) error {
	matches := rx.FindAllStringSubmatch(url, -1)
	if len(matches) == 0 || len(matches[0]) < 4 {
		msg := fmt.Sprintf("Could not parse URL:\n%s\n", url)
		return errors.New(msg)
	}

	url = "https://www.mediafire.com/file/" + matches[0][1]

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("Failed to download from URL:\n%s\nReceived status code %s.\n", url, resp.Status)
		return errors.New(msg)
	}

	temp, err := os.CreateTemp(os.TempDir(), "*.html")
	defer func(temp *os.File) {
		_ = temp.Close()
	}(temp)

	sz, err := io.Copy(temp, resp.Body)
	if err != nil {
		return err
	}
	_ = temp.Close()

	tempFile := temp.Name()
	html, err := os.ReadFile(tempFile)
	if err != nil {
		return err
	}
	htmlStr := string(html)
	match := rxLink.FindStringSubmatch(htmlStr)
	if match == nil || len(match) < 2 {
		msg := fmt.Sprintf("Could not find regexp in HTML file %s", tempFile)
		return errors.New(msg)
	}
	directLink := match[1]
	log.Printf("Direct link: \"%s\"", directLink)

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
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
				return err
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
		return err
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
		return err
	}
	log.Printf("Downloaded %d bytes to\n\t%s\n", sz, out)
	return nil
}
