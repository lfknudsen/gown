package main

import (
	"fmt"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"

	. "gown/pkg"
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

	if len(flag.Args()) == 1 && flag.Args()[0] == "download" {
		DownloadQueue()
		os.Exit(0)
	}

	if len(flag.Args()) == 2 && flag.Args()[0] == "download" {
		num, err := strconv.ParseUint(flag.Args()[1], 10, 64)
		if err != nil {
			if *Verbose >= 1 {
				_, _ = fmt.Fprintf(os.Stderr,
					"Expected an integer argument, but got %s",
					flag.Args()[1])
			}
			os.Exit(1)
		}
		DownloadFirstN(num)
		os.Exit(0)
	}

	if len(flag.Args()) == 1 && flag.Args()[0] == "list" {
		ListPrompt()
		os.Exit(0)
	}

	err := Queue(flag.Args()...)
	FATAL(err)

	DEBUG("Done")
}
