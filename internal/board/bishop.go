package board

import (
	"fmt"
	"strings"
)

type Bishop struct {
	color string
}

func (b *Bishop) Color() (string, error) {
	if strings.ToLower(b.color) != "white" && strings.ToLower(b.color) != "black" {
		return "", fmt.Errorf("malformed bishop struct, incorrect color field")
	}

	return strings.ToLower(b.color), nil
}

func (b *Bishop) Symbol() (rune, error) {
	color, err := b.Color()
	if err != nil {
		return 0, err
	}

	switch color {
	case "white":
		return '♗', nil
	case "black":
		return '♝', nil
	default:
		return 0, fmt.Errorf("unexpected piece color: %s", color)
	}
}

func (b *Bishop) Move(from, to *Position, board *Board) error {
	return nil
}
