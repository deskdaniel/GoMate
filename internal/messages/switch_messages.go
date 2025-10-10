package messages

type SwitchToMainMenu struct{}

type SwitchToGame struct{}

type SwitchToLoginPlayer struct {
	Slot int
}

type SwitchToRegisterUser struct{}

type SwitchToStats struct{}

type SwitchToQuit struct{}
