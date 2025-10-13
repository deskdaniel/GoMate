package navigation

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dragoo23/Go-chess/internal/app"
	"github.com/dragoo23/Go-chess/internal/board"
	"github.com/dragoo23/Go-chess/internal/game"
	"github.com/dragoo23/Go-chess/internal/messages"
	"github.com/dragoo23/Go-chess/internal/player"
)

type navigationModel struct {
	currentModel tea.Model
	ctx          *app.Context
}

func SetupNavigation(ctx *app.Context) tea.Model {
	m := navigationModel{
		currentModel: game.SetupMainMenu(ctx),
		ctx:          ctx,
	}
	return m
}

func (m navigationModel) Init() tea.Cmd {
	return m.currentModel.Init()
}

func (m navigationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.SwitchToGame:
		m.currentModel = board.NewBoardModel(m.ctx)
		return m, nil
	case messages.SwitchToMainMenu:
		m.currentModel = game.SetupMainMenu(m.ctx)
		return m, nil
	case messages.SwitchToLoginPlayer:
		newModel := player.SetupLogin(m.ctx, msg.Slot)
		if newModel != nil {
			m.currentModel = newModel
		}
		return m, nil
	case messages.SwitchToRegisterUser:
		m.currentModel = player.SetupRegister(m.ctx)
		return m, nil
	case messages.SwitchToStats:
		m.currentModel = player.SetupStats(m.ctx)
		return m, nil
	case messages.SwitchToQuit:
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.currentModel, cmd = m.currentModel.Update(msg)
		return m, cmd
	}
}

func (m navigationModel) View() string {
	return m.currentModel.View()
}
