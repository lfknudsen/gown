package util

import (
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
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

// UrlRoot throws away any "http[s]", "www", etc., as well as everything after
// the first / (inclusive).
// So "https://www.duckduckgo.com/search?" would become "duckduckgo.com"
func UrlRoot(s string) string {
	return regexp.MustCompile(`(?:https?://)?(?:ww[w\d]*\.)?([^/]+).*`).FindStringSubmatch(s)[1]
}

func TimeStamp() string {
	out := time.Now().UTC().Format(time.RFC3339)
	return strings.Replace(out, ":", "_", -1)
}

func FormatFractionU(left, right uint64) string {
	l := strconv.FormatUint(left, 10)
	r := strconv.FormatUint(right, 10)
	difference := len(r) - len(l)
	prefix := strings.Repeat(" ", difference)
	return prefix + l + "/" + r
}
