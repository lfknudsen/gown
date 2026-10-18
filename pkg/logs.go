package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"uuid"
)

func DEBUG(msg any) {
	if msg != nil {
		if *Verbose >= 3 {
			fmt.Println(msg)
		}
		log.Println(msg)
	}
}

func INFO(msg any) {
	if msg != nil {
		if *Verbose >= 2 {
			fmt.Println(msg)
		}
		log.Println(msg)
	}
}

func WARN(msg any) {
	if msg != nil {
		if *Verbose >= 1 {
			fmt.Println(msg)
		}
		log.Println(msg)
	}
}

func ERROR(msg any) {
	if msg != nil {
		if *Verbose >= 1 {
			_, _ = fmt.Fprintln(os.Stderr, msg)
		}
		log.Println(msg)
	}
}

func FATAL(msg any) {
	if msg != nil {
		if *Verbose >= 1 {
			_, _ = fmt.Fprintln(os.Stderr, msg)
		}
		log.Fatalln(msg)
	}
}

func PANIC(msg any) {
	if msg != nil {
		if *Verbose >= 1 {
			_, _ = fmt.Fprintln(os.Stderr, msg)
		}
		log.Panicln(msg)
	}
}

const logDirName = "logs"

var LogDir = filepath.Join(GownDirectory, logDirName)

var LogFile = filepath.Join(LogDir, "log_"+uuid.New().String()+".txt")
