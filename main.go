package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/IronWill79/pokedex-go/internal/api"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *api.LocationConfig) error
}

var commands map[string]cliCommand

var conf api.LocationConfig

func main() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 Pokemon world map location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 Pokemon world map location areas",
			callback:    commandMapb,
		},
	}
	conf = api.LocationConfig{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}

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
		if _, ok := commands[command]; ok {
			commands[command].callback(&conf)
		} else {
			fmt.Println("Unknown command")
		}
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

func commandExit(c *api.LocationConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *api.LocationConfig) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, v := range commands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(c *api.LocationConfig) error {
	areas, err := api.GetNextLocations(c)
	if err != nil {
		return err
	}
	for _, area := range areas {
		fmt.Printf("%s\n", area)
	}
	return nil
}

func commandMapb(c *api.LocationConfig) error {
	areas, err := api.GetPreviousLocations(c)
	if err != nil {
		return err
	}
	for _, area := range areas {
		fmt.Printf("%s\n", area)
	}
	return nil
}
