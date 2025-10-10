package player

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dragoo23/Go-chess/internal/app"
	"github.com/dragoo23/Go-chess/internal/database"
	"github.com/dragoo23/Go-chess/internal/messages"
	"github.com/google/uuid"
)

func RegisterPlayer(ctx *app.Context) error {
	if ctx == nil || ctx.Queries == nil {
		return fmt.Errorf("context or Queries is nil")
	}

	userName := ctx.Username
	err := checkUsername(userName)
	if err != nil {
		return err
	}

	_, err = ctx.Queries.GetUserByName(context.Background(), userName)
	if err == nil {
		return fmt.Errorf("username already taken")
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check username availability: %w", err)
	}

	err = checkPassword(ctx.Password)
	if err != nil {
		return err
	}

	hashedPassword, err := HashPassword(ctx.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("failed to generate user ID: %w", err)
	}

	userParams := database.RegisterUserParams{
		ID:             id.String(),
		Username:       userName,
		CreatedAt:      sql.NullString{String: time.Now().Format(time.RFC3339), Valid: true},
		UpdatedAt:      sql.NullString{String: time.Now().Format(time.RFC3339), Valid: true},
		HashedPassword: hashedPassword,
	}

	_, err = ctx.Queries.RegisterUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

func checkUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if len(username) > 20 {
		return fmt.Errorf("username cannot exceed 20 characters")
	}
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	valid := regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username)
	if !valid {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}

	return nil
}

func checkPassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if len(password) < 9 {
		return fmt.Errorf("password must be at least 9 characters")
	}
	if len(password) > 50 {
		return fmt.Errorf("password cannot exceed 50 characters")
	}
	valid := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_\-+=\[{\]}\|\\:;"'<,>.?/]+$`).MatchString(password)
	if !valid {
		return fmt.Errorf("password contains invalid characters")
	}

	return nil
}

type field int

const (
	usernameField field = iota
	passwordField
	confirmPasswordField
	submitField
)

type registerModel struct {
	focusIndex field
	inputs     []textinput.Model
	err        error
	ctx        *app.Context
	success    bool
}

func SetupRegister(ctx *app.Context) tea.Model {
	if ctx == nil || ctx.Queries == nil {
		panic("SetupRegister called with nil ctx or nil ctx.Queries")
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

	confirmPassword := textinput.New()
	confirmPassword.Placeholder = "confirm password"
	confirmPassword.Prompt = "Confirm Password: "
	confirmPassword.EchoMode = textinput.EchoPassword
	confirmPassword.EchoCharacter = '*'
	confirmPassword.CharLimit = 50
	confirmPassword.Width = 30

	m := registerModel{
		inputs: []textinput.Model{
			username,
			password,
			confirmPassword,
		},
		focusIndex: usernameField,
		ctx:        ctx,
	}

	m.inputs[usernameField].PromptStyle = m.inputs[usernameField].PromptStyle.Foreground(lipgloss.Color("37"))
	m.inputs[usernameField].TextStyle = m.inputs[usernameField].TextStyle.Foreground(lipgloss.Color("37"))
	m.inputs[passwordField].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	m.inputs[passwordField].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	m.inputs[confirmPasswordField].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	m.inputs[confirmPasswordField].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	return &m
}

func (m *registerModel) Init() tea.Cmd {
	return textinput.Blink
}

type registerMsg struct {
	Username        string
	Password        string
	ConfirmPassword string
}

func (m *registerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == submitField {
				username := m.inputs[usernameField].Value()
				password := m.inputs[passwordField].Value()
				confirmPassword := m.inputs[confirmPasswordField].Value()
				return m, func() tea.Msg {
					return registerMsg{
						Username:        username,
						Password:        password,
						ConfirmPassword: confirmPassword,
					}
				}
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
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
	case registerMsg:
		if msg.Password != msg.ConfirmPassword {
			m.err = fmt.Errorf("passwords do not match")
			return m, nil
		}
		m.ctx.Username = msg.Username
		m.ctx.Password = msg.Password
		err := RegisterPlayer(m.ctx)
		if err != nil {
			m.err = err
			return m, nil
		}

		m.success = true
		return m, func() tea.Msg {
			return messages.SwitchToMainMenu{}
		}
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

func (m *registerModel) View() string {
	if m.success {
		s := fmt.Sprintf("Registration for %s Successful!\n\n", m.ctx.Username)
		s += "You can now log in with your new account.\n"

		m.ctx.Username = ""
		m.ctx.Password = ""

		s += "\nPress any key to exit.\n"
		return s
	}

	s := "Register New Player\n\n"
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
