package board

import (
	"fmt"
	"strings"
)

type Pawn struct {
	color     string
	direction int // 1 for white (up), -1 for black (down)
	hasMoved  bool
}

func (p *Pawn) Color() (string, error) {
	if strings.ToLower(p.color) != "white" && strings.ToLower(p.color) != "black" {
		return "", fmt.Errorf("malformed pawn struct, incorrect color field")
	}

	return strings.ToLower(p.color), nil
}

func (p *Pawn) Symbol() (rune, error) {
	color, err := p.Color()
	if err != nil {
		return 0, err
	}

	switch color {
	case "white":
		return '♙', nil
	case "black":
		return '♟', nil
	default:
		return 0, fmt.Errorf("unexpected piece color: %s", color)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (p *Pawn) ValidMove(from, to *Position, board *Board) bool {
	fileDiff := to.File - from.File
	rankDiff := to.Rank - from.Rank

	if fileDiff == 0 {
		if rankDiff == p.direction {
			if board.spots[to.Rank][to.File].Piece == nil {
				return true
			}
		} else if rankDiff == 2*p.direction && !p.hasMoved {
			intermediateRank := from.Rank + p.direction
			if board.spots[intermediateRank][to.File].Piece == nil && board.spots[to.Rank][to.File].Piece == nil {
				return true
			}
		}
	}

	if abs(fileDiff) == 1 && rankDiff == p.direction {
		targetPiece := board.spots[to.Rank][to.File].Piece
		if targetPiece != nil {
			targetColor, err := targetPiece.Color()
			if err == nil && targetColor != p.color {
				return true
			}
		} else if board.enPassantTarget != nil {
			if from.Rank == board.enPassantTarget.Rank && to.File == board.enPassantTarget.File {
				return true
			}
		}
	}

	return false
}

func (p *Pawn) Move(from, to *Position, board *Board) error {
	valid := p.ValidMove(from, to, board)
	if !valid {
		return fmt.Errorf("invalid move for pawn")
	}

	if abs(to.File-from.File) == 1 && board.spots[to.Rank][to.File].Piece == nil {
		capturedRank := from.Rank
		capturedFile := to.File
		capture := board.spots[capturedRank][capturedFile]
		if capture.Piece == nil {
			return fmt.Errorf("no piece to capture en passant")
		}
		capture.Piece = nil
	}

	p.hasMoved = true
	to.Piece = p
	from.Piece = nil

	if abs(to.Rank-from.Rank) == 2 {
		board.enPassantTarget = to
	}

	return nil
}
