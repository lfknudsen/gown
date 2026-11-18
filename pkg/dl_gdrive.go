package pkg

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
)

var googleDriveDirName = "Google Drive"
var googleDriveDir = filepath.Join(OutputDir, googleDriveDirName)

var rxGoogle = regexp.MustCompile(`(?:https?://)?(?:ww[.].)?(?:drive|docs)\.google\.com/(?:file/d/([a-zA-Z0-9_-]{10,})|(?:open|uc)\?(?:[^&]*&)*id=([a-zA-Z0-9_-]{10,}))/.+`)

func MatchesGoogleDrive(url string) bool {
	return rxGoogle.MatchString(url)
}

func DownloadGoogle(url string, outputDir string) error {
	matches := rxGoogle.FindStringSubmatch(url)
	if len(matches) < 2 {
		fmt.Printf("There were %d matches in Google Drive link!\n", len(matches))
		return errors.New("invalid GoogleDrive link")
	}
	InitialiseDirectories(googleDriveDir)
	id := matches[1]
	link := fmt.Sprintf("https://drive.google.com/uc?id=%s", id)
	cmd := exec.Command("gdown", link)
	cmd.Dir = googleDriveDir
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error downloading Google Drive: %v\n", err)
		return err
	}
	return nil
}
