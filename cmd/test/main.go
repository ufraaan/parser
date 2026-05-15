package main

import (
	"fmt"
	"log"

	"github.com/ufraaan/parser"
)

func main() {
	title, err := parser.Parse("https://ufraan.dev")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(title)
}
