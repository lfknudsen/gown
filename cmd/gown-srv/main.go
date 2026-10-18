package main

import (
	. "gown/pkg"
	. "gown/pkg/web"
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
