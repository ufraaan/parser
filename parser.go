package parser

import (
	"net/http"
	"regexp"
	"strings"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	Title       string
	Description string
	Links       []string
}

// done only once when package loads
var re = regexp.MustCompile(`\s+`)

func Parse(url string) (Page, string, error) {
	res, err := http.Get(url)
	if err != nil {
		return Page{}, "", err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return Page{}, "", fmt.Errorf("bad status: %d, %s", res.StatusCode, res.Status)
	}

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

	var bodyParts []string
	doc.Find("h1, h2, h3, p, li").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		bodyParts = append(bodyParts, text)
	})

	body = strings.Join(bodyParts, " ")

	// cleanup
	body = re.ReplaceAllString(body, " ")

	return page, body, nil
}
