package board

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dragoo23/Go-chess/internal/app"
	"github.com/dragoo23/Go-chess/internal/database"
	"github.com/dragoo23/Go-chess/internal/messages"
	"github.com/google/uuid"
)

type piece interface {
	colorString() (string, error)
	symbol() (rune, error)
	move(from, to *position, board *board) error
	validMove(from, to *position, board *board) bool
}

type position struct {
	rank  int
	file  int
	piece piece
}

type board struct {
	spots             [8][8]*position
	enPassantTarget   *position
	whiteKingPosition *position
	blackKingPosition *position
	staleTurns        int
}

func (p *position) isValid() error {
	if p.rank < 0 || p.rank > 7 {
		return fmt.Errorf("invalid position: rank out of valid range")
	}

	if p.file < 0 || p.file > 7 {
		return fmt.Errorf("invalid position: file out of valid range")
	}

	return nil
}

func (p *position) string() (string, error) {
	err := p.isValid()
	if err != nil {
		return "", err
	}

	rank := p.rank + 1
	file := rune('a' + p.file)

	position := string(file) + strconv.Itoa(rank)
	return position, nil
}

func positionFromString(pos string, m *boardModel) (*position, error) {
	if len(pos) != 2 {
		m.err = fmt.Sprintf("Incorrect position string %q. It must contain 2 characters: letter(a-h) and number (1-8).\n", pos)
		return nil, fmt.Errorf("incorrect position string")
	}

	lowercase := strings.ToLower(pos)
	file := lowercase[0]
	rank := lowercase[1]

	rankInt := int(rank - '1')
	fileInt := int(file - 'a')

	position := position{
		rank: rankInt,
		file: fileInt,
	}

	err := position.isValid()
	if err != nil {
		m.err = fmt.Sprintf("Incorrect position string %q. It must contain 2 characters: letter(a-h) and number (1-8).\n", pos)
		return nil, fmt.Errorf("incorrect position string")
	}

	return &position, nil
}

func initializeBoard() *board {
	b := board{}

	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	b.spots[0][0].piece = &rook{
		color: "white",
	}
	b.spots[0][1].piece = &knight{
		color: "white",
	}
	b.spots[0][2].piece = &bishop{
		color: "white",
	}
	b.spots[0][3].piece = &queen{
		color: "white",
	}
	b.spots[0][4].piece = &king{
		color: "white",
	}
	b.whiteKingPosition = b.spots[0][4]
	b.spots[0][5].piece = &bishop{
		color: "white",
	}
	b.spots[0][6].piece = &knight{
		color: "white",
	}
	b.spots[0][7].piece = &rook{
		color: "white",
	}

	b.spots[7][0].piece = &rook{
		color: "black",
	}
	b.spots[7][1].piece = &knight{
		color: "black",
	}
	b.spots[7][2].piece = &bishop{
		color: "black",
	}
	b.spots[7][3].piece = &queen{
		color: "black",
	}
	b.spots[7][4].piece = &king{
		color: "black",
	}
	b.blackKingPosition = b.spots[7][4]
	b.spots[7][5].piece = &bishop{
		color: "black",
	}
	b.spots[7][6].piece = &knight{
		color: "black",
	}
	b.spots[7][7].piece = &rook{
		color: "black",
	}

	for j := range b.spots[6] {
		b.spots[6][j].piece = &pawn{
			color:     "black",
			direction: -1,
		}
	}

	for j := range b.spots[1] {
		b.spots[1][j].piece = &pawn{
			color:     "white",
			direction: 1,
		}
	}

	return &b
}

