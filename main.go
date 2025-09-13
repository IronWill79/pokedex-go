package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				fmt.Printf("Error reading input: %v", err)
			}
			fmt.Printf("EOF reached")
		}
		input := scanner.Text()
		if input == "" {
			continue
		}
		cleanedInput := cleanInput(input)
		command := cleanedInput[0]
		fmt.Printf("Your command was: %s", command)
	}
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
