package main

import (
	"flag"
	. "gown/pkg"
)

func main() {
	logFile := InitialiseLogging()
	defer CloseLogFile(logFile)
	flag.Parse()
	Queue(flag.Args()...)
}
