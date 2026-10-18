package main

import (
	"io"
	"os"
)

func closeFile(f *os.File) {
	if f != nil {
		err := f.Close()
		if err != nil {
			ERROR(err)
		}
	}
}

func closeReadCloser(body io.ReadCloser) {
	if body != nil {
		err := body.Close()
		if err != nil {
			ERROR(err)
		}
	}
}
