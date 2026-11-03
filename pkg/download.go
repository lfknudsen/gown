package pkg

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const VERSION = "0.4.0"

func getExecutableDir() string {
	ex, err := os.Executable()
	FATAL(err)
	return filepath.Dir(ex)
}

func determineGownDir() string {
	gownDir := os.Getenv("GOWN_DIR")
	if gownDir == "" {
		return filepath.Join(ExecutableDir, "gown")
	}
	return gownDir
}

var ExecutableDir = getExecutableDir()
var GownDirectory = determineGownDir()
var OutputDir = filepath.Join(GownDirectory, "downloads")
var DataDir = filepath.Join(GownDirectory, "data")

var ChStatus = make(chan int, 1)

func FmtUrl(url string) string {
	url = strings.TrimSpace(url)
	url = strings.TrimSuffix(url, "/")
	url = regexp.MustCompile("http://").ReplaceAllString(url, "https://")
	return url
}

const (
	FileDownloader = iota
	GalleryDownloader
	VideoDownloader
	MediafireDownloader
	MegaDownloader
	InstagramDownloader
)

func handlerStringToEnum(handler string) int {
	handler = strings.ToLower(handler)
	switch handler {
	case "gallery":
		return GalleryDownloader
	case "video":
		return VideoDownloader
	case "mediafire":
		return MediafireDownloader
	case "mega":
		return MegaDownloader
	case "instagram":
		return InstagramDownloader
	default:
		return FileDownloader
	}
}

func downloadUsingDownloader(handler int, url string, chStatus chan int) error {
	var err error
	switch handler {
	case GalleryDownloader:
		err = DownloadGallery(url, OutputDir)
	case VideoDownloader:
		err = DownloadVideo(url, OutputDir)
	case MediafireDownloader:
		err = downloadMediafire(url, OutputDir, chStatus)
	case MegaDownloader:
		err = downloadMega(url)
	case InstagramDownloader:
		err = DownloadInstagram(url)
	default:
		err = downloadFile(url, OutputDir)
	}
	if err == nil {
		WARN("Finished downloading " + url)
	}
	return err
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
		// Some domains are matched by both because they can contain either images or videos.
		if err != nil && matchesVideo(url) {
			err = DownloadVideo(url, OutputDir)
		}
	} else if matchesVideo(url) {
		err = DownloadVideo(url, OutputDir)
	} else if matchesMegaURL(url) {
		err = downloadMega(url)
	} else if strings.Contains(url, "instagram.com") {
		err = DownloadInstagram(url)
	} else {
		// Fallback
		err = downloadFile(url, OutputDir)
	}
	if err == nil {
		WARN("Finished downloading " + url)
	}
	return err
}

func Download(url string, chStatus chan int) error {
	var err error = nil
	if *Using != "" {
		handler := handlerStringToEnum(*Using)
		err = downloadUsingDownloader(handler, url, chStatus)
	} else {
		err = download(url, chStatus)
	}

	if err != nil {
		ERROR("Unsupported URL:\n\t" + url + "\n")
		// If nothing else worked, add the URL to a queue file.
		err = Queue(url)
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
