package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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

func commandExit(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(args []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for key, value := range commands {
		fmt.Println(key + ": " + value.description)
	}
	return nil
}

func commandMap(config *Config, args []string) error {
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

func commandMapb(config *Config, args []string) error {
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

func exploreMap(config *Config, args []string) error {
	if len(args) == 0 {
		fmt.Println("you must provide a location area name")
		return nil
	}
	locationAreaName := args[0]

	fmt.Printf("Exploring %s...\n", locationAreaName)

	details, err := getLocationDetails(locationAreaName, config.Cache)
	if err != nil {
		return fmt.Errorf("error exploring location area: %s", err)
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range details.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func catchPokemon(config *Config, args []string) error {
	if len(args) == 0 {
		fmt.Println("you must provide a pokemon name")
		return nil
	}
	pokemonName := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := getPokemonDetails(pokemonName, config.Cache)
	if err != nil {
		return fmt.Errorf("error retrieving pokemon details")
	}
	catchChance := 100 - (pokemon.BaseExperience / 4)
	if catchChance < 0 {
		catchChance = 0
	}
	if catchChance > 90 {
		catchChance = 90
	}
	randomNum := r.Intn(101)
	if randomNum <= catchChance {
		fmt.Printf("%s was caught!\n", pokemonName)
		config.CaughtPokemon[pokemonName] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	return nil
}

func inspectPokemon(config *Config, args []string) error {
	if len(args) == 0 {
		fmt.Println("you must provide a pokemon name")
		return nil
	}
	pokemonName := args[0]

	pokemon, found := config.CaughtPokemon[pokemonName]
	if !found {
		fmt.Println("You haven't caught this pokemon yet!")
		return nil
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats: \n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}

func pokedex(config *Config, args []string) error {
	if len(config.CaughtPokemon) == 0 {
		fmt.Println("You haven't caught any pokemons yet!")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for name := range config.CaughtPokemon {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

var commands map[string]cliCommand
var r *rand.Rand

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

type Config struct {
	Next          *string
	Previous      *string
	Cache         *pokecache.Cache
	CaughtPokemon map[string]Pokemon
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

type PokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
	} `json:"pokemon"`
}

type Pokemon struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
	BaseExperience int           `json:"base_experience"`
}

type PokemonStat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type PokemonType struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type LocationAreaDetails struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
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

func getLocationDetails(locationName string, cache *pokecache.Cache) (LocationAreaDetails, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", locationName)

	// Check cache first
	if cachedData, found := cache.Get(url); found {
		fmt.Println("Using cached data for:", url)

		var details LocationAreaDetails
		err := json.Unmarshal(cachedData, &details)
		if err != nil {
			return LocationAreaDetails{}, fmt.Errorf("error unmashalling cached content: %w", err)
		}
		return details, nil
	}

	// If not in cached, make HTTP request
	fmt.Println("Fetching data from API:", url)
	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaDetails{}, fmt.Errorf("error getting location details: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaDetails{}, fmt.Errorf("error reading response: %w", err)
	}

	cache.Add(url, body)

	var details LocationAreaDetails
	err = json.Unmarshal(body, &details)
	if err != nil {
		return LocationAreaDetails{}, fmt.Errorf("error unmarshalling content: %w", err)
	}
	return details, nil
}

func getPokemonDetails(pokemonName string, cache *pokecache.Cache) (Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)

	if cachedData, found := cache.Get(url); found {
		fmt.Println("Using cached data for:", url)

		var pokemon Pokemon
		err := json.Unmarshal(cachedData, &pokemon)
		if err != nil {
			return Pokemon{}, fmt.Errorf("error unmashalling cached content: %w", err)
		}
		return pokemon, nil
	}
	fmt.Println("Fetching data from API:", url)
	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error getting pokemon details: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error reading response: %w", err)
	}

	cache.Add(url, body)

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error unmarshalling content: %w", err)
	}
	return pokemon, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(5 * time.Minute)
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	config := &Config{
		Cache:         cache,
		CaughtPokemon: make(map[string]Pokemon),
	}

	commands = map[string]cliCommand{
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
			description: "Show the next 20 Pokemon locations",
			callback:    func(args []string) error { return commandMap(config, args) },
		},
		"mapb": {
			name:        "mapb",
			description: "Show the previous 20 Pokemon locations",
			callback:    func(args []string) error { return commandMapb(config, args) },
		},
		"explore": {
			name:        "explore <location name>",
			description: "Show a list of all Pokemons located here",
			callback:    func(args []string) error { return exploreMap(config, args) },
		},
		"catch": {
			name:        "catch <pokemon name>",
			description: "Try to catch a pokemon",
			callback:    func(args []string) error { return catchPokemon(config, args) },
		},
		"inspect": {
			name:        "inspect <pokemon name>",
			description: "Get details of a pokemon",
			callback:    func(args []string) error { return inspectPokemon(config, args) },
		},
		"pokedex": {
			name:        "pokedex",
			description: "retrieve all your caught pokemons",
			callback:    func(args []string) error { return pokedex(config, args) },
		},
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()

		words := cleanInput(userInput)
		if len(words) > 0 {
			firstWord := words[0]
			args := words[1:]
			cmd, ok := commands[firstWord]
			if ok {
				cmd.callback(args)
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}
