package main

import . "gown/pkg"

func main() {
	logFile := InitialiseLogging()
	defer CloseLogFile(logFile)
	Interactive()
}
