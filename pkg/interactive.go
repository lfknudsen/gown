package pkg

import (
	"fmt"
	"strings"

	"github.com/lmorg/readline/v4"
)

func Interactive() {
	instance := readline.NewInstance()
	instance.SetPrompt(">> URL: ")
	line, err := instance.Readline()
	for err == nil {
		switch line {
		case "exit", "quit":
			return
		default:
			urls := strings.Split(line, " ")
			for _, url := range urls {
				Download(url, make(chan int))
			}
			line, err = instance.Readline()
		}

	}
	fmt.Println(err)
}
