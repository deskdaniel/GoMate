package board

import (
	"fmt"
	"strings"
)

type King struct {
	color    string
	hasMoved bool
}

func (k *King) Color() (string, error) {
	if strings.ToLower(k.color) != "white" && strings.ToLower(k.color) != "black" {
		return "", fmt.Errorf("malformed king struct, incorrect color field")
	}

	return strings.ToLower(k.color), nil
}

func (k *King) Symbol() (rune, error) {
	color, err := k.Color()
	if err != nil {
		return 0, err
	}

	switch color {
	case "white":
		return '♔', nil
	case "black":
		return '♚', nil
	default:
		return 0, fmt.Errorf("unexpected piece color: %s", color)
	}
}

func isUnderAttack(pos *Position, color string, board *Board) bool {
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			piece := board.spots[rank][file].Piece
			if piece == nil {
				continue
			}

			pieceColor, err := piece.Color()
			if err != nil || pieceColor == color {
				continue
			}

			from := board.spots[rank][file]
			switch p := piece.(type) {
			case *Pawn:
				rankDiff := pos.Rank - rank
				fileDiff := pos.File - file
				if p.color == "white" && rankDiff == 1 && abs(fileDiff) == 1 {
					return true
				}
				if p.color == "black" && rankDiff == -1 && abs(fileDiff) == 1 {
					return true
				}
			default:
				if piece.ValidMove(from, pos, board) {
					return true
				}
			}
		}
	}

	return false
}

func (k *King) ValidMove(from, to *Position, board *Board) bool {
	rankDiff := to.Rank - from.Rank
	fileDiff := to.File - from.File

	if abs(rankDiff) > 1 || abs(fileDiff) > 1 {
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

func (k *King) Move(from, to *Position, board *Board) error {
	if !k.ValidMove(from, to, board) {
		return fmt.Errorf("invalid move for king")
	}

	if isUnderAttack(to, k.color, board) {
		return fmt.Errorf("king cannot move into check")
	}

	to.Piece = k
	from.Piece = nil
	k.hasMoved = true
	switch k.color {
	case "white":
		board.whiteKingPosition = to
	case "black":
		board.blackKingPosition = to
	default:
		return fmt.Errorf("malformed king struct, incorrect color field")
	}

	return nil
}