func (b *board) renderString() string {
	gameState := ""

	gameState = fmt.Sprintf("  %-2s %-2s %-2s %-2s %-2s %-2s %-2s %-2s\n", "a", "b", "c", "d", "e", "f", "g", "h")
	for rank := 7; rank >= 0; rank-- {
		rankNumber := rank + 1
		gameState += fmt.Sprintf("%d ", rankNumber)
		for file := 0; file < 8; file++ {
			square := b.spots[rank][file]
			if square.piece != nil {
				symbol, err := square.piece.symbol()
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
	board           *board
	err             string
	drawMsg         string
	drawTimer       int
	check           string
	ctx             *app.Context
	whiteTurn       bool
	input           textinput.Model
	promotionSquare *position
	promotionColor  string
	promotionFocus  int
	offeredDraw     bool
	gameOver        bool
	gameOverMsg     string
}

func NewBoardModel(ctx *app.Context) tea.Model {
	board := initializeBoard()
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
	s := m.board.renderString()
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
	if m.gameOver {
		s += fmt.Sprintf("Game over!\n\n%s\n\nPress any key to exit to main menu.", m.gameOverMsg)
		return s
	}
	if m.check != "" {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true).Render(m.check) + "\n"
	}
	if m.drawMsg != "" {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Bold(true).Render(m.drawMsg) + "\n"
	}
	if m.err != "" {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).Render(m.err) + "\n"
	}
	_, warn := check50MoveFule(m.board.staleTurns)
	if warn {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true).Render(fmt.Sprintf("Warning: %d half-moves without pawn movement or capture. Game will be drawn automatically if it reaches 100.", m.board.staleTurns)) + "\n"
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

type overMsg struct {
	winner  *app.User
	loser   *app.User
	draw    bool
	message string
}

func (m *boardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	if m.gameOver {
		switch msg.(type) {
		case tea.KeyMsg:
			return m, func() tea.Msg {
				return messages.SwitchToMainMenu{}
			}
		}
	}

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
					newPiece = &queen{
						color: m.promotionColor,
					}
				case rookField:
					newPiece = &rook{
						color: m.promotionColor,
					}
				case bishopField:
					newPiece = &bishop{
						color: m.promotionColor,
					}
				case knightField:
					newPiece = &knight{
						color: m.promotionColor,
					}
				}

				m.promotionSquare.piece = newPiece
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
			m.err = "To end the game, type 'surrender'/'surr'/'resign'/'forfeit'/'ff' in the input field and press Enter.\nTo offer a draw, type 'draw' and press Enter.\nTo quit the app copletely, close the window."
			return m, nil
		}
	case gameMsg:
		clearEnPassant(m)

		parts := strings.Fields(msg.input)
		if len(parts) == 1 {
			message := strings.ToLower(parts[0])
			var winnerName string
			var loserName string
			switch message {
			case "resign", "surrender", "surr", "forfeit", "ff":
				var winner *app.User
				var loser *app.User
				if m.whiteTurn {
					winner = m.ctx.User2
					if winner != nil {
						winnerName = winner.Username
					} else {
						winnerName = "Guest 2"
					}
					loser = m.ctx.User1
					if loser != nil {
						loserName = loser.Username
					} else {
						loserName = "Guest 1"
					}
				} else {
					winner = m.ctx.User1
					if winner != nil {
						winnerName = winner.Username
					} else {
						winnerName = "Guest 1"
					}
					loser = m.ctx.User2
					if loser != nil {
						loserName = loser.Username
					} else {
						loserName = "Guest 2"
					}
				}
				return m, func() tea.Msg {
					return overMsg{
						winner: winner,
						loser:  loser,
						draw:   false,
						message: fmt.Sprintf("%s has resigned. %s wins!",
							loserName,
							winnerName),
					}
				}
			case "draw":
				if m.offeredDraw {
					message := "Game ended in a draw by agreement."
					m.input.Blur()
					return m, func() tea.Msg {
						return overMsg{
							winner:  nil,
							loser:   nil,
							draw:    true,
							message: message,
						}
					}
				} else {
					m.offeredDraw = true
					m.drawMsg = "Draw offer sent by opponent. You can accept by typing 'draw'."
					m.drawTimer = 1

					switchTurn(m)
					resetInputField(m)

					return m, nil
				}
			}
		}

		if m.offeredDraw {
			m.offeredDraw = false
			m.drawMsg = "Draw offer declined by opponent."
			m.drawTimer = 1
			switchTurn(m)
			resetInputField(m)
			return m, nil
		}

		if len(parts) != 2 {
			m.input.SetValue("")
			return m, nil
		}

		fromStr := parts[0]
		toStr := parts[1]

		fromPos, err := positionFromString(fromStr, m)
		if err != nil {
			m.input.SetValue("")
			return m, nil
		}

		toPos, err := positionFromString(toStr, m)
		if err != nil {
			m.input.SetValue("")
			return m, nil
		}

		fromSquare := m.board.spots[fromPos.rank][fromPos.file]
		toSquare := m.board.spots[toPos.rank][toPos.file]

		if fromSquare.piece == nil {
			m.input.SetValue("")
			return m, nil
		}

		pieceColor, err := fromSquare.piece.colorString()
		if err != nil {
			m.input.SetValue("")
			return m, nil
		}

		if (m.whiteTurn && pieceColor != "white") || (!m.whiteTurn && pieceColor != "black") {
			m.input.SetValue("")
			return m, nil
		}

		err = fromSquare.piece.move(fromSquare, toSquare, m.board)
		if err != nil {
			m.input.SetValue("")
			errString := err.Error()
			m.err = strings.ToUpper(errString[:1]) + errString[1:]
			return m, nil
		}

		switch m.whiteTurn {
		case true:
			if toSquare.rank == 7 {
				if pawn, ok := toSquare.piece.(*pawn); ok && pawn != nil {
					m.promotionSquare = toSquare
					m.promotionColor = "white"
					return m, nil
				}
			}
		case false:
			if toSquare.rank == 0 {
				if pawn, ok := toSquare.piece.(*pawn); ok && pawn != nil {
					m.promotionSquare = toSquare
					m.promotionColor = "black"
					return m, nil
				}
			}
		}

		if !haveSufficientMaterial(m.board) {
			message := "Draw due to insufficient material! Game over."
			m.input.Blur()
			return m, func() tea.Msg {
				return overMsg{
					winner:  nil,
					loser:   nil,
					draw:    true,
					message: message,
				}
			}
		}

		switchTurn(m)

		if m.check != "" {
			color := "white"
			if !m.whiteTurn {
				color = "black"
			}
			if !hasLegalMove(m.board, color) {
				capitalColor := strings.ToUpper(color[:1]) + color[1:]
				message := fmt.Sprintf("%s king is in checkmate! Game over.", capitalColor)
				m.input.Blur()
				var winner *app.User
				var loser *app.User
				if color == "white" {
					winner = m.ctx.User2
					loser = m.ctx.User1
				} else {
					winner = m.ctx.User1
					loser = m.ctx.User2
				}
				return m, func() tea.Msg {
					return overMsg{
						winner:  winner,
						loser:   loser,
						draw:    false,
						message: message,
					}
				}
			}
		}

		if stalemateCheck(m) {
			message := "Draw due to stalemate! Game over."
			m.input.Blur()
			return m, func() tea.Msg {
				return overMsg{
					winner:  nil,
					loser:   nil,
					draw:    true,
					message: message,
				}
			}
		}

		draw, _ := check50MoveFule(m.board.staleTurns)
		if draw {
			message := "Draw due to fifty-move rule! Game over."
			return m, func() tea.Msg {
				return overMsg{
					winner:  nil,
					loser:   nil,
					draw:    true,
					message: message,
				}
			}
		}

		resetInputField(m)
	case overMsg:
		now := sql.NullString{String: time.Now().Format(time.RFC3339), Valid: true}
		if msg.winner != nil {
			winnerRecord, err := m.ctx.Queries.GetRecordsByUserID(context.Background(), msg.winner.ID)
			switch err {
			case nil:
				winnerRecord.Wins.Int64++
				updatedWinner := database.UpdateRecordParams{
					UserID:    msg.winner.ID,
					Wins:      winnerRecord.Wins,
					Losses:    winnerRecord.Losses,
					Draws:     winnerRecord.Draws,
					UpdatedAt: now,
				}
				_, err = m.ctx.Queries.UpdateRecord(context.Background(), updatedWinner)
				if err != nil {
					return m, nil
				}
			case sql.ErrNoRows:
				id, err := uuid.NewUUID()
				if err != nil {
					return m, nil
				}
				winnerRecord := database.RegisterRecordParams{
					ID:        id.String(),
					UserID:    msg.winner.ID,
					Wins:      sql.NullInt64{Int64: 1, Valid: true},
					Losses:    sql.NullInt64{Int64: 0, Valid: true},
					Draws:     sql.NullInt64{Int64: 0, Valid: true},
					CreatedAt: now,
					UpdatedAt: now,
				}
				_, err = m.ctx.Queries.RegisterRecord(context.Background(), winnerRecord)
				if err != nil {
					return m, nil
				}
			default:
				return m, nil
			}
		}
		if msg.loser != nil {
			loserRecord, err := m.ctx.Queries.GetRecordsByUserID(context.Background(), msg.loser.ID)
			switch err {
			case nil:
				loserRecord.Losses.Int64++
				updatedLoser := database.UpdateRecordParams{
					UserID:    msg.loser.ID,
					Wins:      loserRecord.Wins,
					Losses:    loserRecord.Losses,
					Draws:     loserRecord.Draws,
					UpdatedAt: now,
				}
				_, err = m.ctx.Queries.UpdateRecord(context.Background(), updatedLoser)
				if err != nil {
					return m, nil
				}
			case sql.ErrNoRows:
				id, err := uuid.NewUUID()
				if err != nil {
					return m, nil
				}
				loserRecord := database.RegisterRecordParams{
					ID:        id.String(),
					UserID:    msg.loser.ID,
					Wins:      sql.NullInt64{Int64: 0, Valid: true},
					Losses:    sql.NullInt64{Int64: 1, Valid: true},
					Draws:     sql.NullInt64{Int64: 0, Valid: true},
					CreatedAt: now,
					UpdatedAt: now,
				}
				_, err = m.ctx.Queries.RegisterRecord(context.Background(), loserRecord)
				if err != nil {
					return m, nil
				}
			default:
				return m, nil
			}
		}
		if msg.draw {
			if m.ctx.User1 != nil {
				user1Record, err := m.ctx.Queries.GetRecordsByUserID(context.Background(), m.ctx.User1.ID)
				switch err {
				case nil:
					user1Record.Draws.Int64++
					updatedUser1 := database.UpdateRecordParams{
						UserID:    m.ctx.User1.ID,
						Wins:      user1Record.Wins,
						Losses:    user1Record.Losses,
						Draws:     user1Record.Draws,
						UpdatedAt: now,
					}
					_, err = m.ctx.Queries.UpdateRecord(context.Background(), updatedUser1)
					if err != nil {
						return m, nil
					}
				case sql.ErrNoRows:
					id, err := uuid.NewUUID()
					if err != nil {
						return m, nil
					}
					user1Record := database.RegisterRecordParams{
						ID:        id.String(),
						UserID:    m.ctx.User1.ID,
						Wins:      sql.NullInt64{Int64: 0, Valid: true},
						Losses:    sql.NullInt64{Int64: 0, Valid: true},
						Draws:     sql.NullInt64{Int64: 1, Valid: true},
						CreatedAt: now,
						UpdatedAt: now,
					}
					_, err = m.ctx.Queries.RegisterRecord(context.Background(), user1Record)
					if err != nil {
						return m, nil
					}
				default:
					return m, nil
				}
			}
			if m.ctx.User2 != nil {
				user2Record, err := m.ctx.Queries.GetRecordsByUserID(context.Background(), m.ctx.User2.ID)
				switch err {
				case nil:
					user2Record.Draws.Int64++
					updatedUser2 := database.UpdateRecordParams{
						UserID:    m.ctx.User2.ID,
						Wins:      user2Record.Wins,
						Losses:    user2Record.Losses,
						Draws:     user2Record.Draws,
						UpdatedAt: now,
					}
					_, err = m.ctx.Queries.UpdateRecord(context.Background(), updatedUser2)
					if err != nil {
						return m, nil
					}
				case sql.ErrNoRows:
					id, err := uuid.NewUUID()
					if err != nil {
						return m, nil
					}
					user2Record := database.RegisterRecordParams{
						ID:        id.String(),
						UserID:    m.ctx.User2.ID,
						Wins:      sql.NullInt64{Int64: 0, Valid: true},
						Losses:    sql.NullInt64{Int64: 0, Valid: true},
						Draws:     sql.NullInt64{Int64: 1, Valid: true},
						CreatedAt: now,
						UpdatedAt: now,
					}
					_, err = m.ctx.Queries.RegisterRecord(context.Background(), user2Record)
					if err != nil {
						return m, nil
					}
				default:
					return m, nil
				}
			}
		}
		m.gameOver = true
		m.gameOverMsg = msg.message
	}
	return m, cmd
}

