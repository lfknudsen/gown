package main

import (
	"flag"
	. "gown/pkg"
)

func main() {
	InitialiseLogging()
	flag.Parse()
	Queue(flag.Args()...)
}
