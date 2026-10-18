package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"uuid"
)

var queueDir string = filepath.Join(root, "queue")

func Queue(url string) {
	filename := time.Now().UTC().Format(time.RFC3339)
	filename = strings.Replace(filename, ":", "_", -1)
	uid := uuid.New()
	filename += "-" + uid.String()
	filename = filepath.Join(queueDir, filename+".txt")
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(url + "\n")
	if err != nil {
		_ = f.Close()
		log.Fatal(err)
	}
	_ = f.Close()
}
