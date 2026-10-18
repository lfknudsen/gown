package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const VERSION = "0.3.0"

func getExecutableDir() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(ex)
}

var root = getExecutableDir()
var outputDir string = filepath.Join(root, "downloads")

var chStatus = make(chan int, 1)

func FmtUrl(url string) string {
	url = strings.TrimSpace(url)
	url = strings.TrimSuffix(url, "/")
	url = regexp.MustCompile("http://").ReplaceAllString(url, "https://")
	return url
}

func Download(target string, chStatus chan int) {
	var err error = nil
	if strings.Contains(target, "mediafire.com") {
		err = downloadMediafire(target, outputDir, chStatus)
	} else if rxGit.MatchString(target) {
		err = DownloadGit(target, outputDir)
	} else if matchesGallery(target) {
		err = DownloadGallery(target, outputDir)
	} else if matchesVideo(target) {
		err = DownloadVideo(target, outputDir)
	} else if matchesMegaURL(target) {
		err = downloadMega(target, outputDir)
	} else {
		// Fallback
		err = downloadFile(target, outputDir)
	}
	if err != nil {
		log.Println(err)
		// If nothing else worked, add the URL to a queue file.
		log.Printf("Unsupported URL:\n\t%s\n", target)
		Queue(target)
	}
}

// initialiseDirectories ensures that the given directories exist.
func initialiseDirectories(dir ...string) {
	for _, d := range dir {
		err := os.MkdirAll(d, os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

var help = flag.Bool("h", false, "Print this message and exit.")

var overwrite = flag.Bool("overwrite", false,
	"Overwrite existing files rather than choosing a unique name for the new one.")

var generateLists = flag.Bool("generate-lists", false,
	"Generate/regenerate the lists of supported websites for gallery and video downloads.")

func printHelp() {
	fmt.Printf("Gown v%s\n\nUsage: %s [OPTIONS] <URL(s)...>\n\nOptions:\n", VERSION, os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = printHelp
	flag.Parse()
	if *help {
		printHelp()
		os.Exit(0)
	}

	initialiseDirectories(outputDir, queueDir)

	if *generateLists {
		generateGalleryList()
		generateVideoLists()
	}

	if len(flag.Args()) > 0 {
		wait := sync.WaitGroup{}
		for _, url := range flag.Args() {
			wait.Go(func() { Download(url, chStatus) })
		}
		wait.Wait()
	} else {
		printHelp()
		os.Exit(0)
	}
}
