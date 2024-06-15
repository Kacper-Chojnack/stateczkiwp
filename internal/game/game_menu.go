package game

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
)

// clearScreen clears the console screen
func clearScreen() {
	var cmd *exec.Cmd
	cmd = exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return
	}
}

// DisplayRulesAndDescription displays the game rules and application description
func (a *App) DisplayRulesAndDescription() {
	color.Cyan("Game rules and application description:")
	fmt.Println("Choose a nickname to save your progress. Enter your nickname and description in option 3 in the menu, otherwise you will receive a random nickname and description.")
	fmt.Println("Choose the game mode (with a bot or a real opponent), then choose whether you want to choose the fields with ships yourself, or the game will do it for you.")
	fmt.Println("During the game, to attack the opponent, you have to click on his board.")
	fmt.Println("When you hit correctly, the board will show the symbol H and allow you to shoot again.")
	fmt.Println("You have 60 seconds to take a shot, otherwise you lose.")
	fmt.Println("The game will be won by the person who first sinks all the opponent's ships.")
	fmt.Println("You can leave the game using the keyboard shortcut \"ctrl + c\"")
}

// DisplayMenu displays the game menu
func (a *App) DisplayMenu(ctx context.Context) {

	for {
		clearScreen()
		fmt.Println(" ____        _   _   _           _     _           \n" +
			"|  _ \\      | | | | | |         | |   (_)          \n" +
			"| |_) | __ _| |_| |_| | ___  ___| |__  _ _ __  ___ \n" +
			"|  _ < / _` | __| __| |/ _ \\/ __| '_ \\| | '_ \\/ __|\n" +
			"| |_) | (_| | |_| |_| |  __/\\__ \\ | | | | |_) \\__ \\\n" +
			"|____/ \\__,_|\\__|\\__|_|\\___||___/_| |_|_| .__/|___/\n" +
			"                                        | |        \n" +
			"                                        |_|        ")

		menuItems := []string{
			"Show game rules and application description",
			"Start singleplayer game with bot",
			"Start multiplayer game",
			"Enter player information (nickname and description)",
			"Show top 10 best players",
			"Show player statistics",
			"Show player lobby",
			"Exit",
			"Return to menu",
		}

		prompt := promptui.Select{
			Label: color.GreenString("Welcome to Battleships! Choose an option:"),
			Items: menuItems,
		}

		_, result, err := prompt.Run()
		if err != nil {
			color.Red("Error executing command %v\n", err)
			continue
		}

		a.handleMenuSelection(result, ctx)

		promptContinue := promptui.Prompt{
			Label:     color.YellowString("Press any key to continue"),
			AllowEdit: false,
		}
		_, _ = promptContinue.Run()
	}
}

// handleMenuSelection handles user selection in the menu
func (a *App) handleMenuSelection(selection string, ctx context.Context) {
	switch selection {
	case "Show game rules and application description":
		a.DisplayRulesAndDescription()
	case "Start singleplayer game with bot":
		a.InitGameVersusBot(ctx)
	case "Start multiplayer game":
		a.InitGameVersusPlayer(ctx)
	case "Enter player information (nickname and description)":
		a.EnterPlayerInfo()
	case "Show top 10 best players":
		a.DisplayPlayerRanking()
	case "Show player statistics":
		a.GetPlayerStats()
	case "Show player lobby":
		a.PrintLobby()
	case "Exit":
		a.ExitGame()
	default:
		color.Red("Invalid option.")
	}
}

// ExitGame ends the game
func (a *App) ExitGame() {
	color.Green("See you next time!")
	os.Exit(0)
}

// DisplayPlayerRanking displays player ranking
func (a *App) DisplayPlayerRanking() {
	stats, err := a.game.GetTopPlayerStats()
	if err != nil {
		color.Red("An error occurred:", err)
		return
	}
	fmt.Println(stats)
}

// Menu displays the game menu
func (a *App) Menu(ctx context.Context) {
	a.DisplayMenu(ctx)
}
