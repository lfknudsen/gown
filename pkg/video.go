package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	mapset "github.com/deckarep/golang-set/v3"
)

var supportedVideoURLsFile = filepath.Join(DataDir, "yt-dlp_sites.txt")

const subDirVideos = "yt-dlp"

func DownloadVideo(url string, outputDir string) error {
	outputDir = filepath.Join(outputDir, subDirVideos)
	cmd := exec.Command("yt-dlp",
		"--write-subs", "--write-auto-subs", "--sub-langs", "en.*", "--embed-subs", "--embed-thumbnail",
		"--embed-metadata", "--embed-chapters", "--write-description", "--write-info-json", "--windows-filenames", "-P", outputDir, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var rxStem = regexp.MustCompile(`(?:https?//:)?(?:ww[^.]\.)?([^/ :]+)`)

// GenerateVideoList reads the table of supported websites and saves it as a
// list that can be quickly parsed by matchesVideo.
func GenerateVideoList() {
	cmd := exec.Command("yt-dlp", "--list-extractors")
	r, w := io.Pipe()
	cmd.Stdout = w
	scanner := bufio.NewScanner(r)
	err := cmd.Run()
	PANIC(err)
	set := mapset.NewSet[string]()
	for scanner.Scan() {
		matches := rxStem.FindStringSubmatch(scanner.Text())
		if len(matches) != 2 {
			fmt.Printf("Could not parse domain %s", scanner.Text())
			continue
		}
		set.Add(matches[1])
	}
	_ = w.Close()
	_ = r.Close()

	InitialiseDirectories(DataDir)
	f, err := os.OpenFile(supportedVideoURLsFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	defer closeFile(f)
	wr := bufio.NewWriter(f)
	for entry := range set.Iter() {
		_, err = wr.WriteString(entry + "\n")
		PANIC(err)
	}
	_ = wr.Flush()
}

func matchesVideo(url string) bool {
	f, err := os.OpenFile(supportedVideoURLsFile, os.O_RDONLY, 0666)
	if err != nil {
		GenerateVideoList()
		f, err = os.OpenFile(supportedVideoURLsFile, os.O_RDONLY, 0666)
		if err != nil {
			return false
		}
	}
	defer closeFile(f)
	scanner := bufio.NewScanner(f)
	url = UrlStem(url)
	for scanner.Scan() {
		if scanner.Text() == url {
			return true
		}
	}
	return false
}
