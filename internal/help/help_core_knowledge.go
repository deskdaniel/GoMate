package help

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/deskdaniel/GoMate/internal/messages"
)

func handleCoreKnowledge(m *helpModel, msg tea.Msg) (tea.Model, tea.Cmd) {
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

func coreKnowledgeView() string {
	var s string

	s += "Welcome to the 'Core Knowledge' section of the help menu!\n"
	s += "Press `esc` to return to main help screen or `ctrl+c` to return to the main menu.\n\n"

	s += "Chess Basics:\n"
	s += "- Chess is a two-player game: one controls the white pieces, the other the black pieces.\n"
	s += "- Players take turns, with White always moving first.\n"
	s += "- The chessboard is an 8x8 grid with alternating light and dark squares.\n"
	s += "- Rows are called ranks (numbered 1-8, from White's side to Black's).\n"
	s += "- Columns are called files (labeled a-h, from left to right)."
	s += "- White's pieces start on ranks 1 and 2. Black's pieces start on ranks 7 and 8.\n\n"

	s += "Starting Setup:\n"
	s += "- Pawns occupy rank 2 (White) and rank 7 (Black).\n"
	s += "- On rank 1 (White) and rank 8 (Black), pieces are placed as follows:\n"
	s += "\t- Rooks on a1/h1 (White) and a8/h8 (Black).\n"
	s += "\t- Knights on b1/g1 (White) and b8/g8 (Black).\n"
	s += "\t- Bishops on c1/f1 (White) and c8/f8 (Black).\n"
	s += "\t- Queen on d1 (White) and d8 (Black).\n"
	s += "\t- King on e1 (White) and e8 (Black).\n"
	s += "- White's side is at the bottom (ranks 1-2), Black's at the top (ranks 7-8).\n\n"

	s += "Objective:\n"
	s += "- The goal is to checkmate your opponent's king (see 'Game Ending Conditions' section).\n"
	s += "- A game can also end in a draw, where neither player wins (see 'Game Ending Conditions').\n\n"

	s += "Basic Rules:\n"
	s += "- Each turn, you move one piece to an empty square or capture an opponent's piece by moving to its square. Capturing a piece removes it from the chessboard\n"
	s += "- You cannot move to a square occupied by your own pieces.\n"
	s += "- You cannot make a move that puts or leaves your king in check.\n"
	s += "- A 'check' occurs when a move threatens to capture the opponent's king.\n"
	s += "- If your king is in check, you must respond by:\n"
	s += "\t- Moving the king to a safe square.\n"
	s += "\t- Capturing the attacking piece.\n"
	s += "\t- Blocking the attack with another piece (only possible for checks by a rook, bishop, or queen).\n\n"

	s += "Special Rule - Pawn Promotion:\n"
	s += "- If a pawn reaches the opponent's side (rank 8 for White, rank 1 for Black), it must be promoted to a queen, rook, bishop, or knight of the same color\n\n."

	s += "Moving Pieces in This App:\n"
	s += "- Enter two values: the square of the piece you want to move and the destination square (e.g., 'a2 a4' moves a piece from a2 to a4).\n"
	s += "- Squares are identified by file (a-h) and rank (1-8), like 'a2' (file a, rank 2).\n"
	s += "- Note: This app's input format (e.g., 'a2 a4') differs from standard chess notation (e.g., 'Nf3' for knight to f3) used in chess books.\n"
	s += "- See the 'Piece Movement' section for how each piece moves and special moves like castling and en passant."

	return s
}
