package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// rxMediafire is used to extract the key information from a MediaFire URL so that we
// can bypass redirects and middle-steps.
var rxMediafire *regexp.Regexp = regexp.MustCompile(
	`(?i)\s*(?:https?://)?(?:www\.)?mediafire\.com/(?:view|file)/(([^/]+)/([^/]+))/?(?:file)?/?\s*`)

// rxLink is used to extract the direct download link from a Mediafire page's HTML.
var rxLink *regexp.Regexp = regexp.MustCompile(
	`aria-label="Download file"\s*href="([^"]+)"`)

// downloadFromQuickKey is a convenience function which internally calls downloadMediafire
func downloadFromQuickKey(quickKey string, outputDir string, chStatus chan int) error {
	return downloadMediafire("https://www.mediafire.com/?"+quickKey, outputDir, chStatus)
}

// downloadMediafire downloads the HTML for the relevant page, copies the direct link from the
// download button, and finally downloads the file from that direct link.
func downloadMediafire(url string, outputDir string, chStatus chan int) error {
	matches := rxMediafire.FindAllStringSubmatch(url, -1)
	if len(matches) == 0 || len(matches[0]) < 4 {
		msg := fmt.Sprintf("Could not parse URL:\n%s\n", url)
		return errors.New(msg)
	}

	url = "https://www.mediafire.com/file/" + matches[0][1]

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer closeReadCloser(resp.Body)
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("Failed to download from URL:\n%s\nReceived status code %s.\n", url, resp.Status)
		return errors.New(msg)
	}

	temp, err := os.CreateTemp(os.TempDir(), "*.html")
	_, err = io.Copy(temp, resp.Body)
	closeFile(temp)
	if err != nil {
		return err
	}

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
	DEBUG(fmt.Sprintf("Direct link: \"%s\"", directLink))

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	filename := matches[0][3]
	out := filepath.Join(outputDir, filename)
	DEBUG(fmt.Sprintf("Filename: %s", out))

	f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		// In this case, we may want change the name of the output file depending
		// on the -overwrite flag.
		DEBUG(fmt.Sprintf("Error opening file %s for writing: %s", out, err))
		if *overwrite {
			f, err = os.Create(out)
			if err != nil {
				return err
			}
		} else {
			i := 0
			withoutExt, _ := strings.CutSuffix(filename, filepath.Ext(filename))
			for err != nil {
				i++
				newFilename := fmt.Sprintf("%s (%d)%s", withoutExt, i, filepath.Ext(filename))
				out = filepath.Join(filepath.Dir(out), newFilename)
				f, err = os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
			}
		}
	}
	defer closeFile(f)

	resp2, err := http.Get(directLink)
	if err != nil {
		return err
	}
	defer closeReadCloser(resp2.Body)

	DEBUG("Filename: " + out)
	wr := bufio.NewWriter(f)
	sz, err := io.Copy(wr, resp2.Body)
	if err != nil {
		return err
	}
	INFO(fmt.Sprintf("Downloaded %d bytes to\n\t%s\n", sz, out))
	return nil
}
