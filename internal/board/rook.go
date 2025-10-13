package board

import (
	"fmt"
	"strings"
)

type Rook struct {
	color string
}

func (r *Rook) Color() (string, error) {
	if strings.ToLower(r.color) != "white" && strings.ToLower(r.color) != "black" {
		return "", fmt.Errorf("malformed rook struct, incorrect color field")
	}

	return strings.ToLower(r.color), nil
}

func (r *Rook) Symbol() (rune, error) {
	color, err := r.Color()
	if err != nil {
		return 0, err
	}

	switch color {
	case "white":
		return '♖', nil
	case "black":
		return '♜', nil
	default:
		return 0, fmt.Errorf("unexpected piece color: %s", color)
	}
}

func (r *Rook) Move(from, to *Position, board *Board) error {
	return nil
}