func clearEnPassant(m *boardModel) {
	if m.board.enPassantTarget != nil {
		targetColor, err := m.board.enPassantTarget.piece.colorString()
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
	if m.drawTimer > 0 {
		m.drawTimer--
	} else {
		m.drawMsg = ""
	}
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

func hasLegalMove(board *board, color string) bool {
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := board.spots[rank][file]
			if square.piece != nil {
				pieceColor, err := square.piece.colorString()
				if err != nil || pieceColor != color {
					continue
				}

				for toRank := 0; toRank < 8; toRank++ {
					for toFile := 0; toFile < 8; toFile++ {
						toSquare := board.spots[toRank][toFile]
						if square.piece.validMove(square, toSquare, board) {
							movingPiece := square.piece
							capturedPiece := toSquare.piece
							toSquare.piece = movingPiece
							square.piece = nil
							var kingPosition *position
							if _, isKing := movingPiece.(*king); isKing {
								kingPosition = toSquare
							} else if color == "white" {
								kingPosition = board.whiteKingPosition
							} else {
								kingPosition = board.blackKingPosition
							}
							underAttack := isUnderAttack(kingPosition, color, board)
							square.piece = movingPiece
							toSquare.piece = capturedPiece
							if !underAttack {
								return true
							}
						}
					}
				}
			}
		}
	}

	return false
}

