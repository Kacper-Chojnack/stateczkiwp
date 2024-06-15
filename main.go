package main

import (
	"battleships/internal/game"
	"battleships/internal/httpClient"
	"context"
)

// main is the entry point of the application
func main() {
	// Create a new context for managing the lifecycle of goroutines
	ctx := context.Background()

	// Create communication channels for game status, player shots, and game state
	gameStatusChannel, playerShotsChannel, gameStateChannel := createChannels()

	// Initialize a new game application with the created channels
	app := game.NewApp(gameStatusChannel, playerShotsChannel, gameStateChannel)

	// Start the game menu
	app.Menu(ctx)
}

// createChannels is a helper function that creates and returns channels for game status, player shots, and game state
func createChannels() (chan httpClient.GameStatus, chan string, chan httpClient.GameState) {
	// Create a channel for game status updates
	gameStatusChannel := make(chan httpClient.GameStatus)

	// Create a channel for player shots
	playerShotsChannel := make(chan string)

	// Create a channel for game state updates
	gameStateChannel := make(chan httpClient.GameState)

	// Return the created channels
	return gameStatusChannel, playerShotsChannel, gameStateChannel
}
