package board

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dragoo23/Go-chess/internal/app"
	"github.com/dragoo23/Go-chess/internal/messages"
)

type piece interface {
	Color() (string, error)
	Symbol() (rune, error)
	Move(from, to *Position, board *Board) error
	ValidMove(from, to *Position, board *Board) bool
}

type Position struct {
	Rank  int
	File  int
	Piece piece
}

type Board struct {
	spots             [8][8]*Position
	enPassantTarget   *Position
	whiteKingPosition *Position
	blackKingPosition *Position
}

func (p *Position) IsValid() error {
	if p.Rank < 0 || p.Rank > 7 {
		return fmt.Errorf("invalid position: rank out of valid range")
	}

	if p.File < 0 || p.File > 7 {
		return fmt.Errorf("invalid position: file out of valid range")
	}

	return nil
}

func (p *Position) String() (string, error) {
	err := p.IsValid()
	if err != nil {
		return "", err
	}

	rank := p.Rank + 1
	file := rune('a' + p.File)

	position := string(file) + strconv.Itoa(rank)
	return position, nil
}

func PositionFromString(pos string, m *boardModel) (*Position, error) {
	if len(pos) != 2 {
		m.err = fmt.Sprintf("Incorrect position string %q. It must contain 2 characters: letter(a-h) and number (1-8).\n", pos)
		return nil, fmt.Errorf("incorrect position string")
	}

	lowercase := strings.ToLower(pos)
	file := lowercase[0]
	rank := lowercase[1]

	rankInt := int(rank - '1')
	fileInt := int(file - 'a')

	position := Position{
		Rank: rankInt,
		File: fileInt,
	}

	err := position.IsValid()
	if err != nil {
		m.err = fmt.Sprintf("Incorrect position string %q. It must contain 2 characters: letter(a-h) and number (1-8).\n", pos)
		return nil, fmt.Errorf("incorrect position string")
	}

	return &position, nil
}

func InitializeBoard() *Board {
	b := Board{}

	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &Position{
				Rank:  i,
				File:  j,
				Piece: nil,
			}
		}
	}

	b.spots[0][0].Piece = &Rook{
		color: "white",
	}
	b.spots[0][1].Piece = &Knight{
		color: "white",
	}
	b.spots[0][2].Piece = &Bishop{
		color: "white",
	}
	b.spots[0][3].Piece = &Queen{
		color: "white",
	}
	b.spots[0][4].Piece = &King{
		color: "white",
	}
	b.whiteKingPosition = b.spots[0][4]
	b.spots[0][5].Piece = &Bishop{
		color: "white",
	}
	b.spots[0][6].Piece = &Knight{
		color: "white",
	}
	b.spots[0][7].Piece = &Rook{
		color: "white",
	}

	b.spots[7][0].Piece = &Rook{
		color: "black",
	}
	b.spots[7][1].Piece = &Knight{
		color: "black",
	}
	b.spots[7][2].Piece = &Bishop{
		color: "black",
	}
	b.spots[7][3].Piece = &Queen{
		color: "black",
	}
	b.spots[7][4].Piece = &King{
		color: "black",
	}
	b.blackKingPosition = b.spots[7][4]
	b.spots[7][5].Piece = &Bishop{
		color: "black",
	}
	b.spots[7][6].Piece = &Knight{
		color: "black",
	}
	b.spots[7][7].Piece = &Rook{
		color: "black",
	}

	for j := range b.spots[6] {
		b.spots[6][j].Piece = &Pawn{
			color:     "black",
			direction: -1,
		}
	}

	for j := range b.spots[1] {
		b.spots[1][j].Piece = &Pawn{
			color:     "white",
			direction: 1,
		}
	}

	return &b
}

