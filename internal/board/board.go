package board

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dragoo23/Go-chess/internal/app"
)

type piece interface {
	Color() (string, error)
	Symbol() (rune, error)
}

type Position struct {
	Rank  int
	File  int
	Piece piece
}

type Board struct {
	spots [8][8]*Position
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

func PositionFromString(pos string) (*Position, error) {
	if len(pos) != 2 {
		fmt.Printf("Incorrect position string %q. It must contain 2 characters: letter(a-h) and number (1-8).\n", pos)
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
		fmt.Printf("Incorrect position string %q. It must contain 2 characters: letter(a-h) and number (1-8).\n", pos)
		return nil, fmt.Errorf("incorrect position string")
	}

	return &position, nil
}

// func (b *Board) PlacePiece(piece piece, pos Position) error {
// 	err := pos.IsValid()
// 	if err != nil {
// 		return err
// 	}

// 	if piece == nil {
// 		return fmt.Errorf("cannot place nil piece")
// 	}

// 	b.spots[pos.Rank][pos.File].Piece = piece
// 	return nil
// }

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
			color: "black",
		}
	}

	for j := range b.spots[1] {
		b.spots[1][j].Piece = &Pawn{
			color: "white",
		}
	}

	return &b
}

// func (b *Board) Render() {
// 	fmt.Printf("  %-2s %-2s %-2s %-2s %-2s %-2s %-2s %-2s\n", "a", "b", "c", "d", "e", "f", "g", "h")
// 	for rank := 7; rank >= 0; rank-- {
// 		rankNumber := rank + 1
// 		fmt.Printf("%d ", rankNumber)
// 		for file := 0; file < 8; file++ {
// 			square := b.spots[rank][file]
// 			var str string
// 			if square.Piece != nil {
// 				symbol, err := square.Piece.Symbol()
// 				if err == nil {
// 					str = string(symbol)
// 				} else {
// 					str = "?"
// 				}
// 			} else {
// 				str = "□"
// 			}
// 			fmt.Printf("%-2s", str)
// 		}
// 		fmt.Println()
// 	}
// 	fmt.Printf("  %-2s %-2s %-2s %-2s %-2s %-2s %-2s %-2s\n", "a", "b", "c", "d", "e", "f", "g", "h")
// }

func (b *Board) Render() {
	fmt.Printf("  %-2s %-2s %-2s %-2s %-2s %-2s %-2s %-2s\n", "a", "b", "c", "d", "e", "f", "g", "h")
	for rank := 7; rank >= 0; rank-- {
		rankNumber := rank + 1
		fmt.Printf("%d ", rankNumber)
		for file := 0; file < 8; file++ {
			square := b.spots[rank][file]
			if square.Piece != nil {
				symbol, err := square.Piece.Symbol()
				if err == nil {
					fmt.Print(string(symbol))
					if symbol == '♙' {
						fmt.Print(" ")
					}
					fmt.Print(" ")
				} else {
					fmt.Print("?")
				}
			} else {
				fmt.Printf("%-2s", "□")
			}
		}
		fmt.Println()
	}
	fmt.Printf("  %-2s %-2s %-2s %-2s %-2s %-2s %-2s %-2s\n", "a", "b", "c", "d", "e", "f", "g", "h")
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
				gameState += fmt.Sprintf("%-3s", "□")
			}
		}
		gameState += fmt.Sprintln()
	}
	gameState += fmt.Sprintf("  %-2s %-2s %-2s %-2s %-2s %-2s %-2s %-2s\n", "a", "b", "c", "d", "e", "f", "g", "h")

	return gameState
}

type boardModel struct {
	board *Board
	// selected  *Position
	ctx       *app.Context
	whiteTurn bool
	input     textinput.Model
	// commands  []string
}

func NewBoardModel(ctx *app.Context) tea.Model {
	board := InitializeBoard()
	whiteTurn := true

	input := textinput.New()
	turn := ""
	color := ""
	if whiteTurn {
		turn = ctx.User1.Username
		color = "white"
	} else {
		turn = ctx.User2.Username
		color = "black"
	}
	input.Prompt = fmt.Sprintf("%s's(%s) turn: ", turn, color)
	input.Placeholder = "Enter command (e.g. A2 A3)"
	input.Focus()
	input.CharLimit = 15
	input.Width = 15

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
	s += m.input.View() + "\n"
	return s
}

func (m *boardModel) Init() tea.Cmd {
	return textinput.Blink
}

type gameMsg struct {
	input string
}

func (m *boardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, tea.Quit
		}
	}
	return m, nil
}
