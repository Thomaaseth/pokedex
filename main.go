package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	trimed := strings.TrimSpace(text)
	lowercaseText := strings.ToLower(trimed)
	words := strings.Fields(lowercaseText)
	return words
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()

		words := cleanInput(userInput)
		if len(words) > 0 {
			firstWord := words[0]
			fmt.Println("Your command was:", firstWord)
		}
	}
}
