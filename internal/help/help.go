package help

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dragoo23/Go-chess/internal/app"
)

type helpField int

const (
	mainMenu helpField = iota
	coreKnowledge
	gameEndingConditions
	pieceMovement
	quit
	pieceSelection
	bishopHelp
	kingHelp
	knightHelp
	pawnHelp
	queenHelp
	rookHelp
	previous
)

type helpModel struct {
	focusIndex  helpField
	activeField helpField
	activePiece helpField
	ctx         *app.Context
}

func SetupHelp(ctx *app.Context) tea.Model {
	m := &helpModel{
		focusIndex:  0,
		activeField: 0,
		ctx:         ctx,
	}

	return m
}

func (m *helpModel) Init() tea.Cmd {
	return nil
}

func (m *helpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.activeField {
	case mainMenu:
		return handleMainHelpMenu(m, msg)
	case coreKnowledge:
		return handleCoreKnowledge(m, msg)
	case gameEndingConditions:
		return handleGameEndingConditions(m, msg)
	case pieceMovement:
		if m.activePiece != pieceSelection {
			m.activePiece = pieceSelection
		}
		if m.focusIndex < bishopHelp {
			m.focusIndex = bishopHelp
		}
		return handlePieceMovement(m, msg)
	}

	return m, nil
}

func (m *helpModel) View() string {
	var s string
	switch m.activeField {
	case mainMenu:
		s = mainHelpMenuView(m)
	case coreKnowledge:
		s = coreKnowledgeView()
	case gameEndingConditions:
		s = gameEndingConditionsView()
	case pieceMovement:
		s = pieceMovementView(m)
	}
	return s
}
