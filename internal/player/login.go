package player

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dragoo23/Go-chess/internal/app"
	"github.com/dragoo23/Go-chess/internal/messages"
)

func checkLogin(ctx *app.Context, slot int) error {
	if ctx == nil || ctx.Queries == nil {
		return fmt.Errorf("context or Queries is nil")
	}

	switch slot {
	case 1:
		if ctx.User1 != nil {
			return fmt.Errorf("player 1 is already logged in")
		}
		if ctx.User2 != nil && ctx.User2.Username == ctx.Username {
			return fmt.Errorf("username already logged in as player 2")
		}
	case 2:
		if ctx.User2 != nil {
			return fmt.Errorf("player 2 is already logged in")
		}
		if ctx.User1 != nil && ctx.User1.Username == ctx.Username {
			return fmt.Errorf("username already logged in as player 1")
		}
	default:
		return fmt.Errorf("invalid slot number")
	}

	return nil
}

func LoginPlayer(ctx *app.Context, slot int) error {
	if ctx == nil || ctx.Queries == nil {
		return fmt.Errorf("context or Queries is nil")
	}

	err := checkLogin(ctx, slot)
	if err != nil {
		return err
	}

	user, err := ctx.Queries.GetUserByName(context.Background(), ctx.Username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	err = CheckPassword(ctx.Password, user.HashedPassword)
	if err != nil {
		return fmt.Errorf("invalid password")
	}

	switch slot {
	case 1:
		ctx.User1 = &app.User{
			ID:       user.ID,
			Username: user.Username,
			Slot:     1,
		}
	case 2:
		ctx.User2 = &app.User{
			ID:       user.ID,
			Username: user.Username,
			Slot:     2,
		}
	default:
		return fmt.Errorf("invalid slot number")
	}

	return nil
}

func LogoutPlayer(ctx *app.Context, slot int) error {
	if ctx == nil || ctx.Queries == nil {
		return fmt.Errorf("context or Queries is nil")
	}

	switch slot {
	case 1:
		if ctx.User1 == nil {
			return fmt.Errorf("player 1 is not logged in")
		}
		ctx.User1 = nil
	case 2:
		if ctx.User2 == nil {
			return fmt.Errorf("player 2 is not logged in")
		}
		ctx.User2 = nil
	default:
		return fmt.Errorf("invalid slot number")
	}

	return nil
}

type loginModel struct {
	focusIndex field
	inputs     []textinput.Model
	ctx        *app.Context
	err        error
	success    bool
	slot       int
}

func (m *loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func SetupLogin(ctx *app.Context, slot int) tea.Model {
	if ctx == nil || ctx.Queries == nil {
		panic("SetupLogin called with nil ctx or nil ctx.Queries")
	}

	m := loginModel{}

	switch slot {
	case 1:
		if ctx.User1 != nil {
			err := LogoutPlayer(ctx, slot)
			if err != nil {
				fmt.Println("Warning: failed to logout existing player 1:", err)
				return nil
			}
			return nil
		}
	case 2:
		if ctx.User2 != nil {
			err := LogoutPlayer(ctx, slot)
			if err != nil {
				fmt.Println("Warning: failed to logout existing player 2:", err)
				return nil
			}
			return nil
		}
	default:
		fmt.Println("SetupLogin called with invalid slot number")
		return nil
	}

	username := textinput.New()
	username.Prompt = "Username: "
	username.Placeholder = "username"
	username.Focus()
	username.CharLimit = 20
	username.Width = 30

	password := textinput.New()
	password.Placeholder = "password"
	password.Prompt = "Password: "
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '*'
	password.CharLimit = 50
	password.Width = 30

	m = loginModel{
		inputs: []textinput.Model{
			username,
			password,
		},
		focusIndex: 0,
		ctx:        ctx,
		slot:       slot,
	}

	m.inputs[usernameField].PromptStyle = m.inputs[usernameField].PromptStyle.Foreground(lipgloss.Color("37"))
	m.inputs[usernameField].TextStyle = m.inputs[usernameField].TextStyle.Foreground(lipgloss.Color("37"))
	m.inputs[passwordField].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	m.inputs[passwordField].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	return &m
}

type loginMsg struct {
	Username string
	Password string
}

func (m *loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.success {
		switch msg.(type) {
		case tea.KeyMsg:
			return m, func() tea.Msg {
				return messages.SwitchToMainMenu{}
			}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, func() tea.Msg {
				return messages.SwitchToMainMenu{}
			}
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == submitField {
				username := m.inputs[usernameField].Value()
				password := m.inputs[passwordField].Value()
				return m, func() tea.Msg {
					return loginMsg{
						Username: username,
						Password: password,
					}
				}
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex == confirmPasswordField {
				if s == "up" || s == "shift+tab" {
					m.focusIndex--
				} else {
					m.focusIndex++
				}
			}
			if m.focusIndex > submitField {
				m.focusIndex = usernameField
			}
			if m.focusIndex < usernameField {
				m.focusIndex = submitField
			}

			for i := 0; i < len(m.inputs); i++ {
				if i == int(m.focusIndex) {
					cmd := m.inputs[i].Focus()
					m.inputs[i].PromptStyle = m.inputs[i].PromptStyle.Foreground(lipgloss.Color("37"))
					m.inputs[i].TextStyle = m.inputs[i].TextStyle.Foreground(lipgloss.Color("37"))
					cmds = append(cmds, cmd)
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
					m.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
				}
			}
			return m, tea.Batch(cmds...)
		}
	case loginMsg:
		m.ctx.Username = msg.Username
		m.ctx.Password = msg.Password

		err := LoginPlayer(m.ctx, m.slot)
		if err != nil {
			m.err = err
			return m, nil
		}

		m.success = true
		return m, nil
	case error:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *loginModel) View() string {
	if m.success {
		s := fmt.Sprintf("User %s logged in successfully!\n\n", m.ctx.Username)

		m.ctx.Username = ""
		m.ctx.Password = ""

		s += "\nPress any key to return to main menu.\n"
		return s
	}

	s := fmt.Sprintf("Logging in as Player%d\n\n", m.slot)
	for i := range m.inputs {
		s += m.inputs[i].View() + "\n\n"
	}

	buttonStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	if m.focusIndex == submitField {
		buttonStyle = buttonStyle.Foreground(lipgloss.Color("37")).Bold(true)
	}
	s += buttonStyle.Render("[ Submit ]") + "\n"

	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
		s += "\n" + errStyle.Render(m.err.Error()) + "\n"
	}

	s += "\nPress Esc to quit.\n"

	return s
}
