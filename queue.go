package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"
	"uuid"
)

var queueDir string = filepath.Join(gownDirectory, "queue")

func Queue(url ...string) {
	DEBUG("Entered Queue()")
	for _, url := range url {
		filename := time.Now().UTC().Format(time.RFC3339)
		filename = strings.Replace(filename, ":", "_", -1)
		uid := uuid.New()
		filename += "-" + uid.String()
		filename = filepath.Join(queueDir, filename+".txt")
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		FATAL(err)
		_, err = f.WriteString(url + "\n")
		closeFile(f)
		FATAL(err)
	}
}
