package pkg

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
)

var rxPostID = regexp.MustCompile("instagram.com/p/([^/]+)")

func DownloadInstagram(url string) error {
	match := rxPostID.FindStringSubmatch(url)
	if match == nil || len(match) != 2 {
		return errors.New("invalid Instagram URL")
	}
	cmd := exec.Command("instaloader", "-"+match[1])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