func (b *Board) RenderString() string {
	gameState := ""

	gameState = fmt.Sprintf("  %-2s %-2s %-2s %-2s %-2s %-2s %-2s %-2s\n", "a", "b", "c", "d", "e", "f", "g", "h")
	for rank := 7; rank >= 0; rank-- {
		rankNumber := rank + 1
		gameState += fmt.Sprintf("%d ", rankNumber)
		for file := 0; file < 8; file++ {
			square := b.spots[rank][file]
			if square.Piece != nil {
				symbol, err := square.Piece.Symbol()
				if err == nil {
					gameState += fmt.Sprint(string(symbol)) + "  "
				} else {
					gameState += fmt.Sprintf("%-3s", "?")
				}
			} else {
				if (rank+file)%2 != 1 {
					gameState += fmt.Sprintf("%-3s", "■")
				} else {
					gameState += fmt.Sprintf("%-3s", "□")
				}

			}
		}
		gameState += fmt.Sprintln()
	}
	gameState += fmt.Sprintf("  %-2s %-2s %-2s %-2s %-2s %-2s %-2s %-2s\n\n", "a", "b", "c", "d", "e", "f", "g", "h")

	return gameState
}

type boardModel struct {
	board           *Board
	err             string
	check           string
	ctx             *app.Context
	whiteTurn       bool
	input           textinput.Model
	promotionSquare *Position
	promotionColor  string
	promotionFocus  int
}

func NewBoardModel(ctx *app.Context) tea.Model {
	board := InitializeBoard()
	whiteTurn := true

	input := textinput.New()

	name := ""
	if ctx.User1 != nil {
		name = ctx.User1.Username
	} else {
		name = "Player 1"
	}
	input.Prompt = fmt.Sprintf("%s's(white) turn: ", name)
	input.Placeholder = "Enter command (e.g. A2 A3)"
	input.Focus()
	input.CharLimit = 15
	input.Width = 30

	m := boardModel{
		board:     board,
		ctx:       ctx,
		whiteTurn: whiteTurn,
		input:     input,
	}

	return &m
}

func (m *boardModel) View() string {
	s := m.board.RenderString()
	if m.promotionSquare != nil {
		s += "Pawn promotion! Select a piece to promote to:\n"
		pieces := []string{"Queen", "Rook", "Bishop", "Knight"}
		for i, piece := range pieces {
			if i == m.promotionFocus {
				s += lipgloss.NewStyle().Foreground(lipgloss.Color("37")).Bold(true).Render("> "+piece) + "\n"
			} else {
				s += "  " + piece + "\n"
			}
		}
		return s
	}
	if m.check != "" {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true).Render(m.check) + "\n"
	}
	if m.err != "" {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).Render(m.err) + "\n"
	}
	s += m.input.View() + "\n"
	return s
}

func (m *boardModel) Init() tea.Cmd {
	return textinput.Blink
}

type gameMsg struct {
	input string
}

type promotionField int

const (
	queenField promotionField = iota
	rookField
	bishopField
	knightField
)

