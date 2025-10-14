package board

import (
	"fmt"
	"strings"
)

type Rook struct {
	color    string
	hasMoved bool
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

func (r *Rook) ValidMove(from, to *Position, board *Board) bool {
	rankDiff := to.Rank - from.Rank
	fileDiff := to.File - from.File

	if rankDiff != 0 && fileDiff != 0 {
		return false
	}

	// diff := 0
	// if rankDiff != 0 {
	// 	diff = rankDiff
	// } else {
	// 	diff = fileDiff
	// }
	// step := 1
	// if diff < 0 {
	// 	step = -1
	// }

	// if rankDiff != 0 {
	// 	file := from.File
	// 	for i := from.Rank + step; i != to.Rank; i = i + step {
	// 		if board.spots[i][file].Piece != nil {
	// 			return false
	// 		}
	// 	}
	// } else {
	// 	rank := from.Rank
	// 	for j := from.File + step; j != to.File; j = j + step {
	// 		if board.spots[rank][j].Piece != nil {
	// 			return false
	// 		}
	// 	}
	// }

	// if to.Piece != nil {
	// 	targetColor, err := to.Piece.Color()
	// 	if err != nil || targetColor == r.color {
	// 		return false
	// 	}
	// }

	// return true

	return isPathClear(from, to, board)
}

func (r *Rook) Move(from, to *Position, board *Board) error {
	if !r.ValidMove(from, to, board) {
		return fmt.Errorf("invalid move for rook")
	}

	to.Piece = r
	from.Piece = nil
	r.hasMoved = true

	return nil
}
