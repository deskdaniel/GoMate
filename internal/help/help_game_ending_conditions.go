package help

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/deskdaniel/GoMate/internal/messages"
)

func handleGameEndingConditions(m *helpModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.activeField = mainMenu
			return m, nil
		case "ctrl+c":
			return m, func() tea.Msg {
				return messages.SwitchToMainMenu{}
			}
		}
	}

	return m, nil
}

func gameEndingConditionsView() string {
	var s string

	s += "Welcome to the 'Game Ending Conditions' section of the help menu!\n"
	s += "Press `esc` to return to main help screen or `ctrl+c` to return to the main menu.\n\n"

	s += "The game ends if one of these conditions is met:\n"
	s += "Win conditions:\n"
	s += "- Checkmate: When a king is in check and the player has no legal moves to escape, the game ends with a win for the player who delivered checkmate.\n"

	s += "Draw conditions:\n"
	s += "- Stalemate: When a player is not in check but has no legal moves, the game ends in a draw (neither player wins).\n"

	s += "- Insufficient Material: The game is a draw if neither player has enough pieces to force checkmate. Examples include:\n"
	s += "\t- King vs. king.\n"
	s += "\t- King and one knight vs. king.\n"
	s += "\t- King and one bishop vs. king.\n"
	s += "\t- King and two bishops on same-colored squares vs. king (checkmate cannot be forced).\n"
	s += "\t- King and bishop(s) vs. king and bishop(s), with all bishops on same-colored squares.\n"

	s += "- 50-Move Rule: If both players make 50 moves each (100 total) without capturing a piece or moving a pawn, the game is a draw.\n"

	s += "- Threefold Repetition: The game is a draw if the same board position occurs three times (not necessarily in a row) with:\n"
	s += "\t- All pieces of the same type and color on the same squares.\n"
	s += "\t- The same player to move.\n"
	s += "\t- The same castling and en passant possibilities.\n"
	s += "\t- Note: This rule is not implemented in this app, but you can offer a draw if you notice the repetition."

	s += "- Mutual Agreement: A player can offer a draw by typing `draw` during their turn.\n"
	s += "\tIf the opponent types `draw` to accept, the game ends in a draw. Any other input refuses the offer, and the game continues.\n"

	s += "Lose conditions:\n"
	s += "- Resignation: A player can resign by typing `forfeit`, `ff`, `surrender`, or `surr` during their turn.\n"
	s += "This results in a loss for the resigning player and a win for the opponent.\n"
	s += "Warning: Resigning is irreversible and does not ask for confirmation!"
	s += "- Time loss: In standard chess, a player loses if they run out of time on a chess clock, unless the opponent has insufficient material (then it's a draw).\n"
	s += "This app does not use time controls, so this rule does not apply here.\n"

	return s
}
