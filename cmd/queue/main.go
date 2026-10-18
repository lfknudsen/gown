package main

import (
	"flag"
	. "gown/pkg"
	"os"
)

func main() {
	logFile := InitialiseLogging()
	defer CloseLogFile(logFile)
	flag.Usage = PrintHelp
	flag.Parse()
	DEBUG("Parsed command-line flags.")
	if *Help {
		PrintHelp()
		os.Exit(0)
	}

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

	Queue(flag.Args()...)

	DEBUG("Done")
}
