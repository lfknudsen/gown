package pkg

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var rxPostID = regexp.MustCompile("instagram.com/p/([^/]+)")

var InstagramDir = filepath.Join(filepath.Join(OutputDir, "instagram"))

func DownloadInstagram(url string) error {
	err := os.MkdirAll(InstagramDir, 0755)
	if err != nil {
		return err
	}

	match := rxPostID.FindStringSubmatch(url)
	if match == nil || len(match) != 2 {
		return errors.New("invalid Instagram URL")
	}
	cmd := exec.Command("instaloader", "-"+match[1])
	cmd.Dir = InstagramDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
