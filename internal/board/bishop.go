package board

import (
	"fmt"
	"strings"
)

type Bishop struct {
	color string
}

func (p *Bishop) Color() (string, error) {
	if strings.ToLower(p.color) != "white" && strings.ToLower(p.color) != "black" {
		return "", fmt.Errorf("malformed bishop struct, incorrect color field")
	}

	return strings.ToLower(p.color), nil
}

func (p *Bishop) Symbol() (rune, error) {
	color, err := p.Color()
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
