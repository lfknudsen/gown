package pkg

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var Help = flag.Bool("h", false, "Print this message and exit. (Alias: --help)")

var Overwrite = flag.Bool("overwrite", false,
	"Overwrite existing files rather than choosing a unique name for the new one.")

var GenerateLists = flag.Bool("generate-lists", false,
	"Generate/regenerate the lists of supported websites for gallery and video downloads.")

var Verbose = flag.Int("v", 1, "Set verbosity level. Higher = more messages. 0-3")

var Singlethreaded = flag.Bool("s", false, "If true, downloads are performed sequentially.")

var Retries = flag.Int("retries", 5, "Number of times to retry gallery download. -1 = infinite.")
var RetriesStr = strconv.Itoa(*Retries)

func PrintHelp() {
	fmt.Printf("Gown v%s\n\nUsage: %s [OPTIONS] <URL(s)...>\n\nOptions:\n", VERSION, os.Args[0])
	flag.PrintDefaults()
}

func ParseFlags() {
	flag.Usage = PrintHelp
	flag.Parse()
	if *Retries < -1 {
		*Retries = -1
	}
	RetriesStr = strconv.Itoa(*Retries)
	if *Help {
		PrintHelp()
		os.Exit(0)
	}
}
