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
	if from == to {
		return false
	}

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

// reverting move that would expose king in case of en passant captures en passant target anyway
func (p *Pawn) Move(from, to *Position, board *Board) error {
	valid := p.ValidMove(from, to, board)
	if !valid {
		return fmt.Errorf("invalid move for pawn")
	}

	var capturedPiece piece
	var capturedPos *Position
	if abs(to.File-from.File) == 1 && to.Piece == nil {
		capturedRank := from.Rank
		capturedFile := to.File
		capturedPos = board.spots[capturedRank][capturedFile]
		capturedPiece = capturedPos.Piece
		if capturedPiece == nil {
			return fmt.Errorf("no piece to capture en passant")
		}
		capturedPos.Piece = nil
	}

	backupPiece := to.Piece
	to.Piece = p
	from.Piece = nil

	var kingPos *Position
	switch p.color {
	case "white":
		kingPos = board.whiteKingPosition
	case "black":
		kingPos = board.blackKingPosition
	default:
		return fmt.Errorf("malformed pawn struct, incorrect color field")
	}

	err := exposeKing(p.color, board, from, to, kingPos, p, backupPiece)
	if err != nil {
		if capturedPiece != nil {
			capturedPos.Piece = capturedPiece
		}
		return err
	}

	p.hasMoved = true
	if abs(to.Rank-from.Rank) == 2 {
		board.enPassantTarget = to
	}

	return nil
}
