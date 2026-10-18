package main

import (
	"flag"
	. "gown/pkg"
)

func main() {
	flag.Parse()
	Queue(flag.Args()...)
}
