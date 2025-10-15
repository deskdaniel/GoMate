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

func (k *Knight) ValidMove(from, to *Position, board *Board) bool {
	rankDiff := to.Rank - from.Rank
	fileDiff := to.File - from.File

	if !((rankDiff == 2 && fileDiff == 1) || (rankDiff == 1 && fileDiff == 2)) {
		return false
	}

	if to.Piece != nil {
		targetColor, err := to.Piece.Color()
		if err != nil || targetColor == k.color {
			return false
		}
	}

	return true
}

func (k *Knight) Move(from, to *Position, board *Board) error {
	if !k.ValidMove(from, to, board) {
		return fmt.Errorf("invalid move for knight")
	}

	backupPiece := to.Piece
	to.Piece = k
	from.Piece = nil

	var kingPos *Position
	switch k.color {
	case "white":
		kingPos = board.whiteKingPosition
	case "black":
		kingPos = board.blackKingPosition
	default:
		return fmt.Errorf("malformed knight struct, incorrect color field")
	}

	return exposeKing(k.color, board, from, to, kingPos, k, backupPiece)
}
