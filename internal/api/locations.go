package api

import (
	"encoding/json"
	"net/http"
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

func GetLocations(url string) (locationResults, error) {
	var results locationResults
	res, err := http.Get(url)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&results); err != nil {
		return results, err
	}
	return results, nil
}

func GetNextLocations(c *LocationConfig) ([]string, error) {
	var result []string
	res, err := GetLocations(c.Next)
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

func GetPreviousLocations(c *LocationConfig) ([]string, error) {
	var result []string
	if c.Previous == "" {
		result = append(result, "you're on the first page")
		return result, nil
	}
	res, err := GetLocations(c.Previous)
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
