package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/IronWill79/pokedex-go/internal/api"
	"github.com/IronWill79/pokedex-go/internal/constants"
)

type cliCommand struct {
	name        string
	description string
	callback    func(param string, c *api.LocationConfig) error
}

var commands map[string]cliCommand

var pokeApi api.PokeApi

var pokeDex map[string]api.Pokemon

func main() {
	commands = map[string]cliCommand{
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Try to catch a Pokemon",
			callback:    commandCatch,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore <area_name>",
			description: "Show Pokemon in location area",
			callback:    commandExplore,
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
	pokeApi = api.NewPokeApi()
	pokeDex = make(map[string]api.Pokemon)
	conf := api.LocationConfig{
		Next:     constants.LocationAreaEndpoint,
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
		param := ""
		if len(cleanedInput) > 1 {
			param = cleanedInput[1]
		}
		if _, ok := commands[command]; ok {
			commands[command].callback(param, &conf)
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

func commandExit(param string, c *api.LocationConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(param string, c *api.LocationConfig) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, v := range commands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(param string, c *api.LocationConfig) error {
	areas, err := pokeApi.GetNextLocations(c)
	if err != nil {
		return err
	}
	for _, area := range areas {
		fmt.Printf("%s\n", area)
	}
	return nil
}

func commandMapb(param string, c *api.LocationConfig) error {
	areas, err := pokeApi.GetPreviousLocations(c)
	if err != nil {
		return err
	}
	for _, area := range areas {
		fmt.Printf("%s\n", area)
	}
	return nil
}

func commandExplore(param string, c *api.LocationConfig) error {
	fmt.Printf("Exploring %s...\n", param)
	fmt.Println("Found Pokemon:")
	pokemon, err := pokeApi.GetPokemonFromArea(param)
	if err != nil {
		return err
	}
	for _, name := range pokemon {
		fmt.Printf("- %s\n", name)
	}
	return nil
}

func commandCatch(param string, c *api.LocationConfig) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", param)
	pokemon, err := pokeApi.GetPokemon(param)
	if err != nil {
		return err
	}
	chance := rand.Intn(1000)
	if chance > pokemon.BaseExperience {
		pokeDex[param] = pokemon
		fmt.Printf("%s was caught!\n", param)
	} else {
		fmt.Printf("%s escaped!\n", param)
	}
	return nil
}
