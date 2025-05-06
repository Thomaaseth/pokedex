# Pokedex CLI - Pokemon Explorer

A command-line interface Pokedex application written in Go that allows users to explore the Pokemon world, catch Pokemon, and build their collection.

## Features

- **Map Navigation**: Browse through different Pokemon locations using pagination
- **Location Exploration**: Discover which Pokemon can be found at each location
- **Pokemon Catching**: Try to catch Pokemon with a chance-based system
- **Pokemon Inspection**: View detailed information about caught Pokemon
- **Pokedex Management**: Keep track of all your caught Pokemon
- **Caching System**: Efficient data retrieval with automatic cache expiration

## Installation

### Prerequisites

- Go 1.16 or higher
- Git

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/Thomaaseth/pokedex.git
   cd pokedex

2. Build the application:
    ```bash
    go build

3. Run the application:
    ```bash
    go run main.go


### Usage
Once the Pokedex CLI is running, you'll see a prompt:
Pokedex >

You can use the following commands:
- help Displays a help message with all available commands
- exit Exit the Pokedex application
- map Show the next 20 Pokemon locations
- mapb Show the previous 20 Pokemon locations
- explore <location name> Show a list of all Pokemon found at the specified location
- catch <pokemon name> Try to catch a Pokemon
- inspect <pokemon name> Get details of a caught Pokemon
- pokedex View all your caught Pokemon

#### Examples

**Exploring a location**:
Pokedex > explore canalave-city-area
Exploring canalave-city-area...
Found Pokemon:
 - staravia
 - starly
 - chimchar
 - turtwig
 - piplup

**Catching a Pokemon**:
Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!

**Inspecting a Pokemon**:
Pokedex > inspect pikachu
Name: pikachu
Height: 4
Weight: 60
Stats: 
  -hp: 35
  -attack: 55
  -defense: 40
  -special-attack: 50
  -special-defense: 50
  -speed: 90
Types:
  - electric

### Architecture

#### Major Components

CLI Interface: 
Handles user input/output and command processing
API Client: Communicates with the PokeAPI to fetch Pokemon data
Cache System: Stores API responses to minimize network requests
Pokemon Storage: Manages the user's caught Pokemon collection

Data Flow:
User enters a command
Command is processed and relevant function is called
If data is needed, the system first checks the cache
If not in cache, an API request is made and the result is cached
Data is displayed to the user

Dependencies:
Standard Go libraries: bufio, encoding/json, fmt, io, math/rand, net/http, os, strings, time
Custom cache package: package github.com/{your-username}/pokedex/internal/pokecache

API:
This project utilizes the PokeAPI, a free and open RESTful Pokemon API.

### License
MIT License

### Contributing

Fork the repository:
Create your feature branch (git checkout -b feature/amazing-feature)
Commit your changes (git commit -m 'Add some amazing feature')
Push to the branch (git push origin feature/amazing-feature)
Open a Pull Request

### Acknowledgments

PokeAPI for providing the Pokemon data
Built with Boot.dev as guided project