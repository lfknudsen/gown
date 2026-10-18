package main

import (
	"flag"
	"fmt"
	"os"
)

var help = flag.Bool("h", false, "Print this message and exit. (Alias: --help)")

var overwrite = flag.Bool("overwrite", false,
	"Overwrite existing files rather than choosing a unique name for the new one.")

var generateLists = flag.Bool("generate-lists", false,
	"Generate/regenerate the lists of supported websites for gallery and video downloads.")

var verbose = flag.Int("v", 1, "Set verbosity level. Higher = more messages. 0-3")

func printHelp() {
	fmt.Printf("Gown v%s\n\nUsage: %s [OPTIONS] <URL(s)...>\n\nOptions:\n", VERSION, os.Args[0])
	flag.PrintDefaults()
}
