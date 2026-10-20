package pkg

import (
	"crypto/md5"
	"errors"
	"fmt"
	"gown/pkg/util"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hymkor/trash-go"
	"github.com/manifoldco/promptui"
)

var QueueDir string = filepath.Join(GownDirectory, "queue")

func Queue(url ...string) error {
	for _, url := range url {
		hasher := md5.New()
		hasher.Write([]byte(url))
		hashOfURL := hasher.Sum(nil)
		filename := fmt.Sprintf("%s-%x", util.TimeStamp(), hashOfURL)
		filename = filepath.Join(QueueDir, filename+".txt")
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		_, err = f.WriteString(url + "\n")
		closeFile(f)
		if err != nil {
			return err
		}
		WARN(fmt.Sprintf("Added to queue: %s", url))
	}
	return nil
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

// getQueueItem opens the file with the given filename (not path!) in the queue directory,
// and returns its contents as a string.
func getQueueItem(filename string) (string, error) {
	filename = filepath.Join(QueueDir, filename)
	data, err := os.ReadFile(filename)
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
		filePath := filepath.Join(QueueDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			ERROR(err)
			continue
		}
		url := strings.TrimSpace(string(data))
		err = DownloadFromQueue(url, chStatus)
		if err == nil {
			_ = trash.Throw(filePath)
		}
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
		data, err := getQueueItem(entry.Name())
		if err != nil {
			ERROR(err)
			continue
		}
		url := strings.TrimSpace(string(data))
		fmt.Println(url)
	}
}

// ListPrompt opens an interactive list of the queue. The user
// can download or delete individual entries.
// Each entry will be shown as "<URL>   <Page title>"
func ListPrompt() {
	var err error
	entries, err := os.ReadDir(QueueDir)
	FATAL(err)
	if len(entries) == 0 {
		fmt.Println("Queue is empty")
		return
	}
	elements := make([]string, 0, len(entries))
	maxLength := 0
	// Build array of the URLs contained within the queue files.
	for _, entry := range entries {
		data, err := getQueueItem(entry.Name())
		if err != nil {
			// Couldn't find/read a file with the given name.
			ERROR(err)
			continue
		}
		url := strings.TrimSpace(string(data))
		if maxLength < len(url) {
			maxLength = len(url)
		}
		elements = append(elements, url)
	}
	if *PrintTitle {
		// Padding
		maxLength += 3
		// Concatenate each URL with its page title
		for i := range elements {
			// If the URL contains a filename, it possibly won't have an HTML page,
			// but also, you'll already be able to tell what the URL points to.
			title := ""
			if strings.ContainsRune(elements[i], '.') {
				title = path.Base(elements[i])
			} else {
				title, err = query(elements[i], "title")
				if err != nil {
					continue
				}
			}
			elements[i] = elements[i] + strings.Repeat(" ", maxLength-len(elements[i])) + title
		}
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
			Size:      10,
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
				err = DownloadFromQueue(elements[entryIndex], ChStatus)
				if err == nil {
					elements = slices.Delete(elements, entryIndex, entryIndex+1)
					_ = trash.Throw(filepath.Join(QueueDir, entries[entryIndex].Name()))
				}
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
