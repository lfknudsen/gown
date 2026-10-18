package pkg

import (
	"fmt"

	"github.com/lmorg/readline/v4"
)

func Interactive() {
	instance := readline.NewInstance()
	instance.SetPrompt(">> URL: ")
	line, err := instance.Readline()
	for err == nil {
		Download(line, make(chan int))
		line, err = instance.Readline()
	}
	fmt.Println(err)
}
