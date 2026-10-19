package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/lmorg/readline/v4"
)

func interactiveHelp() {
	fmt.Println("Commands:\n" +
		"download <URL(s)...>    Download the given URL(s).\n" +
		"queue <URL(s)...>       Add the given URL(s) to the queue.\n" +
		"help                    Print this message (alias: h).\n" +
		"exit                    Close the programme (aliases: quit, q, :q).")
}

func Interactive() {
	fmt.Println("Gown interactive environment started.")
	interactiveHelp()
	instance := readline.NewInstance()
	instance.SetPrompt(">>> ")
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
		case "exit", "quit", "q", ":q":
			os.Exit(0)
		case "download", "dl", "down", "gown":
			for _, url := range words[1:] {
				err = Download(url, ChStatus)
				ERROR(err)
			}
		case "queue":
			for _, url := range words[1:] {
				err = Queue(url)
				ERROR(err)
			}
		default:
			fmt.Println("Unknown command. Try typing 'help' (without quotes).")
		}
		line, err = instance.Readline()
	}
	fmt.Println(err)
}
