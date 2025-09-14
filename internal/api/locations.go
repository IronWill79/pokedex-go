package api

import (
	"encoding/json"

	"github.com/IronWill79/pokedex-go/internal/pokecache"
)

type LocationConfig struct {
	Next     string
	Previous string
}

type locationEndpoint struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type locationResults struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous *string            `json:"previous"`
	Results  []locationEndpoint `json:"results"`
}

type encounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type encounterVersionDetails struct {
	Rate    int     `json:"rate"`
	Version version `json:"version"`
}

type encounterMethodRates struct {
	EncounterMethod encounterMethod           `json:"encounter_method"`
	VersionDetails  []encounterVersionDetails `json:"version_details"`
}

type location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type names struct {
	Language language `json:"language"`
	Name     string   `json:"name"`
}

type pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type method struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type encounterDetails struct {
	Chance          int    `json:"chance"`
	ConditionValues []any  `json:"condition_values"`
	MaxLevel        int    `json:"max_level"`
	Method          method `json:"method"`
	MinLevel        int    `json:"min_level"`
}

type pokemonEncounterVersionDetails struct {
	EncounterDetails []encounterDetails `json:"encounter_details"`
	MaxChance        int                `json:"max_chance"`
	Version          version            `json:"version"`
}

type pokemonEncounters struct {
	Pokemon        pokemon                          `json:"pokemon"`
	VersionDetails []pokemonEncounterVersionDetails `json:"version_details"`
}

type locationAreaResults struct {
	EncounterMethodRates []encounterMethodRates `json:"encounter_method_rates"`
	GameIndex            int                    `json:"game_index"`
	ID                   int                    `json:"id"`
	Location             location               `json:"location"`
	Name                 string                 `json:"name"`
	Names                []names                `json:"names"`
	PokemonEncounters    []pokemonEncounters    `json:"pokemon_encounters"`
}

func GetLocationResults(url string, cache *pokecache.Cache) (locationResults, error) {
	var results locationResults
	res, _ := cache.Get(url)
	if err := json.Unmarshal(res, &results); err != nil {
		return results, err
	}
	return results, nil
}

func GetNextLocations(c *LocationConfig, cache *pokecache.Cache) ([]string, error) {
	var result []string
	res, err := GetLocationResults(c.Next, cache)
	if err != nil {
		return nil, err
	}
	if res.Previous != nil {
		c.Previous = string(*res.Previous)
	} else {
		c.Previous = c.Next
	}
	c.Next = res.Next
	for _, r := range res.Results {
		result = append(result, r.Name)
	}
	return result, nil
}

func GetPreviousLocations(c *LocationConfig, cache *pokecache.Cache) ([]string, error) {
	var result []string
	if c.Previous == "" {
		result = append(result, "you're on the first page")
		return result, nil
	}
	res, err := GetLocationResults(c.Previous, cache)
	if err != nil {
		return nil, err
	}
	c.Next = res.Next
	if res.Previous != nil {
		c.Previous = string(*res.Previous)
	} else {
		c.Previous = ""
	}
	for _, r := range res.Results {
		result = append(result, r.Name)
	}
	return result, nil
}

func GetPokemonFromArea(area string, cache *pokecache.Cache) ([]string, error) {
	var result []string
	var results locationAreaResults
	res, _ := cache.Get("https://pokeapi.co/api/v2/location-area/" + area)
	if err := json.Unmarshal(res, &results); err != nil {
		return result, nil
	}
	for _, pokemon := range results.PokemonEncounters {
		result = append(result, pokemon.Pokemon.Name)
	}
	return result, nil
}
