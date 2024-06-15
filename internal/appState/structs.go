package appState

import (
	"sync"
)

// Player represents a player in the game
type Player struct {
	Nick        string
	Description string
}

// Board represents the game board
type Board struct {
	PlayerState   [10][10]string
	OpponentState [10][10]string
}

// NewBoard creates a new game board
func NewBoard() *Board {
	return &Board{
		PlayerState:   [10][10]string{},
		OpponentState: [10][10]string{},
	}
}

// updatePlayerStates updates the player's state on the board
func (b *Board) updatePlayerStates(playerState [10][10]string) {
	b.PlayerState = playerState
}

// Mark marks a cell on the board
func (b *Board) Mark(x int, y int, mark string) {
	b.PlayerState[x][y] = mark
}

// GameState represents the game state
type GameState struct {
	player         *Player
	opponent       *Player
	playerBoard    *Board
	opponentBoard  *Board
	totalShots     int
	hits           int
	m              sync.Mutex
	lastGameStatus string
	oppShipsSun    map[int]int
}
