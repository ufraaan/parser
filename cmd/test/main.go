package main

import (
	"fmt"

	"github.com/ufraaan/parser"
)

func main() {
	text := "Hello, world!! Go_lang is-fast. 123_test, okay???"
	arr, err := parser.Tokenise(text) 
	if err != nil {
		fmt.Print("error: ", err)
	}
	fmt.Println(arr)
}


