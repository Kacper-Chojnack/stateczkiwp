package game

import (
	gui "github.com/grupawp/warships-gui/v2"
)

// mapToState converts coordinates to board coordinates
func mapToState(cord string) (int, int) {
	if len(cord) > 2 {
		return int(cord[0] - 65), 9
	}
	x := int(cord[0] - 65)
	y := int(cord[1] - 49)
	return x, y
}

// isValidPlacement checks if the ship placement is valid
func isValidPlacement(coords []string) bool {
	if len(coords) == 0 || len(coords) > 4 {
		return false
	}

	x, y := mapToState(coords[0])
	directions := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} // Right, Left, Down, Up

	for i := 1; i < len(coords); i++ {
		xi, yi := mapToState(coords[i])

		if xi != x && yi != y {
			return false // Invalid coordinate placement (not horizontal or vertical)
		}

		found := false
		for _, dir := range directions {
			dx, dy := dir[0], dir[1]
			if xi == x+dx && yi == y+dy {
				found = true
				break
			}
		}

		if !found {
			return false // Invalid coordinate placement (not adjacent)
		}

		x, y = xi, yi
	}

	return true
}

// touchesAnotherShip checks if the ship touches another ship
func touchesAnotherShip(coords []string, states [10][10]gui.State) bool {
	for _, v := range coords {
		x, y := mapToState(v)
		if hasShipAround(x, y, states) {
			return true
		}
	}

	return false
}

// hasShipAround checks if there is a ship around a given point
func hasShipAround(x, y int, states [10][10]gui.State) bool {
	return isShip(x-1, y, states) || // Left
		isShip(x+1, y, states) || // Right
		isShip(x, y-1, states) || // Up
		isShip(x, y+1, states) || // Down
		isShip(x-1, y-1, states) || // Upper left corner
		isShip(x-1, y+1, states) || // Lower left corner
		isShip(x+1, y-1, states) || // Upper right corner
		isShip(x+1, y+1, states) // Lower right corner
}

// isShip checks if there is a ship at a given point
func isShip(x, y int, states [10][10]gui.State) bool {
	if x > 9 || y > 9 || x < 0 || y < 0 {
		return false
	}
	return states[x][y] == gui.Ship
}
