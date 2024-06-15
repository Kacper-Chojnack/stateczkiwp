package game

import (
	"battleships/internal/httpClient"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"github.com/manifoldco/promptui"
	"sync"
	"time"
)

// NewApp creates a new instance of the application
func NewApp(gameStatusChannel chan httpClient.GameStatus, playerShotsChannel chan string, gameStateChannel chan httpClient.GameState) *App {
	return &App{
		gui:                NewGui(),
		game:               httpClient.NewGame(),
		playerShotsChannel: playerShotsChannel,
		gameStatusChannel:  gameStatusChannel,
		gameStateChannel:   gameStateChannel,
		errChan:            make(chan error),
		wg:                 &sync.WaitGroup{},
	}
}

// InitGameVersusPlayer starts the game for the player
func (a *App) InitGameVersusPlayer(ctx context.Context) {
	for {
		ctx, cancel := context.WithCancel(ctx)
		var wg sync.WaitGroup

		nick, desc := a.game.GetPlayerInfo()

		prompt := promptui.Select{
			Label: "Do you want to place your ships?",
			Items: []string{"Yes", "No"},
		}
		_, answer, err := prompt.Run()
		if err != nil {
			fmt.Printf("Error executing command %v\n", err)
			return
		}
		if answer == "Yes" {
			a.PlaceShips(ctx)
		}
		coords := a.game.GetPlayerCoords()

		promptNick := promptui.Prompt{
			Label: "Enter opponent's nickname (or leave empty to stay in lobby and wait for a challange) ",
		}
		targetNick, err := promptNick.Run()
		if err != nil {
			fmt.Printf("Error executing command %v\n", err)
			return
		}

		a.game.StartGame(nick, desc, targetNick, coords, false)
		board, err := a.game.LoadPlayerBoard()
		if err != nil {
			a.errChan <- err
		}
		_, err = a.game.SetPlayerBoard(board.Board)

		_, err = a.game.GetDescription()

		wg.Add(8)
		a.runGameRoutines(ctx, cancel, &wg)

		a.gui.gui.Start(ctx, nil)
		a.gui.gui.Draw(gui.NewText(1, 0, "Press ctrl+c to leave the game", nil)) // Add this line

		promptAbort := promptui.Select{
			Label: "Abort?",
			Items: []string{"Yes", "No"},
		}
		_, abort, err := promptAbort.Run()
		if err != nil {
			fmt.Printf("Error executing command %v\n", err)
			return
		}
		if abort == "Yes" {
			a.game.AbortGame()
		}

		wg.Wait()

		promptReplay := promptui.Select{
			Label: "Do you want to play again?",
			Items: []string{"Yes", "No"},
		}
		_, choice, err := promptReplay.Run()
		if err != nil {
			fmt.Printf("Error executing command %v\n", err)
			return
		}
		if choice == "No" {
			break
		}
	}
}

// runGameRoutines runs parallel game threads
func (a *App) runGameRoutines(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	go a.runRoutine(ctx, wg, func(ctx context.Context) { a.updateGameStatus(ctx) })
	go a.runRoutine(ctx, wg, func(ctx context.Context) { a.gui.handleGameState(ctx, a.gameStateChannel) })
	go a.runRoutine(ctx, wg, func(ctx context.Context) { a.updateGameState(ctx, cancel) })
	go a.runRoutine(ctx, wg, func(ctx context.Context) { a.gui.displayBoard() })
	go a.runRoutine(ctx, wg, func(ctx context.Context) { a.gui.handleGameStatus(ctx, a.gameStatusChannel) })
	go a.runRoutine(ctx, wg, func(ctx context.Context) { a.handleError(ctx) })
	go a.runRoutine(ctx, wg, func(ctx context.Context) { a.readPlayerShots(ctx) })
	go a.runRoutine(ctx, wg, func(ctx context.Context) { a.gui.listenPlayerShots(ctx, a.playerShotsChannel) })
}

// runRoutine runs a single thread
func (a *App) runRoutine(ctx context.Context, wg *sync.WaitGroup, f func(context.Context)) {
	defer wg.Done()
	f(ctx)
}

