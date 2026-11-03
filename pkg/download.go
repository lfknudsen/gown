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

// enums for URL download handlers.
const (
	FileDownloader = iota
	GalleryDownloader
	VideoDownloader
	MediafireDownloader
	MegaDownloader
	InstagramDownloader
)

// handlerStringToEnum translates a user-input string to the "enum" used
// internally.
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

// downloadUsingDownloader circumvents the traditional delegation function
// download and instead calls the relevant handler directly.
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

// download delegates the download process to specific programmes based on the URL.
// The status channel is currently unused; nil is acceptable.
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

// Download starts the download process, and adds the URL to the queue
// if the download could not complete for any reason.
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

// DownloadFromQueue starts the download process, but does not add the URL
// to the queue on failure (the intention being that it got the link from the
// queue in the first place).
func DownloadFromQueue(url string, chStatus chan int) error {
	return download(url, chStatus)
}

// InitialiseDirectories ensures that the given directory or directories exist by
// creating them if they do not.
func InitialiseDirectories(dir ...string) {
	for _, d := range dir {
		err := os.MkdirAll(d, os.ModeDir)
		FATAL(err)
	}
}

// SetupDirs creates the major subdirectories needed. Individual handlers
// may initialise and use their own if they so wish, but they can rely on
// e.g. OutputDir existing.
func SetupDirs() {
	InitialiseDirectories(OutputDir, QueueDir, DirMega, RepoDir)
}
