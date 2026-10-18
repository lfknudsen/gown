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

var rxUrlFilename *regexp.Regexp = regexp.MustCompile(
	`(?i)\s*(?:https?://)?(?:ww[^.]\.)?(?:[^/]/)+([^/]+])`)

func downloadFile(url string, outputDir string) error {
	matches := rxUrlFilename.FindStringSubmatch(url)
	if len(matches) != 2 {
		msg := fmt.Sprintf("Could not parse URL:\n%s\n", url)
		return errors.New(msg)
	}
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	filename := matches[1]
	out := filepath.Join(outputDir, filename)
	log.Printf("Filename: %s", out)

	f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		log.Printf("Error opening file %s for writing: %s", out, err)
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
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(f)

	resp2, err := http.Get(url)
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
	sz, err := io.Copy(wr, resp2.Body)
	if err != nil {
		return err
	}
	log.Printf("Downloaded %d bytes to\n\t%s\n", sz, out)
	return nil
}