// InitGameVersusBot starts the game with a bot
func (a *App) InitGameVersusBot(ctx context.Context) {
	for {
		ctx, cancel := context.WithCancel(ctx)
		var wg sync.WaitGroup
		wg.Add(8)
		nick, desc := a.game.GetPlayerInfo()

		prompt := promptui.Select{
			Label: "Do you want to place your ships?",
			Items: []string{"Yes", "No"},
		}
		_, answer, err := prompt.Run()
		if err != nil {
			fmt.Printf("Error executing command %v\n", err)
			return
		}
		if answer == "Yes" {
			a.PlaceShips(ctx)
		}
		coords := a.game.GetPlayerCoords()

		a.game.StartGame(nick, desc, "", coords, true)
		board, err := a.game.LoadPlayerBoard()
		if err != nil {
			a.errChan <- err
		}
		_, err = a.game.SetPlayerBoard(board.Board)

		a.runGameRoutines(ctx, cancel, &wg)

		a.gui.gui.Start(ctx, nil)

		promptAbort := promptui.Select{
			Label: "You left the game.",
		}
		_, _, err = promptAbort.Run()
		if err != nil {
			a.game.AbortGame()
		}

		wg.Wait()

		a.game.LastGameStatus()

		promptReplay := promptui.Select{
			Label: "Do you want to play again?",
			Items: []string{"Yes", "No"},
		}
		_, choice, err := promptReplay.Run()
		if err != nil {
			fmt.Printf("Error executing command %v\n", err)
			return
		}
		if choice == "No" {
			break
		}
	}
}

// updateGameState updates the game state based on information from the server
func (a *App) updateGameState(ctx context.Context, cancel context.CancelFunc) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Game state update ended")
			break loop
		case <-ticker.C:
			state, err := a.game.GetGameStatus()
			a.game.UpdateLastGameStatus(state.LastGameStatus)
			if err != nil {
				a.errChan <- err
				continue
			}
			if state.GameStatus == "ended" {
				a.game.ClearState()
				cancel()
				return
			}
			if state.GameStatus == "game_in_progress" {
				d, _ := a.game.GetDescription()
				a.game.UpdatePlayersDesc(d)
			}
			oppShots := state.OppShots
			a.game.MarkOpponentShots(oppShots)
			a.gameStatusChannel <- state
		}
	}
}

// updateGameStatus updates the game status based on local data
func (a *App) updateGameStatus(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			state, err := a.game.GetGameState()
			if err != nil {
				a.errChan <- err
				continue
			}
			a.gameStateChannel <- httpClient.GameState{
				PlayerBoard:  state.GetPlayerBoard(),
				OppBoard:     state.GetOpponentBoard(),
				TotalHits:    state.GetTotalHits(),
				TotalShots:   state.GetTotalShots(),
				PlayerDesc:   state.GetPlayerDesc(),
				OppDesc:      state.GetOppDesc(),
				OppShipsSunk: state.RetrieveOpponentSunkShipsCount(),
			}
		}
	}
}

// handleError handles application errors
func (a *App) handleError(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case err := <-a.errChan:
			a.gui.gui.Log("Error: %v", err)
		}
	}
}

// readPlayerShots reads player shots
func (a *App) readPlayerShots(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case shot := <-a.playerShotsChannel:
			_, _, err := a.game.FireShot(shot)
			if err != nil {
				a.errChan <- err
			}
		}
	}
}

// EnterPlayerInfo enters player information
func (a *App) EnterPlayerInfo() {
	promptName := promptui.Prompt{
		Label: "Enter your nickname",
	}

	name, err := promptName.Run()
	if err != nil {
		fmt.Printf("Error executing command %v\n", err)
		return
	}

	promptDescription := promptui.Prompt{
		Label: "Enter your description",
	}

	description, err := promptDescription.Run()
	if err != nil {
		fmt.Printf("Error executing command %v\n", err)
		return
	}

	a.game.UpdatePlayerInfo(name, description)
}

// GetPlayerStats gets player statistics
func (a *App) GetPlayerStats() {
	promptName := promptui.Prompt{
		Label: "Enter player's nickname",
	}

	name, err := promptName.Run()
	if err != nil {
		fmt.Printf("Error executing command %v\n", err)
		return
	}

	stats := a.game.GetPlayerStats(name)
	fmt.Println(stats)
}

// PrintLobby displays players in the lobby
func (a *App) PrintLobby() {
	players, err := a.game.Client.GetLobbyPlayers()
	if err != nil {
		fmt.Printf("Error retrieving players: %v\n", err)
		return
	}

	fmt.Println("Players in the lobby:")
	for _, player := range players {
		fmt.Println(player.Nick)
	}
}
