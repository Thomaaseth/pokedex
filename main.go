package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	pokecache "github.com/Thomaaseth/pokedexcli/internal"
)

func cleanInput(text string) []string {
	trimed := strings.TrimSpace(text)
	lowercaseText := strings.ToLower(trimed)
	words := strings.Fields(lowercaseText)
	return words
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for key, value := range commands {
		fmt.Println(key + ": " + value.description)
	}
	return nil
}

func commandMap(config *Config) error {
	url := ""
	if config.Next != nil {
		url = *config.Next
	}

	locs, err := getLocations(url, config.Cache)
	if err != nil {
		return err
	}
	config.Next = locs.Next
	config.Previous = locs.Previous

	for _, loc := range locs.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapb(config *Config) error {
	if config.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}
	url := *config.Previous
	locs, err := getLocations(url, config.Cache)
	if err != nil {
		return err
	}
	config.Next = locs.Next
	config.Previous = locs.Previous

	for _, loc := range locs.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

var commands map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type Config struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
}

type locations struct {
	Results  []locationArea `json:"results"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
}

type locationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func getLocations(url string, cache *pokecache.Cache) (locations, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	if cachedData, found := cache.Get(url); found {
		fmt.Println("Using cached data for:", url)

		var locs locations
		err := json.Unmarshal(cachedData, &locs)
		if err != nil {
			fmt.Println("Error unmarshalling cached content")
			return locations{}, err
		}
		return locs, nil
	}

	// If not in cache, make the HTTP request
	fmt.Println("Fetching data from API:", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting list of locations")
		return locations{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading resp body")
		return locations{}, err
	}

	cache.Add(url, body)

	var locs locations
	err = json.Unmarshal(body, &locs)
	if err != nil {
		fmt.Println("Error unmarshalling content")
		return locations{}, err
	}
	return locs, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(5 * time.Minute)
	config := &Config{
		Cache: cache,
	}

	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    func() error { return commandExit() },
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    func() error { return commandHelp() },
		},
		"map": {
			name:        "map",
			description: "Show the next 20 Pokemon locations",
			callback:    func() error { return commandMap(config) },
		},
		"mapb": {
			name:        "mapb",
			description: "Show the previous 20 Pokemon locations",
			callback:    func() error { return commandMapb(config) },
		},
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()

		words := cleanInput(userInput)
		if len(words) > 0 {
			firstWord := words[0]

			cmd, ok := commands[firstWord]
			if ok {
				cmd.callback()
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}
