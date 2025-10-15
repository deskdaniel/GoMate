package board

import (
	"fmt"
	"strings"
)

type Bishop struct {
	color string
}

func (b *Bishop) Color() (string, error) {
	if strings.ToLower(b.color) != "white" && strings.ToLower(b.color) != "black" {
		return "", fmt.Errorf("malformed bishop struct, incorrect color field")
	}

	return strings.ToLower(b.color), nil
}

func (b *Bishop) Symbol() (rune, error) {
	color, err := b.Color()
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

func isPathClear(from, to *Position, board *Board) bool {
	rankDiff := to.Rank - from.Rank
	fileDiff := to.File - from.File

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

	for i, j := from.Rank+rankStep, from.File+fileStep; i != to.Rank+rankStep || j != to.File+fileStep; i, j = i+rankStep, j+fileStep {
		if board.spots[i][j].Piece != nil {
			if i != to.Rank || j != to.File {
				return false
			} else {
				targetColor, err := to.Piece.Color()
				if err != nil {
					return false
				}
				pieceColor, err := from.Piece.Color()
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

func (b *Bishop) ValidMove(from, to *Position, board *Board) bool {
	rankDiff := to.Rank - from.Rank
	fileDiff := to.File - from.File

	if abs(rankDiff) != abs(fileDiff) {
		return false
	}

	// rankStep := 1
	// if rankDiff < 0 {
	// 	rankStep = -1
	// }
	// fileStep := 1
	// if fileDiff < 0 {
	// 	fileStep = -1
	// }

	// for i, j := from.Rank+rankStep, from.File+fileStep; i != to.Rank; i, j = i+rankStep, j+fileStep {
	// 	if board.spots[i][j].Piece != nil {
	// 		return false
	// 	}
	// }

	// if to.Piece != nil {
	// 	targetColor, err := to.Piece.Color()
	// 	if err != nil || targetColor == b.color {
	// 		return false
	// 	}
	// }

	// return true

	return isPathClear(from, to, board)
}

func exposeKing(color string, board *Board, from, to, kingPos *Position, p, backup piece) error {
	if isUnderAttack(kingPos, color, board) {
		from.Piece = p
		to.Piece = backup
		return fmt.Errorf("this move exposes your king")
	}

	return nil
}

func (b *Bishop) Move(from, to *Position, board *Board) error {
	if !b.ValidMove(from, to, board) {
		return fmt.Errorf("invalid move for bishop")
	}

	backupPiece := to.Piece
	to.Piece = b
	from.Piece = nil

	// switch b.color {
	// case "white":
	// 	if isUnderAttack(board.whiteKingPosition, "white", board) {
	// 		from.Piece = b
	// 		to.Piece = backupPiece
	// 		return fmt.Errorf("this move exposes your king")
	// 	}
	// case "black":
	// 	if isUnderAttack(board.blackKingPosition, "black", board) {
	// 		from.Piece = b
	// 		to.Piece = backupPiece
	// 		return fmt.Errorf("this move exposes your king")
	// 	}
	// }
	var kingPos *Position
	switch b.color {
	case "white":
		kingPos = board.whiteKingPosition
	case "black":
		kingPos = board.blackKingPosition
	default:
		return fmt.Errorf("malformed bishop struct, incorrect color field")
	}

	return exposeKing(b.color, board, from, to, kingPos, b, backupPiece)
}
