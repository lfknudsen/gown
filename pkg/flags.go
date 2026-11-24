package pkg

import (
	"fmt"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"
)

var Help = flag.BoolP("help", "h", false,
	"Print this message and exit. (Alias: --help)")

var Overwrite = flag.BoolP("overwrite", "o", false,
	"Overwrite existing files rather than choosing a unique name for the new one.")

var GenerateLists = flag.BoolP("generate-lists", "g", false,
	"Generate/regenerate the lists of supported websites for gallery and video downloads.")

var Verbose = flag.IntP("verbose", "v", 1, "Set verbosity level. Higher = more messages. 0-3")

var Singlethreaded = flag.BoolP("sequential", "s", false, "If true, downloads are performed sequentially.")

var Retries = flag.IntP("retries", "r", 5, "Number of times to retry gallery download. -1 = infinite.")
var RetriesStr = strconv.Itoa(*Retries)

var Port = flag.IntP("port", "p", 8080, "The port to connect or listen to.")

var PrintTitle = flag.BoolP("title", "t", false, "When listing the queue, print each link's HTML title.")

var Using = flag.StringP("using", "u", "",
	"Use the download handler specified rather than delegating based on URL.\n"+
	"Accepted are: 'gallery', 'video', 'mediafire', 'mega', 'instagram', 'file'.\n"+
	"Falls back to file if argument is unrecognised.")

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

	if *Port < 1000 {
		_, _ = fmt.Fprintln(os.Stderr, "Port must be a four-digit positive integer.")
		PrintHelp()
		os.Exit(1)
	}
}
