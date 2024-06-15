package httpClient

import (
	"battleships/internal/appState"
	"fmt"
	"net/http"
	"strconv"
)

// setStatesFromCoords converts []string to [][]string
func setStatesFromCoords(coords []string, s string) [10][10]string {
	state := [10][10]string{}
	for i := range state {
		state[i] = [10]string{}
	}
	for _, coord := range coords {
		x, y := mapToState(coord)
		state[x][y] = s
	}
	return state
}

// mapToState converts string to int
func mapToState(coord string) (int, int) {
	if len(coord) > 2 {
		return int(coord[0] - 65), 9
	}
	x := int(coord[0] - 65)
	y := int(coord[1] - 49)
	return x, y
}

// NewGame returns a new game instance
func NewGame() *Game {
	return &Game{
		Client: NewClient("https://go-pjatk-server.fly.dev/api", ""),
		state:  appState.InitializeNewGameState(),
	}
}

// FireShot fires a shot at the given coordinate
func (g *Game) FireShot(coord string) (FireResult, int, error) {
	result, err := g.Client.Fire(FireData{Coord: coord})
	if err != nil {
		return FireResult{}, 0, err
	}
	l := g.MarkOpponent(coord, result)
	return result, l, err
}

// StartGame starts the game
func (g *Game) StartGame(nick, desc, targetNick string, coords []string, botGame bool) {
	_, err := g.Client.StartGame(nick, desc, targetNick, coords, botGame)
	if err != nil {
		return
	}
}

// GetGameStatus returns the current game status
func (g *Game) GetGameStatus() (GameStatus, error) {
	gameState, err := g.Client.GetGameStatus()
	if err != nil {
		return GameStatus{}, err
	}
	return gameState, nil
}

// SetPlayerBoard sets the player's board
func (g *Game) SetPlayerBoard(coords []string) ([10][10]string, error) {
	board, err := g.state.UpdatePlayerBoard(setStatesFromCoords(coords, appState.Ship))
	if err != nil {
		return [10][10]string{}, err
	}
	return board, nil
}

// GetDescription returns the game description
func (g *Game) GetDescription() (GameDescription, error) {
	return g.Client.GetGameDescription()
}

// LoadPlayerBoard loads the player's board
func (g *Game) LoadPlayerBoard() (*GameBoard, error) {
	return g.Client.GetGameBoard()
}

// UpdateGameState updates the game state
func (g *Game) UpdateGameState(nick string, desc string, opponent string, oppDesc string) {
	g.state.UpdateGameState(nick, desc, opponent, oppDesc)
}

// GetPlayerBoard returns the player's board
func (g *Game) GetPlayerBoard() [10][10]string {
	return g.state.GetPlayerBoard()
}

// MarkOpponentShots marks the opponent's shots on the player's board
func (g *Game) MarkOpponentShots(shots []string) {
	for _, coord := range shots {
		x, y := mapToState(coord)
		g.state.MarkPlayerBoard(x, y)
	}
}

// GetGameState returns the current game state
func (g *Game) GetGameState() (*appState.GameState, error) {
	return g.state.GetGameState(), nil
}

// GetOpponentBoard returns the opponent's board
func (g *Game) GetOpponentBoard() [10][10]string {
	return g.state.GetOpponentBoard()
}

// MarkOpponent marks the shot result on the opponent's board
func (g *Game) MarkOpponent(shot string, result FireResult) int {
	if shot == "" {
		return 0
	}
	x, y := mapToState(shot)
	var mark string
	switch result.Result {
	case "sunk":
		mark = appState.Sunk
	case "hit":
		mark = appState.Hit
	case "miss":
		mark = appState.Miss
	}
	g.state.IncrementHitCount(result.Result)
	return g.state.MarkOpponentBoard(x, y, mark)
}

// UpdatePlayerInfo updates player information
func (g *Game) UpdatePlayerInfo(name string, description string) {
	g.state.ModifyPlayerInformation(name, description)
}

// GetPlayerInfo returns player information
func (g *Game) GetPlayerInfo() (string, string) {
	return g.state.GetPlayerInfo()
}

// UpdatePlayersDesc updates players' descriptions
func (g *Game) UpdatePlayersDesc(d GameDescription) {
	g.state.UpdatePlayersDesc(d.Desc, d.OppDesc)
}

// GetTopPlayerStats returns the top players' statistics
func (g *Game) GetTopPlayerStats() (TopPlayerStats, error) {
	stats, err := g.Client.GetTopPlayerStats()
	if err != nil {
		return TopPlayerStats{}, err
	}
	return stats, nil
}

// MarkPlayerShip marks the player's ship on the board
func (g *Game) MarkPlayerShip(coords string) {
	x, y := mapToState(coords)
	g.state.AddShip(x, y)
}

// GetPlayerCoords returns the player's ship coordinates
func (g *Game) GetPlayerCoords() []string {
	states := g.state.GetPlayerBoard()

	var coords []string
	for i, row := range states {
		for j, s := range row {
			if s == appState.Ship {
				coords = append(coords, mapFromState(i, j))
			}
		}
	}

	// Ensures that the coordinates are within the range from A1 to J10
	var fixedCoords []string
	for _, coord := range coords {
		if isValidCoord(coord) {
			fixedCoords = append(fixedCoords, coord)
		}
	}

	return fixedCoords
}

// GetPlayerStats returns the player's statistics
func (g *Game) GetPlayerStats(name string) GameStats {
	stats, err := g.Client.GetPlayerStats(name)
	if err != nil {
		_ = fmt.Errorf("error while fetching player's statistics: %v", err)
		return GameStats{}
	}
	return stats
}

// GetPlayerLobby returns the list of players in the lobby
func (c *Client) GetPlayerLobby() ([]LobbyPlayer, error) {
	req, err := c.getRequest("/lobby")
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var players []LobbyPlayer
	if err := handleResponse(resp, http.StatusOK, &players); err != nil {
		return nil, err
	}
	return players, nil
}

// ClearState clears the game state
func (g *Game) ClearState() {
	g.state.ClearState()
}

// UpdateLastGameStatus updates the status of the last game
func (g *Game) UpdateLastGameStatus(status string) {
	g.state.UpdateLastGameStatus(status)
}

// LastGameStatus returns the status of the last game
func (g *Game) LastGameStatus() string {
	return g.state.LastGameStatus()
}

// AbortGame aborts the game
func (g *Game) AbortGame() {
	err := g.Client.AbortGame()
	if err != nil {
		return
	}
}

// mapFromState converts int to string
func mapFromState(x, y int) string {
	return string(byte(x+65)) + strconv.Itoa(y+1)
}

// isValidCoord checks if the coordinates are valid
func isValidCoord(coord string) bool {
	if len(coord) < 2 || len(coord) > 3 {
		return false
	}

	column := coord[0]
	row := coord[1:]
	if column < 'A' || column > 'J' {
		return false
	}

	num, err := strconv.Atoi(row)
	if err != nil || num < 1 || num > 10 {
		return false
	}

	return true
}
