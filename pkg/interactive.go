package pkg

import (
	"fmt"
	"strings"

	"github.com/lmorg/readline/v4"
)

func Interactive(urlHandler func(string, chan int)) {
	instance := readline.NewInstance()
	instance.SetPrompt(">> URL: ")
	line, err := instance.Readline()
	for err == nil {
		line = strings.TrimSpace(line)
		switch line {
		case "exit", "quit":
			return
		default:
			urls := strings.Split(line, " ")
			for _, url := range urls {
				urlHandler(url, ChStatus)
			}
			line, err = instance.Readline()
		}

	}
	fmt.Println(err)
}
