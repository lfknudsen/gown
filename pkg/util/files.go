package util

import (
	"io"
	"os"
)

func CloseFile(f *os.File) {
	if f != nil {
		_ = f.Close()
	}
}

func CloseReadCloser(body io.ReadCloser) {
	if body != nil {
		_ = body.Close()
	}
}
