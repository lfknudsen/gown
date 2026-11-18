package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	. "gown/pkg/util"

	mapset "github.com/deckarep/golang-set/v3"
)

var supportedVideoURLsFile = filepath.Join(DataDir, "yt-dlp_sites.txt")

const subDirVideos = "yt-dlp"

func DownloadVideo(url string, outputDir string) error {
	outputDir = filepath.Join(outputDir, subDirVideos)
	cmd := exec.Command("yt-dlp",
		"--write-subs", "--write-auto-subs", "--sub-langs", "en.*", "--embed-subs",
		"--embed-thumbnail",
		"--embed-metadata", "--embed-chapters", "--write-description", "--write-info-json",
		"--windows-filenames",
		"-P", outputDir,
		url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var rxStem = regexp.MustCompile(`(?:https?//:)?(?:ww[^.]\.)?([^/ :]+)`)

// GenerateVideoList reads the table of supported websites and saves it as a
// list that can be quickly parsed by matchesVideo.
func GenerateVideoList() {
	cmd := exec.Command("yt-dlp", "--list-extractors")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	PANIC(err)

	set := mapset.NewSet[string]()
	scanner := bufio.NewScanner(&out)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		matches := rxStem.FindStringSubmatch(text)
		if len(matches) != 2 {
			fmt.Printf("Could not parse domain %s", text)
			continue
		}
		set.Add(strings.ToLower(matches[1]))
	}

	InitialiseDirectories(DataDir)
	f, err := os.OpenFile(supportedVideoURLsFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	defer CloseFile(f)
	wr := bufio.NewWriter(f)
	for entry := range set.Iter() {
		_, err = wr.WriteString(entry + "\n")
		PANIC(err)
	}
	_ = wr.Flush()
}

func MatchesVideo(url string) bool {
	f, err := os.OpenFile(supportedVideoURLsFile, os.O_RDONLY, 0666)
	if err != nil {
		GenerateVideoList()
		f, err = os.OpenFile(supportedVideoURLsFile, os.O_RDONLY, 0666)
		if err != nil {
			return false
		}
	}
	defer CloseFile(f)
	scanner := bufio.NewScanner(f)
	url = strings.ToLower(UrlStem(url))
	for scanner.Scan() {
		if scanner.Text() == url {
			return true
		}
	}
	return false
}
