package pkg

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const VERSION = "0.3.1"

func getExecutableDir() string {
	ex, err := os.Executable()
	FATAL(err)
	return filepath.Dir(ex)
}

var ExecutableDir = getExecutableDir()
var GownDirectory = filepath.Join(ExecutableDir, "gown")
var OutputDir = filepath.Join(GownDirectory, "downloads")
var DataDir = filepath.Join(GownDirectory, "data")

var ChStatus = make(chan int, 1)

func FmtUrl(url string) string {
	url = strings.TrimSpace(url)
	url = strings.TrimSuffix(url, "/")
	url = regexp.MustCompile("http://").ReplaceAllString(url, "https://")
	return url
}

func download(url string, chStatus chan int) error {
	WARN("Downloading " + url)
	var err error = nil
	if strings.Contains(url, "mediafire.com") {
		err = downloadMediafire(url, OutputDir, chStatus)
	} else if rxGit.MatchString(url) {
		err = DownloadGit(url)
	} else if matchesGallery(url) {
		err = DownloadGallery(url, OutputDir)
	} else if matchesVideo(url) {
		err = DownloadVideo(url, OutputDir)
	} else if matchesMegaURL(url) {
		err = downloadMega(url)
	} else {
		// Fallback
		err = downloadFile(url, OutputDir)
	}
	return err
}

func Download(url string, chStatus chan int) error {
	err := download(url, chStatus)
	if err != nil {
		ERROR("Unsupported URL:\n\t" + url + "\n")
		// If nothing else worked, add the URL to a queue file.
		Queue(url)
	}
	return err
}

func DownloadFromQueue(url string, chStatus chan int) error {
	return download(url, chStatus)
}

// InitialiseDirectories ensures that the given directories exist.
func InitialiseDirectories(dir ...string) {
	for _, d := range dir {
		err := os.MkdirAll(d, os.ModeDir)
		FATAL(err)
	}
}

func SetupDirs() {
	InitialiseDirectories(OutputDir, QueueDir, DirMega, RepoDir)
}
