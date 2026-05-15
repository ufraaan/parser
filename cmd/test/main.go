package main

import (
	"fmt"
	"os"

	"github.com/ufraaan/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: parser <url>")
		os.Exit(1)
	}

	url := os.Args[1]
	page, body, err := parser.Parse(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("title:", page.Title)
	fmt.Println("description:", page.Description)
	fmt.Println("links:")
	for _, l := range page.Links {
		fmt.Println("  ", l)
	}

	tokens, _ := parser.Tokenise(body)
	fmt.Println("tokens:", tokens)
}


