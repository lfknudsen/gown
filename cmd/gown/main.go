package main

import (
	"flag"
	. "gown/pkg"
	"os"
	"sync"
)

func main() {
	logFile := InitialiseLogging()
	defer CloseLogFile(logFile)
	ParseFlags()
	SetupDirs()

	if *GenerateLists {
		GenerateGalleryList()
		GenerateVideoList()
		WARN("Finished generating lists of supported domains.")
		if len(flag.Args()) == 0 {
			os.Exit(0)
		}
	}

	if len(flag.Args()) == 0 {
		PrintHelp()
		os.Exit(0)
	}

	if flag.Args()[0] == "queue" {
		if len(flag.Args()) == 2 && flag.Args()[1] == "download" {
			DownloadQueue()
		} else if len(flag.Args()) == 2 && flag.Args()[1] == "list" {
			ListQueue()
		} else {
			Queue(flag.Args()[1:]...)
		}
		os.Exit(0)
	}

	if len(flag.Args()) >= 1 && flag.Args()[0] == "interactive" {
		Interactive(Download)
		os.Exit(0)
	}

	if *Singlethreaded {
		for _, url := range flag.Args() {
			Download(url, ChStatus)
		}
	} else {
		wait := sync.WaitGroup{}
		for _, url := range flag.Args() {
			wait.Go(func() { Download(url, ChStatus) })
		}
		wait.Wait()
	}

	DEBUG("Done")
}
