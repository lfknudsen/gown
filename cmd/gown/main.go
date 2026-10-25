package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"

	. "gown/pkg"
	"gown/pkg/util"
	"gown/pkg/web"
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
		} else if len(flag.Args()) == 3 && flag.Args()[1] == "download" {
			num, err := strconv.ParseUint(flag.Args()[2], 10, 64)
			if err != nil {
				if *Verbose >= 1 {
					_, _ = fmt.Fprintf(os.Stderr,
						"Expected an integer argument, but got %s",
						flag.Args()[1])
				}
				os.Exit(1)
			}
			DownloadFirstN(num)
		} else if len(flag.Args()) == 2 && flag.Args()[1] == "list" {
			ListPrompt()
		} else {
			err := Queue(flag.Args()[1:]...)
			FATAL(err)
		}
		os.Exit(0)
	}

	if len(flag.Args()) >= 1 && flag.Args()[0] == "interactive" {
		Interactive()
		os.Exit(0)
	}

	if len(flag.Args()) >= 1 && flag.Args()[0] == "server" {
		wait := web.StartServer(*Port)
		select {
		case <-wait:
			os.Exit(0)
		}
	}

	if len(flag.Args()) >= 1 && flag.Args()[0] == "client" {
		client := web.NewGownClient(*Port)
		for _, domain := range flag.Args()[1:] {
			client.RequestDownload(domain)
		}
	}

	if *Singlethreaded {
		for i, url := range flag.Args() {
			fmt.Printf("(%s) ", util.FormatFraction(i+1, len(flag.Args())))
			err := Download(url, ChStatus)
			ERROR(err)
		}
	} else {
		wait := sync.WaitGroup{}
		for _, url := range flag.Args() {
			wait.Go(func() {
				err := Download(url, ChStatus)
				ERROR(err)
			})
		}
		wait.Wait()
	}

	DEBUG("Done")
}
