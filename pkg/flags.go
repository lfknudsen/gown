package pkg

import (
	"flag"
	"fmt"
	"os"
)

var Help = flag.Bool("h", false, "Print this message and exit. (Alias: --help)")

var Overwrite = flag.Bool("overwrite", false,
	"Overwrite existing files rather than choosing a unique name for the new one.")

var GenerateLists = flag.Bool("generate-lists", false,
	"Generate/regenerate the lists of supported websites for gallery and video downloads.")

var Verbose = flag.Int("v", 1, "Set verbosity level. Higher = more messages. 0-3")

var Singlethreaded = flag.Bool("s", false, "If true, downloads are performed sequentially.")

func PrintHelp() {
	fmt.Printf("Gown v%s\n\nUsage: %s [OPTIONS] <URL(s)...>\n\nOptions:\n", VERSION, os.Args[0])
	flag.PrintDefaults()
}
