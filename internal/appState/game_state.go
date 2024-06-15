package appState

import (
	"sync"
)

// Constants representing the state of a cell in the game board
const (
	Ship  = "Ship"
	Empty = ""
	Hit   = "Hit"
	Miss  = "Miss"
	Sunk  = "Sunk"
)

// InitializeNewGameState initializes a new game state
func InitializeNewGameState() *GameState {
	return &GameState{
		m:             sync.Mutex{},
		player:        &Player{},
		opponent:      &Player{},
		playerBoard:   NewBoard(),
		opponentBoard: NewBoard(),
		oppShipsSun: map[int]int{
			1: 4,
			2: 3,
			3: 2,
			4: 1,
		},
	}
}

// GetGameState returns the current game state
func (g *GameState) GetGameState() *GameState {
	g.m.Lock()
	defer g.m.Unlock()
	return g
}

// UpdateGameState updates the game state
func (g *GameState) UpdateGameState(nick, desc, opp, oppdesc string) {
	g.m.Lock()
	defer g.m.Unlock()
	g.player.Nick = nick
	g.player.Description = desc
	g.opponent.Nick = opp
	g.opponent.Description = oppdesc
}

// UpdatePlayerBoard updates the player's board
func (g *GameState) UpdatePlayerBoard(playerState [10][10]string) ([10][10]string, error) {
	g.m.Lock()
	defer g.m.Unlock()
	g.playerBoard.updatePlayerStates(playerState)
	return g.playerBoard.PlayerState, nil
}

// UpdateOpponentBoard updates the opponent's board
func (g *GameState) UpdateOpponentBoard(opponentState [10][10]string) ([10][10]string, error) {
	g.m.Lock()
	defer g.m.Unlock()
	g.opponentBoard.updatePlayerStates(opponentState)
	return g.opponentBoard.PlayerState, nil
}

// GetPlayerBoard returns the player's board
func (g *GameState) GetPlayerBoard() [10][10]string {
	g.m.Lock()
	defer g.m.Unlock()
	return g.playerBoard.PlayerState
}

// MarkPlayerBoard marks a cell on the player's board
func (g *GameState) MarkPlayerBoard(x, y int) {
	g.m.Lock()
	defer g.m.Unlock()
	switch g.playerBoard.PlayerState[x][y] {
	case Ship:
		g.playerBoard.PlayerState[x][y] = Hit
	case Empty:
		g.playerBoard.PlayerState[x][y] = Miss
	}
}

// GetOpponentBoard returns the opponent's board
func (g *GameState) GetOpponentBoard() [10][10]string {
	g.m.Lock()
	defer g.m.Unlock()
	return g.opponentBoard.PlayerState
}

// MarkOpponentBoard marks a cell on the opponent's board
func (g *GameState) MarkOpponentBoard(x int, y int, result string) int {
	g.m.Lock()
	defer g.m.Unlock()
	if result == Sunk {
		g.opponentBoard.PlayerState[x][y] = result
		_, l := g.opponentBoard.DrawBorder(x, y)
		g.oppShipsSun[l]--
		return l
	}
	g.opponentBoard.PlayerState[x][y] = result

	return 0
}

// CheckIfAlreadyHit checks if a cell has already been hit
func (g *GameState) CheckIfAlreadyHit(x, y int) bool {
	g.m.Lock()
	defer g.m.Unlock()
	s := g.opponentBoard.PlayerState[x][y]
	return s == Hit || s == Miss
}

// IncrementHitCount increases the hit count
func (g *GameState) IncrementHitCount(str string) {
	g.m.Lock()
	defer g.m.Unlock()
	if str == "hit" || str == "sunk" {
		g.hits++
	}
	g.totalShots++
}

// GetTotalShots returns the total number of shots
func (g *GameState) GetTotalShots() int {
	g.m.Lock()
	defer g.m.Unlock()
	return g.totalShots
}

// GetTotalHits returns the total number of hits
func (g *GameState) GetTotalHits() int {
	g.m.Lock()
	defer g.m.Unlock()
	return g.hits
}

// ModifyPlayerInformation modifies player information
func (g *GameState) ModifyPlayerInformation(name string, description string) {
	g.m.Lock()
	defer g.m.Unlock()
	g.player.Nick = name
	g.player.Description = description
}

// GetPlayerInfo returns player information
func (g *GameState) GetPlayerInfo() (string, string) {
	g.m.Lock()
	defer g.m.Unlock()
	return g.player.Nick, g.player.Description
}

// GetOppDesc returns the opponent's description
func (g *GameState) GetOppDesc() string {
	g.m.Lock()
	defer g.m.Unlock()
	return g.opponent.Description
}

// GetPlayerDesc returns the player's description
func (g *GameState) GetPlayerDesc() string {
	g.m.Lock()
	defer g.m.Unlock()
	return g.player.Description
}

// UpdatePlayersDesc updates players' descriptions
func (g *GameState) UpdatePlayersDesc(desc, oppDesc string) {
	g.m.Lock()
	defer g.m.Unlock()
	g.player.Description = desc
	g.opponent.Description = oppDesc
}

// AddShip adds a ship to the board
func (g *GameState) AddShip(x int, y int) {
	g.m.Lock()
	defer g.m.Unlock()
	g.playerBoard.PlayerState[x][y] = Ship
}

// ClearState resets the game state
func (g *GameState) ClearState() {
	g.m.Lock()
	defer g.m.Unlock()
	g.playerBoard = NewBoard()
	g.opponentBoard = NewBoard()
	g.totalShots = 0
	g.hits = 0
}

// RetrieveOpponentSunkShipsCount returns the number of sunk ships of the opponent
func (g *GameState) RetrieveOpponentSunkShipsCount() map[int]int {
	g.m.Lock()
	defer g.m.Unlock()
	return g.oppShipsSun
}

// UpdateLastGameStatus updates the status of the last game
func (g *GameState) UpdateLastGameStatus(status string) {
	g.m.Lock()
	defer g.m.Unlock()
	g.lastGameStatus = status
}

// LastGameStatus returns the status of the last game
func (g *GameState) LastGameStatus() string {
	g.m.Lock()
	defer g.m.Unlock()
	return g.lastGameStatus
}
