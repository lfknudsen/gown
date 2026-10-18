package main

import (
	"flag"
	. "gown/pkg"
	. "gown/pkg/web"
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
	}

	wait := StartServer(8080)
	Interactive(Download)
	for {
		select {
		case <-wait:
			return
		}
	}

	DEBUG("Done")
}
