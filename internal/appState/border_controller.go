package appState

// Directions represents possible directions to check for a ship
var Directions = [][]int{
	{1, 0},
	{0, 1},
	{-1, 0},
	{0, -1},
	{1, 1},
	{-1, 1},
	{-1, -1},
	{1, -1},
}

// DrawBorder draws a border around the ship on the board
func (b *Board) DrawBorder(x, y int) ([][]int, int) {
	// Finds the ship on the board
	shipFound, l := b.LocateShipOnBoard(x, y)
	// For each direction and for each segment of the ship
	for _, v := range Directions {
		for _, s := range shipFound {
			xA, yA := s[0]+v[0], s[1]+v[1]
			// If the coordinates are within range and there is no ship there, mark it as a miss
			if IsWithinBoardLimits(xA, yA) && !IsShipAtCoordinates(xA, yA, b) {
				b.Mark(xA, yA, Miss)
			}
		}
	}
	// Returns the found ship coordinates and its length
	return shipFound, l
}

// LocateShipOnBoard finds the ship on the board
func (b *Board) LocateShipOnBoard(x, y int) ([][]int, int) {
	// Starts from the given coordinates
	shipPlacement := [][]int{{x, y}}
	// For each direction, searches for the ship recursively
	for _, v := range Directions {
		shipPlacement = append(shipPlacement, findShipRecursive(x, y, v, b)...)
	}
	// Returns the ship coordinates and its length
	return shipPlacement, len(shipPlacement)
}

// findShipRecursive searches for the ship recursively in a given direction
func findShipRecursive(x, y int, v []int, b *Board) [][]int {
	// Updates the coordinates
	x, y = x+v[0], y+v[1]
	// If the coordinates are out of range or there is no ship there, returns an empty array
	if !IsWithinBoardLimits(x, y) || !IsShipAtCoordinates(x, y, b) {
		return [][]int{}
	}
	// Adds these coordinates to the array and continues searching in the same direction
	return append([][]int{{x, y}}, findShipRecursive(x, y, v, b)...)
}

// IsWithinBoardLimits checks if the coordinates are within the board range
func IsWithinBoardLimits(x, y int) bool {
	return x >= 0 && x < 10 && y >= 0 && y < 10
}

// IsShipAtCoordinates checks if there is a ship at the given coordinates
func IsShipAtCoordinates(x, y int, b *Board) bool {
	return b.PlayerState[x][y] == Hit || b.PlayerState[x][y] == Sunk
}
