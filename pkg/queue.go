package pkg

import (
	"os"
	"path/filepath"
	"strings"
	"time"
	"uuid"
)

var QueueDir string = filepath.Join(GownDirectory, "queue")

func Queue(url ...string) {
	for _, url := range url {
		filename := time.Now().UTC().Format(time.RFC3339)
		filename = strings.Replace(filename, ":", "_", -1)
		uid := uuid.New()
		filename += "-" + uid.String()
		filename = filepath.Join(QueueDir, filename+".txt")
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		FATAL(err)
		_, err = f.WriteString(url + "\n")
		closeFile(f)
		FATAL(err)
	}
}
