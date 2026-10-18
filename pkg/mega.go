package pkg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var rxMegaURL = regexp.MustCompile(
	`(?:https?://)?(?:www\.)?mega(?:\.co)?\.nz/(folder|file)/(.+)`)

func matchesMegaURL(url string) bool {
	return rxMegaURL.MatchString(url)
}

var dirMega = "mega"

var rxFolderName = regexp.MustCompile(`<span class="not-loading selectable-txt">([^<]+)</span>`)

func downloadMega(url string, outputDir string) error {
	matches := rxMegaURL.FindStringSubmatch(url)[1:]
	if len(matches) < 1 || len(matches) > 2 {
		return errors.New(fmt.Sprintf("weirdly formatted MEGA URL: %s", url))
	}
	idSection := len(matches) - 1
	idx := strings.IndexRune(matches[idSection], '#')
	folderName := matches[idSection][:idx]
	outputDir = filepath.Join(outputDir, dirMega, folderName)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command("megatools", "dl", "--no-progress", "--path", outputDir, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
