package help

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deskdaniel/GoMate/internal/messages"
)

func handlePieceMovement(m *helpModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			if m.activePiece == pieceSelection {
				m.focusIndex = bishopHelp
				m.activePiece = bishopHelp
			}
			return m, nil
		case "2":
			if m.activePiece == pieceSelection {
				m.focusIndex = kingHelp
				m.activePiece = kingHelp
			}
			return m, nil
		case "3":
			if m.activePiece == pieceSelection {
				m.focusIndex = knightHelp
				m.activePiece = knightHelp
			}
			return m, nil
		case "4":
			if m.activePiece == pieceSelection {
				m.focusIndex = pawnHelp
				m.activePiece = pawnHelp
			}
			return m, nil
		case "5":
			if m.activePiece == pieceSelection {
				m.focusIndex = queenHelp
				m.activePiece = queenHelp
			}
			return m, nil
		case "6":
			if m.activePiece == pieceSelection {
				m.focusIndex = rookHelp
				m.activePiece = rookHelp
			}
			return m, nil
		case "7":
			m.focusIndex = pieceMovement
			m.activeField = mainMenu
			m.activePiece = pieceSelection
			return m, nil
		case "esc":
			if m.activePiece == pieceSelection {
				m.activeField = mainMenu
				m.focusIndex = pieceMovement
				m.activePiece = pieceSelection
				return m, nil
			} else {
				m.focusIndex = m.activePiece
				m.activePiece = pieceSelection
				return m, nil
			}
		case "ctrl+c":
			return m, func() tea.Msg {
				return messages.SwitchToMainMenu{}
			}
		case "up", "down":
			s := msg.String()

			if s == "up" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex < bishopHelp {
				m.focusIndex = previous
			} else if m.focusIndex > previous {
				m.focusIndex = bishopHelp
			}
		case "enter":
			switch m.focusIndex {
			case bishopHelp, kingHelp, knightHelp, pawnHelp, queenHelp, rookHelp:
				m.activePiece = m.focusIndex
				return m, nil
			case previous:
				m.activeField = mainMenu
				m.activePiece = pieceSelection
				return m, nil
			}
		}
	}
	return m, nil
}

func pieceMovementView(m *helpModel) string {
	var s string

	if m.activePiece == pieceSelection {
		s += "Welcome to the 'Piece Movement' section of the help menu!\n"
		s += "Use up/down arrows to switch between options and `enter` to select one or directly type number next to desired option, `esc` or `7` to return to the main help screen, or `ctrl+c` to return to the main menu.\n\n"

		buttonStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		highlightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("37")).Bold(true)
		var label string

		for field := bishopHelp; field <= previous; field++ {
			switch field {
			case bishopHelp:
				label = "1. Bishop (♗/♝)"
			case kingHelp:
				label = "2. King (♔/♚)"
			case knightHelp:
				label = "3. Knight (♘/♞)"
			case pawnHelp:
				label = "4. Pawn (♙/♟)"
			case queenHelp:
				label = "5. Queen (♕/♛)"
			case rookHelp:
				label = "6. Rook (♖/♜)"
			case previous:
				label = "[ Return ]"
			}

			if m.focusIndex == field {
				s += highlightStyle.Render(label) + "\n"
			} else {
				s += buttonStyle.Render(label) + "\n"
			}
		}
	}

	switch m.activePiece {
	case bishopHelp:
		s += "Bishop (white: ♗, black: ♝):\n"
		s += "* The bishop moves any number of squares diagonally, if the path is clear.\n"
		s += "* Each bishop stays on squares of the same color (light or dark) for the entire game.\n"
		s += "* Example: if bishop is on c1 (light square), it can move to f4 (if clear), but not to f3.\n"
		s += pieceHelpReturnInstructions()
	case kingHelp:
		s += "King (white: ♔, black: ♚):\n"
		s += "* The king moves one square in any direction: horizontally, vertically, diagonally.\n"
		s += "* Special move - Castling: The king can move two squares toward a rook (on the same rank), and the rook moves to the king's other side if:\n"
		s += "\t- Neither the king nor the rook has moved yet.\n"
		s += "\t- No pieces are between the king and the rook.\n"
		s += "\t- The king is not in check and the squares it passes through are not attacked.\n"
		s += "\t-Example: For white `e1 g1` castles kingside (rook from h1 to f1) or `e1 c1` castles queenside (rook from a1 to d1).\n"
		s += "*Note: The king cannot move into check. Castling uses a single input for the king's start and end squares.\n"
		s += pieceHelpReturnInstructions()
	case knightHelp:
		s += "Knight (white: ♘, black: ♞):\n"
		s += "* The knight moves in an L-shape: two squares in one direction (horizontally or vertically) and one square perpendicular.\n"
		s += "* The knight can jump over other pieces, making it unique in this aspect.\n"
		s += "* Example: if the knight is on g1, it can move to f3, h3, e2 or e4 (if those squares are empty or occupied by an opponent's piece)/\n"
		s += "* Note: a knight's check cannot be blocked, as it jumps over other pieces.\n"
		s += pieceHelpReturnInstructions()
	case pawnHelp:
		s += "Pawn (white: ♙, black: ♟):\n"
		s += "* Pawns move forward one square (toward the opponent's side - rank 8 for white, rank 1 for black).\n"
		s += "* From starting rank (rank 2 for white, rank 7 for black), pawns can move forward 2 squares.\n"
		s += "* Pawns capture diagonally one square forward.\n"
		s += "* Example: a white pawn on a2 can move to a3 or a4 (if clear), or capture to b3 if an opponent's piece is there.\n"
		s += "* Special move - En Passant: if an opponent's pawn moves two squares forward and lands next to your pawn (in the same rank). You can capture it as if it moved one square, on your next turn only.\n"
		s += "*Example: if a black pawn moves from b7 to b5, white pawn on c5 can capture it with `c5 b6`.\n"
		s += "*Special move - Promotion: when a pawn reaches last rank on opponent's side (8 for white, 1 for black), it must be promoted to a queen, rook, bishop or knight. This app prompts you to select the piece after the move.\n"
		s += pieceHelpReturnInstructions()
	case queenHelp:
		s += "Queen (white: ♕, black: ♛):\n"
		s += "* The queen moves any number of squares horizontally, vertically, or diagonally, if the path is clear.\n"
		s += "* Example: if the queen is on d1, it can move to d8, a1, h1, h5, a4 or any square along those files, ranks or diagonals (if clear).\n"
		s += pieceHelpReturnInstructions()
	case rookHelp:
		s += "Rook (white: ♖, black: ♜):\n"
		s += "* The rook moves any number of squares horizontally or vertically, if the path is clear.\n"
		s += "* Example: If the rook is on a1, it can move to a8, h1, or any square along those files or ranks (if clear).\n"
		s += "* Special Move - Castling: rook takes part in king's special move - castling. This move description can be found in king's section of help.\n"
		s += pieceHelpReturnInstructions()
	}

	s += "\nInput move: Enter the piece's current square and destination, e.g., `a2 a4` moves a piece from a2 to a4.\n"

	return s
}

func pieceHelpReturnInstructions() string {
	var s string

	s += "\nPress `esc` or `7` to return to main help menu, `ctrl+c` to return to main menu, numbers 1-6 to switch to another piece's description or any other key to return to piece selection\n"
	return s
}
