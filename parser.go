package parser

import (
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	Title       string
	Description string
	Links       []string
}

func Parse(url string) (Page, string, error) {
	res, err := http.Get(url)
	if err != nil {
		return Page{}, "", err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return Page{}, "", err
	}
	page := Page{}

	// non content nodes
	doc.Find("script").Remove()
	doc.Find("style").Remove()

	// get the title
	page.Title = doc.Find("title").Text()

	// get the meta description content
	desc, exists := doc.Find("meta[name=description]").Attr("content")
	if exists {
		page.Description = desc
	}

	// get urls
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			page.Links = append(page.Links, href)
		}
	})

	body := doc.Find("body").Text()

	re := regexp.MustCompile(`\s+`)
	body = re.ReplaceAllString(body, " ")

	return page, body, nil
}
