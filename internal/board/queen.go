package board

import (
	"fmt"
	"strings"
)

type Queen struct {
	color string
}

func (q *Queen) Color() (string, error) {
	if strings.ToLower(q.color) != "white" && strings.ToLower(q.color) != "black" {
		return "", fmt.Errorf("malformed queen struct, incorrect color field")
	}

	return strings.ToLower(q.color), nil
}

func (q *Queen) Symbol() (rune, error) {
	color, err := q.Color()
	if err != nil {
		return 0, err
	}

	switch color {
	case "white":
		return '♕', nil
	case "black":
		return '♛', nil
	default:
		return 0, fmt.Errorf("unexpected piece color: %s", color)
	}
}

func (q *Queen) ValidMove(from, to *Position, board *Board) bool {
	rankDiff := to.Rank - from.Rank
	fileDiff := to.File - from.File

	if rankDiff != 0 && fileDiff != 0 {
		if abs(rankDiff) != abs(fileDiff) {
			return false
		}
	}

	return isPathClear(from, to, board)
}

func (q *Queen) Move(from, to *Position, board *Board) error {
	if !q.ValidMove(from, to, board) {
		return fmt.Errorf("invalid move for queen")
	}

	to.Piece = q
	from.Piece = nil

	return nil
}
