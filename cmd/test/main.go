package main

import (
	"fmt"
	"log"

	"github.com/ufraaan/parser"
)

func main() {
	page, body, err := parser.Parse("https://ufraan.dev")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Page: %#v", page)
	fmt.Println("Body Text: ", body)
}
