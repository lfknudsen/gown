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

	client := NewGownClient(8080)
	for _, url := range flag.Args() {
		client.RequestDownload(url)
	}

}
