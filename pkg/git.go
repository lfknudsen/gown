package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	gogit "github.com/go-git/go-git/v6"
)

const RepoDirName = "repositories"

var RepoDir = filepath.Join(OutputDir, RepoDirName)

var rxGit = regexp.MustCompile(`(?:https?://)?(?:gitlab\.com|github\.com|codeberg\.org)/[^/]+/([^/]+)(?:\.git)?|.+/([^/]+)\.git`)

func DownloadPublicGit(url string) error {
	opts := gogit.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}
	r, err := gogit.PlainClone(RepoDir, &opts)
	if err != nil {
		return err
	}
	return r.Close()
}

func DownloadGit(url string) error {
	subs := rxGit.FindStringSubmatch(url)
	if len(subs) <= 1 {
		return fmt.Errorf("invalid git url: %s", url)
	}
	name := UrlStem(subs[1])
	dir := filepath.Join(RepoDir, name)
	cmd := exec.Command("git", "clone", url, dir)
	return cmd.Run()
}
