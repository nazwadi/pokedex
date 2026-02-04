package main

import (
	"time"

	"github.com/nazwadi/pokedexcli/internal/pokeapi"
	"github.com/nazwadi/pokedexcli/internal/pokecache"
)

type Config struct {
	client   pokeapi.Client
	cache    *pokecache.Cache
	pokemon  map[string]pokeapi.Pokemon
	Next     *string
	Previous *string
}

func main() {
	client := pokeapi.NewClient(5 * time.Second)
	newCache := pokecache.NewCache(5 * time.Second)
	pokedex := make(map[string]pokeapi.Pokemon)

	cfg := &Config{
		client:   client,
		cache:    newCache,
		pokemon:  pokedex,
		Next:     nil,
		Previous: nil,
	}

	startRepl(cfg)
}
