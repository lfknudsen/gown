package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/lmorg/readline/v4"
)

func interactiveHelp() {
	fmt.Println("Commands:\n" +
		"download <URL(s)...>       Download the given URL(s) (alias: dl).\n" +
		"queue <URL(s)...>          Add the given URL(s) to the queue (alias: q).\n" +
		"download                   Enter download mode, where every line given will be\n" +
		"                           treated as if you'd prepended it with the download command.\n" +
		"                           All commands still function as normal in this mode, though.\n" +
		"                           Exit with 'cancel'.\n" +
		"queue                      Enter queue mode, where every line given will be\n" +
		"                           treated as if you'd prepended it with the queue command.\n" +
		"                           All commands still function as normal in this mode, though.\n" +
		"                           Exit with 'cancel'.\n" +
		"list					    List the contents of the queue.\n" +
		"help                       Print this message (alias: h).\n" +
		"exit                       Close the programme (aliases: quit, :q).")
}

const (
	lineMode = iota
	downloadMode
	queueMode
)

var promptText = map[int]string{
	lineMode:     "\033[94m>>>\033[0m ",
	downloadMode: "\033[32mDL\033[94m>\033[0m ",
	queueMode:    "\033[35mQ\033[94m>>\033[0m ",
}

func Interactive() {
	fmt.Println("Gown interactive environment started.")
	interactiveHelp()
	mode := lineMode
	instance := readline.NewInstance()
	instance.SetPrompt(promptText[mode])
	line, err := instance.Readline()
	for err == nil {
		words := strings.Fields(line)
		if len(words) == 0 {
			line, err = instance.Readline()
			continue
		}
		switch words[0] {
		case "help", "h", ":h":
			interactiveHelp()
		case "exit", "quit", ":q":
			os.Exit(0)
		case "download", "dl", "down", "gown":
			if len(words[1:]) == 0 {
				mode = downloadMode
				instance.SetPrompt(promptText[downloadMode])
				break
			}
			for _, url := range words[1:] {
				err = Download(url, ChStatus)
				ERROR(err)
			}
		case "queue", "q":
			if len(words[1:]) == 0 {
				mode = queueMode
				instance.SetPrompt(promptText[queueMode])
				break
			}
			for _, url := range words[1:] {
				_, err = Queue(url)
				ERROR(err)
			}
		case "list":
			ListPrompt()
		case "cancel":
			mode = lineMode
			instance.SetPrompt(promptText[lineMode])
		default:
			switch mode {
			case downloadMode:
				filenames, err := Queue(words...)
				ERROR(err)
				for _, url := range filenames {
					err = DownloadQueueItem(url)
					ERROR(err)
				}
			case queueMode:
				for _, url := range words {
					_, err = Queue(url)
					ERROR(err)
				}
			default:
				fmt.Println("Unknown command. Try typing 'help' (without quotes).")
			}
		}
		line, err = instance.Readline()
	}
	fmt.Println(err)
}
