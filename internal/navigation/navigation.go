package navigation

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deskdaniel/GoMate/internal/app"
	"github.com/deskdaniel/GoMate/internal/board"
	"github.com/deskdaniel/GoMate/internal/game"
	"github.com/deskdaniel/GoMate/internal/help"
	"github.com/deskdaniel/GoMate/internal/messages"
	"github.com/deskdaniel/GoMate/internal/player"
)

type navigationModel struct {
	currentModel tea.Model
	ctx          *app.Context
	viewport     viewport.Model
	width        int
	height       int
}

const minWidth = 30
const minHeight = 15

func SetupNavigation(ctx *app.Context) tea.Model {
	m := navigationModel{
		currentModel: game.SetupMainMenu(ctx),
		ctx:          ctx,
		viewport:     viewport.New(0, 0),
		width:        minWidth,
		height:       minHeight,
	}
	return m
}

func (m navigationModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.currentModel.Init(),
	)
}

func (m navigationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = max(msg.Width, minWidth)
		m.height = max(msg.Height, minHeight)
		// Set viewport dimension, subtract for padding/borders
		m.viewport.Width = m.width - 4
		m.viewport.Height = m.height - 1
		m.viewport.SetContent(m.renderWrappedContent())
		return m, nil
	case messages.SwitchToGame:
		m.currentModel = board.NewBoardModel(m.ctx)
		m.viewport.SetContent(m.renderWrappedContent())
		return m, nil
	case messages.SwitchToMainMenu:
		m.currentModel = game.SetupMainMenu(m.ctx)
		m.viewport.SetContent(m.renderWrappedContent())
		return m, nil
	case messages.SwitchToLoginPlayer:
		newModel := player.SetupLogin(m.ctx, msg.Slot)
		if newModel != nil {
			m.currentModel = newModel
			m.viewport.SetContent(m.renderWrappedContent())
		}
		return m, nil
	case messages.SwitchToRegisterUser:
		m.currentModel = player.SetupRegister(m.ctx)
		m.viewport.SetContent(m.renderWrappedContent())
		return m, nil
	case messages.SwitchToStats:
		m.currentModel = player.SetupStats(m.ctx)
		m.viewport.SetContent(m.renderWrappedContent())
		return m, nil
	case messages.SwitchToHelp:
		m.currentModel = help.SetupHelp(m.ctx)
		m.viewport.SetContent(m.renderWrappedContent())
		return m, nil
	case messages.SwitchToQuit:
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.currentModel, cmd = m.currentModel.Update(msg)
		m.viewport.SetContent(m.renderWrappedContent())
		cmds = append(cmds, cmd)
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m navigationModel) View() string {
	return m.viewport.View()
}

func (m navigationModel) renderWrappedContent() string {
	wrapStyle := lipgloss.NewStyle().Width(max(m.width - 4))
	return wrapStyle.Render(m.currentModel.View())
}
