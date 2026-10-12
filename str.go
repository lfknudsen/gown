package main

import (
	"path"
	"path/filepath"
	"strings"
)

func UrlStem(s string) string {
	base := path.Base(s)
	extIdx := strings.IndexRune(base, '.')
	if extIdx == -1 {
		return base
	}
	return base[:extIdx]
}

func FileStem(s string) string {
	return s[strings.LastIndexByte(s, filepath.Separator)+1 : strings.LastIndexByte(s, '.')]
}

func Ext(s string) string {
	return s[strings.IndexRune(s, '.'):]
}
