package board

import (
	"fmt"
	"strings"
)

type Rook struct {
	color string
}

func (p *Rook) Color() (string, error) {
	if strings.ToLower(p.color) != "white" && strings.ToLower(p.color) != "black" {
		return "", fmt.Errorf("malformed rook struct, incorrect color field")
	}

	return strings.ToLower(p.color), nil
}

func (p *Rook) Symbol() (rune, error) {
	color, err := p.Color()
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

// func NewRook(color string) (Rook, error) {
// 	colorLowercase := strings.ToLower(color)

// 	if colorLowercase != "white" && colorLowercase != "black" {
// 		return Rook{}, fmt.Errorf("incorrect pawn color")
// 	}

// 	rook := Rook{
// 		color: colorLowercase,
// 	}

// 	return rook, nil
// }
