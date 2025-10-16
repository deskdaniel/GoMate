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
	if from == to {
		return false
	}

	rankDiff := to.Rank - from.Rank
	fileDiff := to.File - from.File

	if rankDiff != 0 && fileDiff != 0 {
		return false
	}

	return isPathClear(from, to, board)
}

func (r *Rook) Move(from, to *Position, board *Board) error {
	if !r.ValidMove(from, to, board) {
		return fmt.Errorf("invalid move for rook")
	}

	backupPiece := to.Piece
	to.Piece = r
	from.Piece = nil

	var kingPos *Position
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
	}

	r.hasMoved = true

	return nil
}
