package help

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dragoo23/Go-chess/internal/messages"
)

func handleMainHelpMenu(m *helpModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.activeField = coreKnowledge
			return m, nil
		case "2":
			m.activeField = gameEndingConditions
			return m, nil
		case "3":
			m.activeField = pieceMovement
			m.activePiece = pieceSelection
			return m, nil
		case "esc", "ctrl+c", "4":
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

			if m.focusIndex > quit {
				m.focusIndex = coreKnowledge
			} else if m.focusIndex < coreKnowledge {
				m.focusIndex = quit
			}
		case "enter":
			switch m.focusIndex {
			case coreKnowledge:
				m.activeField = coreKnowledge
				return m, nil
			case gameEndingConditions:
				m.activeField = gameEndingConditions
				return m, nil
			case pieceMovement:
				m.activeField = pieceMovement
				m.activePiece = pieceSelection
				return m, nil
			case quit:
				return m, func() tea.Msg {
					return messages.SwitchToMainMenu{}
				}
			case mainMenu:
				return m, nil
			}
		}
	}
	return m, nil
}

func mainHelpMenuView(m *helpModel) string {
	s := ""
	s += "Welcome to help menu!\n"
	s += "Here you will learn everything you need to know to play chess using this app.\n"
	s += "Use `arrow keys` and `enter` or type a number (1-3) to select help section you want to read.\n\n"

	buttonStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	highlightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("37")).Bold(true)

	for field := coreKnowledge; field <= quit; field++ {
		var label string
		switch field {
		case coreKnowledge:
			label = "1. Core Knowledge"
		case gameEndingConditions:
			label = "2. Game Ending Conditions"
		case pieceMovement:
			label = "3. Piece Movement"
		case quit:
			label = "[ Quit ]"
		}

		if m.focusIndex == field {
			s += highlightStyle.Render(label) + "\n"
		} else {
			s += buttonStyle.Render(label) + "\n"
		}
	}

	s += "Press `esc` or `ctrl+c` to exit to main menu.\n"

	return s
}
