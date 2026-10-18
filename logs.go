package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"uuid"
)

func DEBUG(msg any) {
	if msg != nil {
		if *verbose >= 3 {
			fmt.Println(msg)
		}
		log.Println(msg)
	}
}

func INFO(msg any) {
	if msg != nil {
		if *verbose >= 2 {
			fmt.Println(msg)
		}
		log.Println(msg)
	}
}

func WARN(msg any) {
	if msg != nil {
		if *verbose >= 1 {
			fmt.Println(msg)
		}
		log.Println(msg)
	}
}

func ERROR(msg any) {
	if msg != nil {
		if *verbose >= 1 {
			_, _ = fmt.Fprintln(os.Stderr, msg)
		}
		log.Println(msg)
	}
}

func FATAL(msg any) {
	if msg != nil {
		if *verbose >= 1 {
			_, _ = fmt.Fprintln(os.Stderr, msg)
		}
		log.Fatalln(msg)
	}
}

func PANIC(msg any) {
	if msg != nil {
		if *verbose >= 1 {
			_, _ = fmt.Fprintln(os.Stderr, msg)
		}
		log.Panicln(msg)
	}
}

const logDirName = "logs"

var logDir = filepath.Join(gownDirectory, logDirName)

var logFile = filepath.Join(logDir, "log_"+uuid.New().String()+".txt")
