package board

import (
	"fmt"
	"strings"
)

type bishop struct {
	color string
}

func (b *bishop) colorString() (string, error) {
	if strings.ToLower(b.color) != "white" && strings.ToLower(b.color) != "black" {
		return "", fmt.Errorf("malformed bishop struct, incorrect color field")
	}

	return strings.ToLower(b.color), nil
}

func (b *bishop) symbol() (rune, error) {
	color, err := b.colorString()
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

func isPathClear(from, to *position, board *board) bool {
	rankDiff := to.rank - from.rank
	fileDiff := to.file - from.file

	rankStep := 0
	if rankDiff > 0 {
		rankStep = 1
	} else if rankDiff < 0 {
		rankStep = -1
	}

	fileStep := 0
	if fileDiff > 0 {
		fileStep = 1
	} else if fileDiff < 0 {
		fileStep = -1
	}

	for i, j := from.rank+rankStep, from.file+fileStep; i != to.rank+rankStep || j != to.file+fileStep; i, j = i+rankStep, j+fileStep {
		if board.spots[i][j].piece != nil {
			if i != to.rank || j != to.file {
				return false
			} else {
				targetColor, err := to.piece.colorString()
				if err != nil {
					return false
				}
				pieceColor, err := from.piece.colorString()
				if err != nil {
					return false
				}
				if targetColor == pieceColor {
					return false
				}
			}
		}
	}

	return true
}

func (b *bishop) validMove(from, to *position, board *board) bool {
	if from == to {
		return false
	}

	rankDiff := to.rank - from.rank
	fileDiff := to.file - from.file

	if abs(rankDiff) != abs(fileDiff) {
		return false
	}

	return isPathClear(from, to, board)
}

func exposeKing(color string, board *board, from, to, kingPos *position, p, backup piece) error {
	if isUnderAttack(kingPos, color, board) {
		from.piece = p
		to.piece = backup
		return fmt.Errorf("this move exposes your king")
	}

	return nil
}

func (b *bishop) move(from, to *position, board *board) error {
	if !b.validMove(from, to, board) {
		return fmt.Errorf("invalid move for bishop")
	}

	backupPiece := to.piece
	to.piece = b
	from.piece = nil

	var kingPos *position
	switch b.color {
	case "white":
		kingPos = board.whiteKingPosition
	case "black":
		kingPos = board.blackKingPosition
	default:
		return fmt.Errorf("malformed bishop struct, incorrect color field")
	}

	err := exposeKing(b.color, board, from, to, kingPos, b, backupPiece)
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
