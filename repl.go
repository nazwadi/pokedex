package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

func startRepl(cfg *Config) {
	m := map[string]cliCommand{
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
			description: "Map a Pokedex using the given arguments",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Map backwards a Pokedex using the given arguments",
			callback:    commandMapb,
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
		value, ok := m[cleanedInput[0]]
		if ok {
			err := value.callback(cfg)
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

func commandMap(cfg *Config) error {
	var url string
	if cfg.Next != nil && *cfg.Next != "" {
		url = *cfg.Next
	} else {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	return _pokeMap(cfg, url)
}

func commandMapb(cfg *Config) error {
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
	locationsResp, err := cfg.client.ListLocations(&url)
	if err != nil {
		return err
	}
	cfg.Next = locationsResp.Next
	cfg.Previous = locationsResp.Previous

	for _, result := range locationsResp.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // unreachable, but keeps the IDE happy
}

func commandHelp(cfg *Config) error {
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
