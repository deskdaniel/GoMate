package player

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dragoo23/Go-chess/internal/app"
	"github.com/dragoo23/Go-chess/internal/database"
	"github.com/dragoo23/Go-chess/internal/messages"
	"github.com/google/uuid"
)

type Stats struct {
	Username string
	Wins     int
	Losses   int
	Draws    int
}

func UpdateUserRecord(username string, ctx *app.Context, win, loss, draw bool) error {
	if ctx == nil || ctx.Queries == nil {
		return fmt.Errorf("context or Queries is nil")
	}

	trueCount := 0
	if win {
		trueCount++
	}
	if loss {
		trueCount++
	}
	if draw {
		trueCount++
	}
	if trueCount != 1 {
		return fmt.Errorf("exactly one of win, loss, or draw must be true")
	}

	user, err := ctx.Queries.GetUserByName(context.Background(), username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	stats, err := ctx.Queries.GetRecordsByUserID(context.Background(), user.ID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get stats: %w", err)
	} else if err == sql.ErrNoRows {
		id, err := uuid.NewUUID()
		if err != nil {
			return fmt.Errorf("failed to generate record ID: %w", err)
		}

		statsParams := database.RegisterRecordParams{
			ID:        id.String(),
			UserID:    user.ID,
			CreatedAt: sql.NullString{String: time.Now().Format(time.RFC3339), Valid: true},
			UpdatedAt: sql.NullString{String: time.Now().Format(time.RFC3339), Valid: true},
			Wins:      sql.NullInt64{Int64: 0, Valid: true},
			Losses:    sql.NullInt64{Int64: 0, Valid: true},
			Draws:     sql.NullInt64{Int64: 0, Valid: true},
		}

		if win {
			statsParams.Wins.Int64 = 1
		} else if loss {
			statsParams.Losses.Int64 = 1
		} else if draw {
			statsParams.Draws.Int64 = 1
		}

		_, err = ctx.Queries.RegisterRecord(context.Background(), statsParams)
		if err != nil {
			return fmt.Errorf("failed to create stats record: %w", err)
		}
		return nil
	}

	if win {
		stats.Wins.Int64++
	} else if loss {
		stats.Losses.Int64++
	} else if draw {
		stats.Draws.Int64++
	}
	stats.UpdatedAt = sql.NullString{String: time.Now().Format(time.RFC3339), Valid: true}

	updateParams := database.UpdateRecordParams{
		UpdatedAt: stats.UpdatedAt,
		Wins:      stats.Wins,
		Losses:    stats.Losses,
		Draws:     stats.Draws,
		UserID:    user.ID,
	}
	_, err = ctx.Queries.UpdateRecord(context.Background(), updateParams)
	if err != nil {
		return fmt.Errorf("failed to update stats record: %w", err)
	}

	return nil
}

func CheckStats(username string, ctx *app.Context) (Stats, error) {
	if ctx == nil || ctx.Queries == nil {
		return Stats{}, fmt.Errorf("context or Queries is nil")
	}

	var stats Stats

	user, err := ctx.Queries.GetUserByName(context.Background(), username)
	if err != nil {
		return Stats{}, fmt.Errorf("failed to get user: %w", err)
	}

	stats.Username = user.Username

	sqlStats, err := ctx.Queries.GetRecordsByUserID(context.Background(), user.ID)
	if err == sql.ErrNoRows {
		stats.Wins = 0
		stats.Losses = 0
		stats.Draws = 0
	} else if err != nil {
		return Stats{}, fmt.Errorf("failed to get stats: %w", err)
	} else {
		stats.Wins = int(sqlStats.Wins.Int64)
		stats.Losses = int(sqlStats.Losses.Int64)
		stats.Draws = int(sqlStats.Draws.Int64)
	}

	return stats, nil
}

type statsField int

const (
	user1Field statsField = iota
	user2Field
	inputUsernameField
	quitField
)

type statsModel struct {
	ctx        *app.Context
	fields     []statsField
	focusIndex int
	input      textinput.Model
	err        error
	stats      *Stats
	found      bool
}

func SetupStats(ctx *app.Context) tea.Model {
	username := textinput.New()
	username.Prompt = "Username: "
	username.Placeholder = "Enter username"
	username.CharLimit = 20
	username.Width = 30
	username.Blur()

	var fields []statsField
	if ctx.User1 != nil {
		fields = append(fields, user1Field)
	}
	if ctx.User2 != nil {
		fields = append(fields, user2Field)
	}
	fields = append(fields, inputUsernameField, quitField)

	if fields[0] == inputUsernameField {
		username.Focus()
		username.PromptStyle = username.PromptStyle.Foreground(lipgloss.Color("37"))
		username.TextStyle = username.TextStyle.Foreground(lipgloss.Color("37"))
	}

	m := statsModel{
		ctx:        ctx,
		fields:     fields,
		focusIndex: 0,
		input:      username,
	}

	return &m
}

func (m *statsModel) Init() tea.Cmd {
	if m.fields[m.focusIndex] == inputUsernameField {
		return textinput.Blink
	}
	return nil
}

type statsMsg struct {
	Username string
}

func (m *statsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}
			if m.focusIndex >= len(m.fields) {
				m.focusIndex = 0
			}
			if m.focusIndex < 0 {
				m.focusIndex = len(m.fields) - 1
			}

			if m.fields[m.focusIndex] == inputUsernameField {
				m.input.Focus()
				m.input.PromptStyle = m.input.PromptStyle.Foreground(lipgloss.Color("37"))
				m.input.TextStyle = m.input.TextStyle.Foreground(lipgloss.Color("37"))
			} else {
				m.input.Blur()
				m.input.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
				m.input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
			}

		case "enter":
			switch m.fields[m.focusIndex] {
			case user1Field:
				return m, func() tea.Msg {
					return statsMsg{
						Username: m.ctx.User1.Username,
					}
				}
			case user2Field:
				return m, func() tea.Msg {
					return statsMsg{
						Username: m.ctx.User2.Username,
					}
				}
			case inputUsernameField:
				username := m.input.Value()
				return m, func() tea.Msg {
					return statsMsg{
						Username: username,
					}
				}
			case quitField:
				return m, tea.Quit
			}
		}
	case statsMsg:
		stats, err := CheckStats(msg.Username, m.ctx)
		if err != nil {
			m.err = err
			m.stats = nil
		} else {
			m.err = nil
			m.stats = &stats
			m.found = true
		}

		return m, func() tea.Msg {
			return messages.SwitchToMainMenu{}
		}
	case error:
		m.err = msg
		return m, nil
	}

	if m.fields[m.focusIndex] == inputUsernameField {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, tea.Batch(cmds...)
}