// Need checkmate logic
// Need stalemate logic
// Need draw logic (threefold repetition, fifty-move rule, insufficient material)
// Need pawn promotion logic - needs testing
// Need resign logic
func (m *boardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	if m.promotionSquare != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				m.promotionFocus--
				if m.promotionFocus < 0 {
					m.promotionFocus = 3
				}
			case "down":
				m.promotionFocus++
				if m.promotionFocus > 3 {
					m.promotionFocus = 0
				}
			case "enter":
				var newPiece piece
				switch promotionField(m.promotionFocus) {
				case queenField:
					newPiece = &Queen{
						color: m.promotionColor,
					}
				case rookField:
					newPiece = &Rook{
						color: m.promotionColor,
					}
				case bishopField:
					newPiece = &Bishop{
						color: m.promotionColor,
					}
				case knightField:
					newPiece = &Knight{
						color: m.promotionColor,
					}
				}

				m.promotionSquare.Piece = newPiece
				m.promotionSquare = nil
				m.promotionColor = ""
				m.promotionFocus = 0

				switchTurn(m)
				resetInputField(m)

				return m, nil
			}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			input := strings.TrimSpace(m.input.Value())
			if input == "" {
				return m, nil
			}
			return m, func() tea.Msg {
				return gameMsg{
					input: input,
				}
			}
		case "ctrl+c", "esc":
			return m, func() tea.Msg {
				return messages.SwitchToMainMenu{}
			}
		}
	case gameMsg:
		clearEnPassant(m)

		parts := strings.Fields(msg.input)
		if len(parts) != 2 {
			m.input.SetValue("")
			return m, nil
		}

		fromStr := parts[0]
		toStr := parts[1]

		fromPos, err := PositionFromString(fromStr, m)
		if err != nil {
			m.input.SetValue("")
			return m, nil
		}

		toPos, err := PositionFromString(toStr, m)
		if err != nil {
			m.input.SetValue("")
			return m, nil
		}

		fromSquare := m.board.spots[fromPos.Rank][fromPos.File]
		toSquare := m.board.spots[toPos.Rank][toPos.File]

		if fromSquare.Piece == nil {
			m.input.SetValue("")
			return m, nil
		}

		pieceColor, err := fromSquare.Piece.Color()
		if err != nil {
			m.input.SetValue("")
			return m, nil
		}

		if (m.whiteTurn && pieceColor != "white") || (!m.whiteTurn && pieceColor != "black") {
			m.input.SetValue("")
			return m, nil
		}

		err = fromSquare.Piece.Move(fromSquare, toSquare, m.board)
		if err != nil {
			m.input.SetValue("")
			errString := err.Error()
			m.err = strings.ToUpper(errString[:1]) + errString[1:]
			return m, nil
		}

		switch m.whiteTurn {
		case true:
			if toSquare.Rank == 7 && toSquare.Piece.(*Pawn) != nil {
				m.promotionSquare = toSquare
				m.promotionColor = "white"
				return m, nil
			}
		case false:
			if toSquare.Rank == 0 && toSquare.Piece.(*Pawn) != nil {
				m.promotionSquare = toSquare
				m.promotionColor = "black"
				return m, nil
			}
		}

		switchTurn(m)

		resetInputField(m)
	}
	return m, cmd
}

func clearEnPassant(m *boardModel) {
	if m.board.enPassantTarget != nil {
		targetColor, err := m.board.enPassantTarget.Piece.Color()
		if m.whiteTurn {
			if err != nil || targetColor == "white" {
				m.board.enPassantTarget = nil
			}
		} else {
			if err != nil || targetColor == "black" {
				m.board.enPassantTarget = nil
			}
		}
	}
}

func resetInputField(m *boardModel) {
	m.input.SetValue("")
	name := ""
	switch {
	case m.whiteTurn:
		if m.ctx.User1 != nil {
			name = m.ctx.User1.Username
		} else {
			name = "Player 1"
		}
		color := "white"
		m.input.Prompt = fmt.Sprintf("%s's(%s) turn: ", name, color)
		m.input.Placeholder = "Enter command (e.g. A2 A3)"
	case !m.whiteTurn:
		if m.ctx.User2 != nil {
			name = m.ctx.User2.Username
		} else {
			name = "Player 2"
		}
		color := "black"
		m.input.Prompt = fmt.Sprintf("%s's(%s) turn: ", name, color)
		m.input.Placeholder = "Enter command (e.g. A2 A3)"
	}
	m.err = ""
}

func switchTurn(m *boardModel) {
	if m.whiteTurn {
		m.whiteTurn = false
		if isUnderAttack(m.board.blackKingPosition, "black", m.board) {
			m.check = "Black king is under check!"
		} else {
			m.check = ""
		}
	} else {
		m.whiteTurn = true
		if isUnderAttack(m.board.whiteKingPosition, "white", m.board) {
			m.check = "White king is under check!"
		} else {
			m.check = ""
		}
	}
}
