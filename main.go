package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	var cleanedInput []string
	text = strings.TrimSpace(text)
	for word := range strings.SplitSeq(text, " ") {
		word = strings.TrimSpace(word)
		if len(word) > 0 {
			cleanedInput = append(cleanedInput, strings.ToLower(word))
		}
	}
	return cleanedInput
}
