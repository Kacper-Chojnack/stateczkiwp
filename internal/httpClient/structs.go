package httpClient

import (
	"battleships/internal/appState"
	"fmt"
	"net/http"
	"sort"
)

// GameStatus represents the game state
type GameStatus struct {
	GameStatus     string   `json:"game_status"`
	LastGameStatus string   `json:"last_game_status"`
	Nick           string   `json:"nick"`
	OppShots       []string `json:"opp_shots"`
	Opponent       string   `json:"opponent"`
	ShouldFire     bool     `json:"should_fire"`
	Timer          int      `json:"timer"`
}

// StartGameData represents the data to start a game
type StartGameData struct {
	Coords     []string `json:"coords"`
	Desc       string   `json:"desc"`
	Nick       string   `json:"nick"`
	TargetNick string   `json:"target_nick"`
	WPBot      bool     `json:"wpbot"`
}

// GameBoard represents the game board
type GameBoard struct {
	Board []string `json:"board"`
}

// FireData represents the shot data
type FireData struct {
	Coord string `json:"coord"`
}

// FireResult represents the shot result
type FireResult struct {
	Result string `json:"result"`
}

// GameDescription represents the game description
type GameDescription struct {
	Desc     string `json:"desc"`
	Nick     string `json:"nick"`
	OppDesc  string `json:"opp_desc"`
	Opponent string `json:"opponent"`
}

// GameList represents the list of games
type GameList []struct {
	Guest  string `json:"guest"`
	Host   string `json:"host"`
	ID     string `json:"id"`
	Status string `json:"status"`
}

// LobbyPlayer represents a player in the lobby
type LobbyPlayer struct {
	GameStatus string `json:"game_status"`
	Nick       string `json:"nick"`
}

// PlayerStats represents the player's statistics
type PlayerStats struct {
	Games  int    `json:"games"`
	Nick   string `json:"nick"`
	Points int    `json:"points"`
	Rank   int    `json:"rank"`
	Wins   int    `json:"wins"`
}

// String returns the formatted player's statistics
func (ps PlayerStats) String() string {
	return fmt.Sprintf("Nick: %s, Number of games: %d, Number of points: %d, Ranking position: %d, Number of wins: %d",
		ps.Nick, ps.Games, ps.Points, ps.Rank, ps.Wins)
}

// TopPlayerStats represents the statistics of the top 10 players
type TopPlayerStats struct {
	Stats []PlayerStats `json:"stats"`
}

// String returns the formatted top players' statistics
func (tps TopPlayerStats) String() string {
	var result string
	for i, ps := range tps.Stats {
		result += fmt.Sprintf("%d.%s\n", i+1, ps)
	}
	return result
}

// Shot represents a shot
type Shot struct {
	Coord string `json:"coord"`
}

// GameStat represents the game statistics
type GameStat struct {
	Games  int    `json:"games"`
	Nick   string `json:"nick"`
	Points int    `json:"points"`
	Rank   int    `json:"rank"`
	Wins   int    `json:"wins"`
}

// GameStats represents the list of game statistics
type GameStats []GameStat

// Len returns the length of the GameStats list
func (g GameStats) Len() int {
	return len(g)
}

// Less compares two elements of the GameStats list
func (g GameStats) Less(i, j int) bool {
	if g[i].Wins != g[j].Wins {
		return g[i].Wins > g[j].Wins
	}
	return g[i].Nick < g[j].Nick
}

// Swap swaps two elements of the GameStats list
func (g GameStats) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

// String returns the formatted game statistics
func (g GameStats) String() string {
	// Sorting game statistics before printing
	sort.Sort(g)

	// Formatting game statistics as a string
	var result string
	for _, stat := range g {
		result += fmt.Sprintf("Nick: %s\n", stat.Nick)
		result += fmt.Sprintf("Number of games: %d\n", stat.Games)
		result += fmt.Sprintf("Number of points: %d\n", stat.Points)
		result += fmt.Sprintf("Ranking position: %d\n", stat.Rank)
		result += fmt.Sprintf("Number of wins: %d\n", stat.Wins)
		result += "\n"
	}

	return result
}

// Game represents a game
type Game struct {
	Client *Client
	state  *appState.GameState
}

// GameState represents the game state
type GameState struct {
	PlayerBoard  [10][10]string `json:"player_board"`
	OppBoard     [10][10]string `json:"opp_board"`
	TotalShots   int            `json:"total_shots"`
	TotalHits    int            `json:"total_hits"`
	PlayerDesc   string         `json:"player_desc"`
	OppDesc      string         `json:"opp_desc"`
	OppShipsSunk map[int]int
}

// Structures from the api_client.go file

// Client represents an API client
type Client struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

// ErrorMessage represents an error message
type ErrorMessage struct {
	Message string `json:"message"`
}

// ApiError represents an API error
type ApiError struct {
	ErrorMessage
	ErrorType string
}

// Error returns the formatted API error
func (e ApiError) Error() string {
	return e.ErrorType + ": " + e.Message
}

// UnauthorizedError represents an authorization error
type UnauthorizedError struct {
	ApiError
}

// Error returns the formatted authorization error
func (e UnauthorizedError) Error() string {
	return e.ApiError.Error()
}

// ForbiddenError represents a forbidden access error
type ForbiddenError struct {
	ApiError
}

// Error returns the formatted forbidden access error
func (e ForbiddenError) Error() string {
	return e.ApiError.Error()
}

// RateLimitExceededError represents a rate limit exceeded error
type RateLimitExceededError struct {
	ApiError
}

// Error returns the formatted rate limit exceeded error
func (e RateLimitExceededError) Error() string {
	return e.ApiError.Error()
}

// BadRequestError represents a bad request error
type BadRequestError struct {
	ApiError
}

// Error returns the formatted bad request error
func (e BadRequestError) Error() string {
	return e.ApiError.Error()
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	ApiError
}

// Error returns the formatted resource not found error
func (e NotFoundError) Error() string {
	return e.ApiError.Error()
}
