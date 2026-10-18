package main

import (
	"bufio"
	"flag"
	. "gown/pkg"
	"log"
	"os"
	"sync"
)

func main() {
	InitialiseDirectories(LogDir)
	wr, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	FATAL(err)
	defer func() {
		if wr != nil {
			DEBUG("Closing log file")
			err = wr.Close()
			FATAL(err)
		}
	}()
	DEBUG("Opened the new log file.")
	logWriter := bufio.NewWriter(wr)
	log.SetOutput(logWriter)
	log.SetFlags(log.Lmicroseconds | log.LstdFlags)
	DEBUG("Finished setting up log file.")

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
