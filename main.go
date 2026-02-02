package main

import (
	"time"

	"github.com/nazwadi/pokedexcli/internal/pokeapi"
)

type Config struct {
	client   pokeapi.Client
	Next     *string
	Previous *string
}

func main() {
	client := pokeapi.NewClient(5 * time.Second)

	cfg := &Config{
		client:   client,
		Next:     nil,
		Previous: nil,
	}

	startRepl(cfg)
}
