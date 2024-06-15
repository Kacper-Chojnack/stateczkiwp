package game

import (
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"sync"
)

// PlaceShips places ships on the board
func (a *App) PlaceShips(ctx context.Context) {
	var mutex sync.Mutex

	// Map of ships to place: key is the length of the ship, value is the number of such ships
	ships := map[int]int{
		4: 1,
		3: 2,
		2: 3,
		1: 4,
	}

	currStates := [10][10]gui.State{} // Current state of the board
	newStates := [10][10]gui.State{}  // New state of the board after placing a ship
	var fullCoords []string           // Full list of coordinates of all ships

	board := gui.NewBoard(0, 0, nil)               // Creating a new board
	hint := gui.NewText(50, 0, "Place ships", nil) // Hint for the player
	invalid := gui.NewText(50, 1, "", nil)         // Error message
	placeGui := gui.NewGUI(false)                  // Creating a new user interface
	placeGui.Draw(board)
	placeGui.Draw(hint)

	// Goroutine for placing ships on the board
	go func() {
		for k, v := range ships {
			hint.SetText(fmt.Sprintf("Place %v ship(s) of length %v", v, k))
			placeGui.Draw(hint)
			for i := 0; i < v; i++ {
				var coords []string
				for j := 0; j < k; j++ {
					coord := board.Listen(ctx)
					coords = append(coords, coord)
					x, y := mapToState(coord)

					mutex.Lock()
					newStates[x][y] = gui.Ship
					mutex.Unlock()
					board.SetStates(newStates)
				}
				mutex.Lock()
				valid := isValidPlacement(coords) && !touchesAnotherShip(coords, currStates)
				mutex.Unlock()
				if valid {
					invalid.SetText("")
					currStates = newStates
					fullCoords = append(fullCoords, coords...)
				} else {
					invalid.SetText("Invalid placement, try again")
					placeGui.Draw(invalid)
					for _, coord := range coords {
						x, y := mapToState(coord)
						newStates[x][y] = gui.Empty
					}
					board.SetStates(newStates)
					i--
				}
			}
		}
		hint.SetText("Finished placing ships. Press ctrl+c to save and return to the game!")
		a.game.SetPlayerBoard(fullCoords)
	}()
	placeGui.Start(ctx, nil)
}
