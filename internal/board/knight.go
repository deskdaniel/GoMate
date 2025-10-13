package board

import (
	"fmt"
	"strings"
)

type Knight struct {
	color string
}

func (k *Knight) Color() (string, error) {
	if strings.ToLower(k.color) != "white" && strings.ToLower(k.color) != "black" {
		return "", fmt.Errorf("malformed knight struct, incorrect color field")
	}

	return strings.ToLower(k.color), nil
}

func (k *Knight) Symbol() (rune, error) {
	color, err := k.Color()
	if err != nil {
		return 0, err
	}

	switch color {
	case "white":
		return '♘', nil
	case "black":
		return '♞', nil
	default:
		return 0, fmt.Errorf("unexpected piece color: %s", color)
	}
}

func (k *Knight) Move(from, to *Position, board *Board) error {
	return nil
}
