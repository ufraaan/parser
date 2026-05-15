package parser

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

func Parse(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	title := doc.Find("title").Text()

	return title, nil
}


