package parser

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	Title       string
	Description string
	Links       []string
}

// done only once when package loads
var re = regexp.MustCompile(`\s+`)

func normalise(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	// normalise scheme
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	// remove fragment
	u.Fragment = ""

	// clean query
	u.RawQuery = ""

	// normalise host
	u.Host = strings.ToLower(u.Hostname()) // Hostname strips port

	// normalise path
	if u.Path == "" {
		u.Path = "/"
	}

	return u.String(), nil
}

func Parse(rawUrl string) (Page, string, error) {

	// normalise url
	normalisedUrl, err := normalise(rawUrl)
	if err != nil {
		return Page{}, "", err
	}

	baseURL, err := url.Parse(normalisedUrl)
	if err != nil {
		return Page{}, "", err
	}

	res, err := http.Get(normalisedUrl)
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
			linkURL, err := url.Parse(href)
			if err == nil {
				resolved := baseURL.ResolveReference(linkURL).String()
				page.Links = append(page.Links, resolved)
			}
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

func Tokenise(text string) ([]string, error) {
	text = strings.ToLower(text)
	
	words := regexp.MustCompile(`[^\w]+`).Split(text, -1)
	
	var tokens []string
	for _, w := range words {
		if len(w) > 2 {
			tokens = append(tokens, w)
		}
	}
	
	return tokens, nil
}