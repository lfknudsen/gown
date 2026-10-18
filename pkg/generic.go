package pkg

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
	DEBUG("Filename: " + out)

	f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		DEBUG(fmt.Sprintf("Error opening file %s for writing: %s", out, err))
		if *Overwrite {
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

	resp2, err := http.Get(url)
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
