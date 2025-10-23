package board

import (
	"fmt"
	"strings"
)

type knight struct {
	color string
}

func (k *knight) colorString() (string, error) {
	if strings.ToLower(k.color) != "white" && strings.ToLower(k.color) != "black" {
		return "", fmt.Errorf("malformed knight struct, incorrect color field")
	}

	return strings.ToLower(k.color), nil
}

func (k *knight) symbol() (rune, error) {
	color, err := k.colorString()
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

func (k *knight) validMove(from, to *position, board *board) bool {
	if from == to {
		return false
	}

	rankDiff := abs(to.rank - from.rank)
	fileDiff := abs(to.file - from.file)

	if !((rankDiff == 2 && fileDiff == 1) || (rankDiff == 1 && fileDiff == 2)) {
		return false
	}

	if to.piece != nil {
		targetColor, err := to.piece.colorString()
		if err != nil || targetColor == k.color {
			return false
		}
	}

	return true
}

func (k *knight) move(from, to *position, board *board) error {
	if !k.validMove(from, to, board) {
		return fmt.Errorf("invalid move for knight")
	}

	backupPiece := to.piece
	to.piece = k
	from.piece = nil

	var kingPos *position
	switch k.color {
	case "white":
		kingPos = board.whiteKingPosition
	case "black":
		kingPos = board.blackKingPosition
	default:
		return fmt.Errorf("malformed knight struct, incorrect color field")
	}

	err := exposeKing(k.color, board, from, to, kingPos, k, backupPiece)
	if err != nil {
		return err
	} else {
		if backupPiece != nil {
			board.staleTurns = 0
		} else {
			board.staleTurns++
		}
	}

	return nil
}
