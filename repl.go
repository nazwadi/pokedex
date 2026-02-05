package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/nazwadi/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *Config, args []string) error
}

func startRepl(cfg *Config) {
	m := map[string]cliCommand{
		"catch": {
			name:        "catch",
			description: "Catch pokemon",
			callback:    commandCatch,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore",
			description: "Explore the Pokedex",
			callback:    commandExplore,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect the pokemon",
			callback:    commandInspect,
		},
		"map": {
			name:        "map",
			description: "Map a Pokedex using the given arguments",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Map backwards a Pokedex using the given arguments",
			callback:    commandMapb,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Pokedex using the given arguments",
			callback:    commandPokedex,
		},
	}
	fmt.Println("Welcome to the Pokedex!")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		cleanedInput := cleanInput(scanner.Text())
		if len(cleanedInput) == 0 {
			continue
		}
		cmd := cleanedInput[0]
		args := cleanedInput[1:]

		value, ok := m[cmd]
		if ok {
			err := value.callback(cfg, args)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandCatch(cfg *Config, args []string) error {
	if len(args) == 0 {
		return errors.New("no pokemon name given")
	}
	pokemonName := args[0]
	var url = "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	var pokemon pokeapi.Pokemon
	var err error

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	data, ok := cfg.cache.Get(url)
	if ok {
		err := json.Unmarshal(data, &pokemon)
		if err != nil {
			return err
		}
	} else {
		pokemon, err = cfg.client.CatchPokemon(&url)
		if err != nil {
			return err
		}
		var jsonData []byte
		jsonData, err = json.Marshal(pokemon)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, jsonData)
	}
	res := rand.Intn(pokemon.BaseExperience)
	if res > 40 {
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemonName)
	fmt.Println("You may now inspect it with the inspect command.")
	cfg.pokemon[pokemon.Name] = pokemon
	return nil
}

func commandMap(cfg *Config, _ []string) error {
	var url string
	if cfg.Next != nil && *cfg.Next != "" {
		url = *cfg.Next
	} else {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	return _pokeMap(cfg, url)
}

func commandMapb(cfg *Config, _ []string) error {
	var url string
	if cfg.Previous != nil && *cfg.Previous != "" {
		url = *cfg.Previous
	} else {
		fmt.Println("you're on the first page")
		return nil
	}
	return _pokeMap(cfg, url)
}

func _pokeMap(cfg *Config, url string) error {
	data, ok := cfg.cache.Get(url)
	var locationsResp pokeapi.RespShallowLocations
	var err error
	if ok {
		err = json.Unmarshal(data, &locationsResp)
		if err != nil {
			return err
		}
	} else {
		locationsResp, err = cfg.client.ListLocations(&url)
		if err != nil {
			return err
		}
		var jsonData []byte
		jsonData, err = json.Marshal(locationsResp)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, jsonData)
	}
	cfg.Next = locationsResp.Next
	cfg.Previous = locationsResp.Previous

	for _, result := range locationsResp.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExplore(cfg *Config, locationArea []string) error {
	var url string = "https://pokeapi.co/api/v2/location-area/" + locationArea[0]
	var respDeepLocations pokeapi.RespDeepLocations
	var err error
	data, ok := cfg.cache.Get(url)
	if ok {
		err = json.Unmarshal(data, &respDeepLocations)
		if err != nil {
			return err
		}
	} else {
		respDeepLocations, err = cfg.client.LocationExplore(&url)
		if err != nil {
			return err
		}
		jsonData, err := json.Marshal(respDeepLocations)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, jsonData)
	}

	fmt.Printf("Exploring %s...\n", locationArea[0])
	fmt.Println("Found Pokemon:")
	for _, result := range respDeepLocations.PokemonEncounters {
		fmt.Printf(" - %s\n", result.Pokemon.Name)
	}
	return nil
}

func commandExit(cfg *Config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // unreachable, but keeps the IDE happy
}

func commandHelp(cfg *Config, _ []string) error {
	//	fmt.Println("Usage: Pokedex [command]")
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("\nUsage:")
	fmt.Println()
	fmt.Println("help: Displays a help message")
	fmt.Println("map: Display the next page of locations")
	fmt.Println("mapb: Display the previous page of locations")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println()
	return nil
}

func commandInspect(cfg *Config, args []string) error {
	pokemonName := args[0]
	data, ok := cfg.pokemon[pokemonName]
	if ok {
		fmt.Printf("Name: %s\n", data.Name)
		fmt.Printf("Height: %d\n", data.Height)
		fmt.Printf("Weight: %d\n", data.Weight)
		fmt.Println("Stats:")
		for _, stat := range data.Stats {
			fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, pType := range data.Types {
			fmt.Printf("  - %s\n", pType.Type.Name)
		}
		return nil
	}

	fmt.Println("you have not caught that pokemon")
	return nil
}

func commandPokedex(cfg *Config, args []string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.pokemon {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}
