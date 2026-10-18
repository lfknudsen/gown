package main

import (
	"flag"
	. "gown/pkg"
	"os"
	"sync"
)

func main() {
	InitialiseLogging()
	flag.Usage = PrintHelp
	flag.Parse()
	DEBUG("Parsed command-line flags.")
	if *Help {
		PrintHelp()
		os.Exit(0)
	}

	InitialiseDirectories(OutputDir, QueueDir)

	if *GenerateLists {
		GenerateGalleryList()
		GenerateVideoList()
	}

	if len(flag.Args()) == 0 {
		PrintHelp()
		os.Exit(0)
	}

	if flag.Args()[0] == "queue" {
		Queue(flag.Args()[1:]...)
		os.Exit(0)
	}

	wait := sync.WaitGroup{}
	for _, url := range flag.Args() {
		wait.Go(func() { Download(url, ChStatus) })
	}
	wait.Wait()

	DEBUG("Done")
}
