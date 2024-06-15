package game

import (
	"battleships/internal/appState"
	"battleships/internal/httpClient"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"strconv"
	"sync"
)

// NewGui creates a new user interface
func NewGui() *Gui {
	return &Gui{
		gui:            gui.NewGUI(false),
		playerNick:     gui.NewText(1, 27, "Player", nil),
		playerDesc:     gui.NewText(1, 28, "Your board", nil),
		opponentNick:   gui.NewText(50, 27, "Opponent", nil),
		opponentDesc:   gui.NewText(50, 28, "Opponent's board", nil),
		playerBoard:    gui.NewBoard(1, 5, nil),
		opponentBoard:  gui.NewBoard(50, 5, nil),
		waiting:        gui.NewText(10, 10, "Waiting for opponent...", nil),
		turn:           gui.NewText(1, 3, "", nil),
		timer:          gui.NewText(1, 1, "", nil),
		mu:             sync.Mutex{},
		numberOf1Ships: gui.NewText(100, 9, "4 ships of length 1", nil),
		numberOf2Ships: gui.NewText(100, 10, "3 ships of length 2", nil),
		numberOf3Ships: gui.NewText(100, 11, "2 ships of length 3", nil),
		numberOf4Ships: gui.NewText(100, 12, "1 ship of length 4", nil),
	}
}

// SetPlayerBoard sets the player's board
func (g *Gui) SetPlayerBoard(states [10][10]gui.State) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.playerBoard.SetStates(states)
}

// sendPlayerShots sends player shots
func (g *Gui) sendPlayerShots(ctx context.Context, shotsChannel chan string) {
	coord := g.playerBoard.Listen(ctx)
	shotsChannel <- coord
}

// displayBoard displays the board
func (g *Gui) displayBoard() {
	g.gui.Draw(g.playerBoard)
	g.gui.Draw(g.opponentBoard)
}

// handleGameStatus handles the game status
func (g *Gui) handleGameStatus(ctx context.Context, events chan httpClient.GameStatus) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case status := <-events:
			g.mu.Lock()
			g.updateTurn(status)
			g.updateTimer(status)
			g.updatePlayers(status)
			g.gui.Draw(g.waiting)
			g.gui.Draw(g.numberOf1Ships)
			g.gui.Draw(g.numberOf2Ships)
			g.gui.Draw(g.numberOf3Ships)
			g.gui.Draw(g.numberOf4Ships)
			if status.GameStatus == "ended" {
				g.gui.Draw(gui.NewText(5, 10, "Game ended. Press ctrl + c to return to the menu", nil))
			}
			if status.GameStatus == "game_in_progress" {
				g.waiting.SetText("")
				g.gui.Draw(g.waiting)
			}
		}
		g.mu.Unlock()
	}
}

// updatePlayers updates players
func (g *Gui) updatePlayers(status httpClient.GameStatus) {
	g.playerNick.SetText(status.Nick)
	g.gui.Draw(g.playerNick)
	g.opponentNick.SetText(status.Opponent)
	g.gui.Draw(g.opponentNick)
}

// updateTimer updates the timer
func (g *Gui) updateTimer(status httpClient.GameStatus) {
	g.timer.SetText(fmt.Sprintf("Time: %d", status.Timer))
	g.gui.Draw(g.timer)
}

// updateTurn updates the turn
func (g *Gui) updateTurn(status httpClient.GameStatus) {
	if status.ShouldFire {
		g.turn.SetText("Your turn")
		g.gui.Draw(g.turn)
	} else {
		g.turn.SetText("Opponent's turn")
		g.gui.Draw(g.turn)
	}
}

// handleGameState handles the game state
func (g *Gui) handleGameState(ctx context.Context, state chan httpClient.GameState) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case gameState := <-state:
			g.mu.Lock()
			g.playerBoard.SetStates(mapStatesToGuiMarks(gameState.PlayerBoard))
			g.opponentBoard.SetStates(mapStatesToGuiMarks(gameState.OppBoard))
			g.opponentBoard.SetStates(mapStatesToGuiMarks(gameState.OppBoard))
			g.gui.Draw(gui.NewText(1, 28, gameState.PlayerDesc, nil))
			g.gui.Draw(gui.NewText(58, 28, gameState.OppDesc, nil))
			g.gui.Draw(gui.NewText(1, 2, fmt.Sprintf("Accuracy: %s %%",
				getAccuracy(gameState.TotalHits, gameState.TotalShots)), nil))
			g.mu.Unlock()

			g.numberOf1Ships.SetText(strconv.Itoa(gameState.OppShipsSunk[1]) + " ships of length 1")
			g.numberOf2Ships.SetText(strconv.Itoa(gameState.OppShipsSunk[2]) + " ships of length 2")
			g.numberOf3Ships.SetText(strconv.Itoa(gameState.OppShipsSunk[3]) + " ships of length 3")
			g.numberOf4Ships.SetText(strconv.Itoa(gameState.OppShipsSunk[4]) + " ships of length 4")
			g.drawLegend()
		}
	}
}

// getAccuracy calculates accuracy
func getAccuracy(hits, shots int) string {
	if shots == 0 {
		return "0.00"
	}
	accuracy := float64(hits) / float64(shots) * 100
	return fmt.Sprintf("%.2f", accuracy)
}

// listenPlayerShots listens for player shots
func (g *Gui) listenPlayerShots(ctx context.Context, shots chan string) {
	var s []string
loop:

	for {

		select {
		case <-ctx.Done():

			break loop
		default:
			shot := g.opponentBoard.Listen(ctx)
			if shot != "" && !contains(s, shot) {
				s = append(s, shot)
				shots <- shot
			}
		}
	}
}

// drawLegend draws the legend
func (g *Gui) drawLegend() {
	g.gui.Draw(gui.NewText(100, 4, "H - Hit", nil))
	g.gui.Draw(gui.NewText(100, 5, "M - Miss", nil))
	g.gui.Draw(gui.NewText(100, 6, "S - Ship", nil))
	g.gui.Draw(gui.NewText(100, 7, "~ - Empty", nil))
}

// mapStatesToGuiMarks maps states to GUI marks
func mapStatesToGuiMarks(sts [10][10]string) [10][10]gui.State {
	var mapped [10][10]gui.State
	for i, row := range sts {
		for j, s := range row {
			if s == appState.Sunk {
				s = appState.Hit
			}
			mapped[i][j] = gui.State(s)
		}
	}
	return mapped
}

// contains checks if an item is in the set
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
