package main

import (
	"fmt"

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

	wait := StartServer(*Port)
	fmt.Println()
	Interactive()
	for {
		select {
		case <-wait:
			DEBUG("Done")
		}
	}
}
