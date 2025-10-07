package display

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dragoo23/Go-chess/internal/board"
	"github.com/dragoo23/Go-chess/internal/game"
)

type State int

const (
	MainMenu State = iota
	LoggingIn
	Registering
	Playing
	QuitConfirm
)

type Model struct {
	Board       *board.Board
	selected    *board.Position
	turn        string
	inputBuffer string
	commands    []string
	State       State
}

func (m Model) View() string {
	switch m.State {
	case MainMenu:
		s := game.RenderMainMenu()
		s += "\n" + m.inputBuffer
		return s
	case Playing:
		s := m.Board.RenderString()
		s += "\n" + m.inputBuffer
		return s
	default:
		return "Error reading game state"
	}

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyRunes:
			m.inputBuffer += string(msg.Runes)
		case tea.KeyBackspace:
			if len(m.inputBuffer) > 0 {
				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
			}
		case tea.KeyEnter:
			m.commands = strings.Fields(m.inputBuffer)
			m.inputBuffer = ""
		}
	}
	return m, nil
}
