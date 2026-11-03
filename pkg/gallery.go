package pkg

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	. "gown/pkg/util"

	"github.com/deckarep/golang-set/v3"
)

// supportedURLsFile is the path to the file which contains the list of URLs supported by gallery-dl
var supportedURLsFile = filepath.Join(DataDir, "gallery-dl.txt")

// DownloadGallery runs the `gallery-dl` command with the given url and output directory.
func DownloadGallery(url string, outputDir string) error {
	cmd := exec.Command("gallery-dl", "-R", RetriesStr, "-d", outputDir, url)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// matchesGallery returns true if the given URL supported by `gallery-dl`.
func matchesGallery(url string) bool {
	f, err := os.OpenFile(supportedURLsFile, os.O_RDONLY, 0666)
	if err != nil {
		GenerateGalleryList()
		f, err = os.OpenFile(supportedURLsFile, os.O_RDONLY, 0666)
		if err != nil {
			ERROR(err)
			return false
		}
	}
	defer CloseFile(f)
	scanner := bufio.NewScanner(f)
	url = UrlRoot(url)
	for scanner.Scan() {
		if scanner.Text() == url {
			return true
		}
	}
	return false
}

// GenerateGalleryList reads the table of supported websites and saves it as a
// list that can be quickly parsed by matchesGallery.
func GenerateGalleryList() {
	res, err := http.Get("https://codeberg.org/mikf/gallery-dl/raw/branch/master/docs/supportedsites.md")
	FATAL(err)
	defer CloseReadCloser(res.Body)
	reader := bufio.NewReader(res.Body)
	content, err := io.ReadAll(reader)
	FATAL(err)

	InitialiseDirectories(DataDir)
	f, err := os.OpenFile(supportedURLsFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	defer CloseFile(f)

	set := mapset.NewSet[string]()

	var rxURL = regexp.MustCompile(`https://(?:www\.)?([^/"]+)`)
	found := rxURL.FindAllSubmatch(content, -1)
	for _, match := range found {
		set.Add(string(match[1]))
	}
	wr := bufio.NewWriter(f)
	for entry := range set.Iter() {
		_, err = wr.WriteString(entry + "\n")
		PANIC(err)
	}
	_ = wr.Flush()
}
