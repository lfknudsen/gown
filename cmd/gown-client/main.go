package main

import (
	"flag"
	. "gown/pkg"
	. "gown/pkg/web"
)

func main() {
	logFile := InitialiseLogging()
	defer CloseLogFile(logFile)
	ParseFlags()

	client := NewGownClient(8080)
	for _, url := range flag.Args() {
		client.RequestDownload(url)
	}

}
