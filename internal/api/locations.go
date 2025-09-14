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
