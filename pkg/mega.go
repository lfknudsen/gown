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

var DirMegaName = "mega"

var DirMega = filepath.Join(OutputDir, DirMegaName)

var rxFolderName = regexp.MustCompile(`<span class="not-loading selectable-txt">([^<]+)</span>`)

func downloadMega(url string) error {
	matches := rxMegaURL.FindStringSubmatch(url)[1:]
	if len(matches) < 1 || len(matches) > 2 {
		return errors.New(fmt.Sprintf("weirdly formatted MEGA URL: %s", url))
	}
	idSection := len(matches) - 1
	folders := strings.Split(matches[idSection], "/folder/")
	dirStructure := make([]string, 0, len(folders))
	for _, folder := range folders {
		folderName := folder
		idx := strings.IndexRune(folder, '#')
		if idx != -1 {
			folderName = folder[:idx]
		}
		dirStructure = append(dirStructure, folderName)
	}
	outputDir := filepath.Join(dirStructure...)
	outputDir = filepath.Join(DirMega, outputDir)
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