func stalemateCheck(m *boardModel) bool {
	if m.check == "" {
		color := "white"
		if !m.whiteTurn {
			color = "black"
		}
		if !hasLegalMove(m.board, color) {
			m.check = fmt.Sprintf("%s is in stalemate! Game over.", strings.ToUpper(color[:1])+color[1:])
			m.input.Blur()
			return true
		}
	}
	return false
}

func haveSufficientMaterial(board *board) bool {
	var minorPiecePositions []*position
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := board.spots[rank][file]
			if square.piece != nil {
				switch square.piece.(type) {
				case *pawn, *rook, *queen:
					return true
				case *bishop:
					minorPiecePositions = append(minorPiecePositions, square)
				case *knight:
					minorPiecePositions = append(minorPiecePositions, square)
				}
			}
		}
	}
	if len(minorPiecePositions) > 2 {
		return true
	}
	if len(minorPiecePositions) <= 1 {
		return false
	}
	piece1 := minorPiecePositions[0].piece
	piece2 := minorPiecePositions[1].piece
	_, ok1 := piece1.(*bishop)
	_, ok2 := piece2.(*bishop)
	if ok1 && ok2 {
		color1 := (minorPiecePositions[0].rank + minorPiecePositions[0].file) % 2
		color2 := (minorPiecePositions[1].rank + minorPiecePositions[1].file) % 2
		return color1 != color2
	}
	return true
}

func check50MoveFule(staleTurns int) (draw bool, warning bool) {
	switch {
	case staleTurns >= 100:
		return true, false
	case staleTurns >= 60:
		return false, true
	default:
		return false, false
	}
}