func (m *statsModel) View() string {
	s := ""
	if m.found && m.stats != nil {
		s = fmt.Sprintf("Stats for %s:\n", m.stats.Username)
		if m.stats.Wins == 0 && m.stats.Losses == 0 && m.stats.Draws == 0 {
			s += "No games played yet.\nPress any key to exit.\n"
			return s
		} else {
			s += fmt.Sprintf("Wins: %d\n", m.stats.Wins)
			s += fmt.Sprintf("Losses: %d\n", m.stats.Losses)
			s += fmt.Sprintf("Draws: %d\n", m.stats.Draws)
		}
		s += "Press any key to exit.\n"
		return s
	}

	buttonStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	highlightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("37")).Bold(true)

	s = "Check Player Stats\n\n"
	for _, field := range m.fields {
		var label string
		switch field {
		case user1Field:
			label = m.ctx.User1.Username
		case user2Field:
			label = m.ctx.User2.Username
		case inputUsernameField:
			label = m.input.View()
		case quitField:
			label = "[ Quit ]"
		}

		if field == m.fields[m.focusIndex] {
			s += highlightStyle.Render(label) + "\n"
		} else {
			s += buttonStyle.Render(label) + "\n"
		}
	}

	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
		s += "\n" + errStyle.Render(m.err.Error()) + "\n"
	}

	return s
}
