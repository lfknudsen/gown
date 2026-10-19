package pkg

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func query(url string, tag string) (string, error) {
	output := ""
	resp, err := http.Get(url)
	if err != nil {
		return output, err
	}
	defer closeReadCloser(resp.Body)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return output, err
	}

	return strings.TrimSpace(doc.Find(tag).Text()), err
}
