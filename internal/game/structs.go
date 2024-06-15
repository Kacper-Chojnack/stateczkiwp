package game

import (
	"battleships/internal/appState"
	"battleships/internal/httpClient"
	"sync"

	gui "github.com/grupawp/warships-gui/v2"
)

// App represents the main structure of the application
type App struct {
	gui                *Gui                       // Game user interface
	game               *httpClient.Game           // Game object
	playerShotsChannel chan string                // Channel for player shots communication
	gameStatusChannel  chan httpClient.GameStatus // Channel for game status communication
	gameStateChannel   chan httpClient.GameState  // Channel for game state communication
	errChan            chan error                 // Channel for error communication
	wg                 *sync.WaitGroup            // WaitGroup for waiting all goroutines to finish
}

// Gui represents the game user interface
type Gui struct {
	gui            *gui.GUI                   // User interface object
	playerBoard    *gui.Board                 // Player's board
	opponentBoard  *gui.Board                 // Opponent's board
	playerNick     *gui.Text                  // Player's nickname
	playerDesc     *gui.Text                  // Player's description
	opponentNick   *gui.Text                  // Opponent's nickname
	opponentDesc   *gui.Text                  // Opponent's description
	turn           *gui.Text                  // Turn information
	timer          *gui.Text                  // Game timer
	waiting        *gui.Text                  // Waiting for opponent information
	numberOf1Ships *gui.Text                  // Number of ships of length 1
	numberOf2Ships *gui.Text                  // Number of ships of length 2
	numberOf3Ships *gui.Text                  // Number of ships of length 3
	numberOf4Ships *gui.Text                  // Number of ships of length 4
	gameStateChan  <-chan *appState.GameState // Channel for game state communication
	timerChan      <-chan int                 // Channel for game time communication
	gameStatusChan chan httpClient.GameStatus // Channel for game status communication
	mu             sync.Mutex                 // Mutex for data access synchronization
}

// GameEvent represents an event in the game
type GameEvent struct {
	PlayerStates   [10][10]string // Player's board state
	OpponentStates [10][10]string // Opponent's board state
	PlayerName     string         // Player's nickname
	PlayerDesc     string         // Player's description
	OpponentName   string         // Opponent's nickname
	OpponentDesc   string         // Opponent's description
	TimeLeft       int            // Remaining game time
	ShouldFire     bool           // Information if player should fire
	GameState      string         // Current game state
	Result         string         // Game result
}
