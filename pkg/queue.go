package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"uuid"

	"github.com/hymkor/trash-go"
	"github.com/manifoldco/promptui"
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
		WARN(fmt.Sprintf("Added to queue: %s", url))
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
		err = trash.Throw(path)
		ERROR(err)
	}
}

func ListQueue() {
	entries, err := os.ReadDir(QueueDir)
	FATAL(err)
	if len(entries) == 0 {
		fmt.Println("Queue is empty")
		return
	}
	for _, entry := range entries {
		path := filepath.Join(QueueDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			ERROR(err)
			continue
		}
		url := strings.TrimSpace(string(data))
		fmt.Println(url)
	}
}

func ListPrompt() {
	var err error
	entries, err := os.ReadDir(QueueDir)
	FATAL(err)
	if len(entries) == 0 {
		fmt.Println("Queue is empty")
		return
	}
	elements := make([]string, 0, len(entries))
	for _, entry := range entries {
		path := filepath.Join(QueueDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			ERROR(err)
		}
		url := strings.TrimSpace(string(data))
		elements = append(elements, url)
	}

	promptAction := promptui.Select{
		Label:     "Select action",
		Items:     []string{"Download", "Delete", "Cancel"},
		IsVimMode: false,
	}

	var entryIndex int = 0
	var actionIndex int = -1
	scrollPosition := 0
	for {
		promptList := promptui.Select{
			Label:     "Select queue element",
			Items:     elements,
			IsVimMode: false,
		}
		entryIndex, _, err = promptList.RunCursorAt(entryIndex, scrollPosition)
		if err != nil {
			return
		}
		scrollPosition = promptList.ScrollPosition()
		actionIndex, _, err = promptAction.Run()
		if actionIndex >= 0 {
			switch actionIndex {
			case 0:
				Download(elements[entryIndex], ChStatus)
				elements = slices.Delete(elements, entryIndex, entryIndex+1)
				_ = trash.Throw(filepath.Join(QueueDir, entries[entryIndex].Name()))
			case 1:
				elements = slices.Delete(elements, entryIndex, entryIndex+1)
				_ = trash.Throw(filepath.Join(QueueDir, entries[entryIndex].Name()))
			}
		}
		if err != nil {
			break
		}
	}

}
