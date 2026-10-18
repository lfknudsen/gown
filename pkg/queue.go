package pkg

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
	"uuid"

	"github.com/hymkor/trash-go"
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

func GetFirst() (string, error) {
	entries, err := os.ReadDir(QueueDir)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", errors.New("queue is empty")
	}
	data, err := os.ReadFile(entries[0].Name())
	return string(data), err
}

func DownloadQueue() {
	entries, err := os.ReadDir(QueueDir)
	FATAL(err)
	if len(entries) == 0 {
		return
	}
	chStatus := make(chan int)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			ERROR(err)
			continue
		}
		if !info.Mode().IsRegular() {
			continue
		}
		path := filepath.Join(QueueDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			ERROR(err)
			continue
		}
		url := strings.TrimSpace(string(data))
		Download(url, chStatus)
		err = trash.Throw(entry.Name())
		ERROR(err)
	}
}
