package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	flag "github.com/spf13/pflag"

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
			_, err := Queue(flag.Args()[1:]...)
			FATAL(err)
		}
		os.Exit(0)
	}

	if len(flag.Args()) >= 1 && (flag.Args()[0] == "interactive" || flag.Args()[0] == "repl") {
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

	filenames, err := Queue(flag.Args()...)
	FATAL(err)
	if *Singlethreaded {
		for i, filePath := range filenames {
			fmt.Printf("(%s) ", util.FormatFraction(i+1, len(flag.Args())))
			err = DownloadQueueItem(filePath)
			ERROR(err)
		}
	} else {
		wait := sync.WaitGroup{}
		for _, filePath := range filenames {
			wait.Go(func() {
				err = DownloadQueueItem(filePath)
				ERROR(err)
			})
		}
		wait.Wait()
	}

	DEBUG("Done")
}
