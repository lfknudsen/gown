package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const VERSION = "0.3.1"

func getExecutableDir() string {
	ex, err := os.Executable()
	if err != nil {
		FATAL("Error while retrieving the path to the Gown executable:", err)
	}
	return filepath.Dir(ex)
}

var executableDir = getExecutableDir()
var gownDirectory = filepath.Join(executableDir, "gown")
var outputDir = filepath.Join(gownDirectory, "downloads")

var chStatus = make(chan int, 1)

func FmtUrl(url string) string {
	url = strings.TrimSpace(url)
	url = strings.TrimSuffix(url, "/")
	url = regexp.MustCompile("http://").ReplaceAllString(url, "https://")
	return url
}

func Download(url string, chStatus chan int) {
	DEBUG("Downloading", url)
	var err error = nil
	if strings.Contains(url, "mediafire.com") {
		err = downloadMediafire(url, outputDir, chStatus)
	} else if rxGit.MatchString(url) {
		err = DownloadGit(url, outputDir)
	} else if matchesGallery(url) {
		err = DownloadGallery(url, outputDir)
	} else if matchesVideo(url) {
		err = DownloadVideo(url, outputDir)
	} else if matchesMegaURL(url) {
		err = downloadMega(url, outputDir)
	} else {
		// Fallback
		err = downloadFile(url, outputDir)
	}
	if err != nil {
		ERROR(err)
		ERROR("Unsupported URL:\n\t%s\n", url)
		// If nothing else worked, add the URL to a queue file.
		Queue(url)
	}
}

// initialiseDirectories ensures that the given directories exist.
func initialiseDirectories(dir ...string) {
	for _, d := range dir {
		err := os.MkdirAll(d, os.ModeDir)
		FATAL(err)
	}
}

func main() {
	wr, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	FATAL(err)
	defer func() {
		if wr != nil {
			DEBUG("Closing log file")
			err = wr.Close()
			FATAL(err)
		}
	}()
	logWriter := bufio.NewWriter(wr)
	log.SetOutput(logWriter)
	log.SetFlags(log.Lmicroseconds | log.LstdFlags)

	flag.Usage = printHelp
	flag.Parse()
	if *help {
		printHelp()
		os.Exit(0)
	}

	initialiseDirectories(outputDir, queueDir)

	if *generateLists {
		generateGalleryList()
		generateVideoList()
	}

	if len(flag.Args()) == 0 {
		printHelp()
		os.Exit(0)
	}

	if flag.Args()[0] == "queue" {
		Queue(flag.Args()[1:]...)
		os.Exit(0)
	}

	wait := sync.WaitGroup{}
	for _, url := range flag.Args() {
		wait.Go(func() { Download(url, chStatus) })
	}
	wait.Wait()

	DEBUG("Done")
}
