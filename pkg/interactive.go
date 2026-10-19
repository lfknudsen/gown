package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/lmorg/readline/v4"
)

func interactiveHelp() {
	fmt.Println("Enter a URL and press ENTER to download it.\n\n" +
		"Commands:\n" +
		"help\t\tPrint this message (alias: h).\n" +
		"exit\t\tClose the programme (aliases: quit, q, :q).")
}
func Interactive(urlHandler func(string, chan int) error) {
	fmt.Println("Interactive environment started.")
	interactiveHelp()
	instance := readline.NewInstance()
	instance.SetPrompt(">> URL: ")
	line, err := instance.Readline()
	for err == nil {
		line = strings.TrimSpace(line)
		switch line {
		case "":
			break
		case "help", "h":
			interactiveHelp()
		case "exit", "quit", ":q", "q":
			os.Exit(0)
		default:
			urls := strings.Split(line, " ")
			for _, url := range urls {
				err = urlHandler(url, ChStatus)
				ERROR(err)
			}
		}
		line, err = instance.Readline()

	}
	fmt.Println(err)
}
