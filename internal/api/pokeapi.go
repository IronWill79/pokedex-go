package api

import (
	"encoding/json"
	"time"

	"github.com/IronWill79/pokedex-go/internal/constants"
	"github.com/IronWill79/pokedex-go/internal/pokecache"
)

type PokeApi struct {
	cache *pokecache.Cache
}

func NewPokeApi() PokeApi {
	interval := 5 * time.Minute
	cache := pokecache.NewCache(interval)

	return PokeApi{
		cache: cache,
	}
}

func (p *PokeApi) GetLocationResults(url string) (LocationAreasResults, error) {
	var results LocationAreasResults
	res, _ := p.cache.Get(url)
	if err := json.Unmarshal(res, &results); err != nil {
		return results, err
	}
	return results, nil
}

func (p *PokeApi) GetNextLocations(c *LocationConfig) ([]string, error) {
	var result []string
	res, err := p.GetLocationResults(c.Next)
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

func (p *PokeApi) GetPreviousLocations(c *LocationConfig) ([]string, error) {
	var result []string
	if c.Previous == "" {
		result = append(result, "you're on the first page")
		return result, nil
	}
	res, err := p.GetLocationResults(c.Previous)
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

func (p *PokeApi) GetPokemonFromArea(area string) ([]string, error) {
	var result []string
	var results LocationArea
	res, _ := p.cache.Get(constants.LocationAreaEndpoint + area)
	if err := json.Unmarshal(res, &results); err != nil {
		return result, err
	}
	for _, pokemon := range results.PokemonEncounters {
		result = append(result, pokemon.Pokemon.Name)
	}
	return result, nil
}

func (p *PokeApi) GetPokemon(name string) (Pokemon, error) {
	var pokemon Pokemon
	res, _ := p.cache.Get(constants.PokemonEndpoint + name)
	if err := json.Unmarshal(res, &pokemon); err != nil {
		return pokemon, err
	}
	return pokemon, nil
}
