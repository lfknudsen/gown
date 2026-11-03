package main

import (
	. "gown/pkg"
	. "gown/pkg/web"

	flag "github.com/spf13/pflag"
)

func main() {
	logFile := InitialiseLogging()
	defer CloseLogFile(logFile)
	ParseFlags()

	client := NewGownClient(*Port)
	for _, url := range flag.Args() {
		client.RequestDownload(url)
	}

}
