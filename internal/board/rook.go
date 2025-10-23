package board

import (
	"fmt"
	"strings"
)

type rook struct {
	color    string
	hasMoved bool
}

func (r *rook) colorString() (string, error) {
	if strings.ToLower(r.color) != "white" && strings.ToLower(r.color) != "black" {
		return "", fmt.Errorf("malformed rook struct, incorrect color field")
	}

	return strings.ToLower(r.color), nil
}

func (r *rook) symbol() (rune, error) {
	color, err := r.colorString()
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

func (r *rook) validMove(from, to *position, board *board) bool {
	if from == to {
		return false
	}

	rankDiff := to.rank - from.rank
	fileDiff := to.file - from.file

	if rankDiff != 0 && fileDiff != 0 {
		return false
	}

	return isPathClear(from, to, board)
}

func (r *rook) move(from, to *position, board *board) error {
	if !r.validMove(from, to, board) {
		return fmt.Errorf("invalid move for rook")
	}

	backupPiece := to.piece
	to.piece = r
	from.piece = nil

	var kingPos *position
	switch r.color {
	case "white":
		kingPos = board.whiteKingPosition
	case "black":
		kingPos = board.blackKingPosition
	default:
		return fmt.Errorf("malformed rook struct, incorrect color field")
	}

	err := exposeKing(r.color, board, from, to, kingPos, r, backupPiece)
	if err != nil {
		return err
	} else {
		if backupPiece != nil {
			board.staleTurns = 0
		} else {
			board.staleTurns++
		}
	}

	r.hasMoved = true

	return nil
}
