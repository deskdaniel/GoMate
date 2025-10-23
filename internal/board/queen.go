package board

import (
	"fmt"
	"strings"
)

type queen struct {
	color string
}

func (q *queen) colorString() (string, error) {
	if strings.ToLower(q.color) != "white" && strings.ToLower(q.color) != "black" {
		return "", fmt.Errorf("malformed queen struct, incorrect color field")
	}

	return strings.ToLower(q.color), nil
}

func (q *queen) symbol() (rune, error) {
	color, err := q.colorString()
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

func (q *queen) validMove(from, to *position, board *board) bool {
	if from == to {
		return false
	}

	rankDiff := to.rank - from.rank
	fileDiff := to.file - from.file

	if rankDiff != 0 && fileDiff != 0 {
		if abs(rankDiff) != abs(fileDiff) {
			return false
		}
	}

	return isPathClear(from, to, board)
}

func (q *queen) move(from, to *position, board *board) error {
	if !q.validMove(from, to, board) {
		return fmt.Errorf("invalid move for queen")
	}

	backupPiece := to.piece
	to.piece = q
	from.piece = nil

	var kingPos *position
	switch q.color {
	case "white":
		kingPos = board.whiteKingPosition
	case "black":
		kingPos = board.blackKingPosition
	default:
		return fmt.Errorf("malformed queen struct, incorrect color field")
	}

	err := exposeKing(q.color, board, from, to, kingPos, q, backupPiece)
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
