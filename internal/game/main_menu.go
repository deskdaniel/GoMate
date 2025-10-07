package game

import (
	"fmt"
)

func RenderMainMenu() string {
	mainmenu := ""
	mainmenu += fmt.Sprintln("\tWelcome to Go-Chess")
	mainmenu += fmt.Sprintln("1. Sign in - player 1")
	mainmenu += fmt.Sprintln("2. Sign in - player 2")
	mainmenu += fmt.Sprintln("3. Start game")
	mainmenu += fmt.Sprintln("4. Register user")
	mainmenu += fmt.Sprintln("5. Stats")
	mainmenu += fmt.Sprintln("6. Quit")

	return mainmenu
}
