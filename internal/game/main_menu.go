package game

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dragoo23/Go-chess/internal/app"
	"github.com/dragoo23/Go-chess/internal/messages"
)

type mainMenuFields int

const (
	startNewGame mainMenuFields = iota
	loginPlayer1
	loginPlayer2
	registerUser
	viewStats
	viewHelp
	quit
)

type mainMenuModel struct {
	focusIndex mainMenuFields
	fields     []mainMenuFields
	ctx        *app.Context
}

func SetupMainMenu(ctx *app.Context) tea.Model {
	fields := []mainMenuFields{
		startNewGame,
		loginPlayer1,
		loginPlayer2,
		registerUser,
		viewStats,
		viewHelp,
		quit,
	}

	m := mainMenuModel{
		focusIndex: startNewGame,
		fields:     fields,
		ctx:        ctx,
	}
	return m
}

func (m mainMenuModel) Init() tea.Cmd {
	return nil
}

func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			return m, func() tea.Msg {
				return messages.SwitchToGame{}
			}
		case "2":
			return m, func() tea.Msg {
				return messages.SwitchToLoginPlayer{Slot: 1}
			}
		case "3":
			return m, func() tea.Msg {
				return messages.SwitchToLoginPlayer{Slot: 2}
			}
		case "4":
			return m, func() tea.Msg {
				return messages.SwitchToRegisterUser{}
			}
		case "5":
			return m, func() tea.Msg {
				return messages.SwitchToStats{}
			}
		case "6":
			return m, func() tea.Msg {
				return messages.SwitchToHelp{}
			}
		case "7", "q", "esc", "ctrl+c":
			return m, func() tea.Msg {
				return messages.SwitchToQuit{}
			}
		case "up", "down":
			s := msg.String()

			if s == "up" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > quit {
				m.focusIndex = startNewGame
			} else if m.focusIndex < startNewGame {
				m.focusIndex = quit
			}

			return m, nil

		case "enter":
			switch m.focusIndex {
			case startNewGame:
				return m, func() tea.Msg {
					return messages.SwitchToGame{}
				}
			case loginPlayer1:
				return m, func() tea.Msg {
					return messages.SwitchToLoginPlayer{Slot: 1}
				}
			case loginPlayer2:
				return m, func() tea.Msg {
					return messages.SwitchToLoginPlayer{Slot: 2}
				}
			case registerUser:
				return m, func() tea.Msg {
					return messages.SwitchToRegisterUser{}
				}
			case viewStats:
				return m, func() tea.Msg {
					return messages.SwitchToStats{}
				}
			case viewHelp:
				return m, func() tea.Msg {
					return messages.SwitchToHelp{}
				}
			case quit:
				return m, func() tea.Msg {
					return messages.SwitchToQuit{}
				}
			}
		}
	}

	return m, nil
}

func (m mainMenuModel) View() string {
	s := "Welcome to Go-Chess!\n\n"

	buttonStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	highlightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("37")).Bold(true)

	for i, field := range m.fields {
		var label string
		switch field {
		case startNewGame:
			label = "1. Start game"
		case loginPlayer1:
			if m.ctx.User1 != nil {
				label = fmt.Sprintf("2. Sign out - %s", m.ctx.User1.Username)
			} else {
				label = "2. Sign in - player 1"
			}
		case loginPlayer2:
			if m.ctx.User2 != nil {
				label = fmt.Sprintf("3. Sign out - %s", m.ctx.User2.Username)
			} else {
				label = "3. Sign in - player 2"
			}
		case registerUser:
			label = "4. Register user"
		case viewStats:
			label = "5. Stats"
		case viewHelp:
			label = "6. Help"
		case quit:
			label = "7. Quit"
		}

		if i == int(m.focusIndex) {
			s += highlightStyle.Render(label) + "\n"
		} else {
			s += buttonStyle.Render(label) + "\n"
		}
	}

	s += "\nUse up/down arrows to navigate, enter to select.\n"
	s += "Alternatively, press the number key for the option.\n"
	s += "Press 7, q, esc or ctrl+c to quit.\n"

	return s
}
