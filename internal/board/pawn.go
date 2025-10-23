package board

import (
	"fmt"
	"strings"
)

type pawn struct {
	color     string
	direction int // 1 for white (up), -1 for black (down)
	hasMoved  bool
}

func (p *pawn) colorString() (string, error) {
	if strings.ToLower(p.color) != "white" && strings.ToLower(p.color) != "black" {
		return "", fmt.Errorf("malformed pawn struct, incorrect color field")
	}

	return strings.ToLower(p.color), nil
}

func (p *pawn) symbol() (rune, error) {
	color, err := p.colorString()
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

func (p *pawn) validMove(from, to *position, board *board) bool {
	if from == to {
		return false
	}

	fileDiff := to.file - from.file
	rankDiff := to.rank - from.rank

	if fileDiff == 0 {
		if rankDiff == p.direction {
			if board.spots[to.rank][to.file].piece == nil {
				return true
			}
		} else if rankDiff == 2*p.direction && !p.hasMoved {
			intermediateRank := from.rank + p.direction
			if board.spots[intermediateRank][to.file].piece == nil && board.spots[to.rank][to.file].piece == nil {
				return true
			}
		}
	}

	if abs(fileDiff) == 1 && rankDiff == p.direction {
		targetPiece := board.spots[to.rank][to.file].piece
		if targetPiece != nil {
			targetColor, err := targetPiece.colorString()
			if err == nil && targetColor != p.color {
				return true
			}
		} else if board.enPassantTarget != nil {
			if from.rank == board.enPassantTarget.rank && to.file == board.enPassantTarget.file {
				return true
			}
		}
	}

	return false
}

func (p *pawn) move(from, to *position, board *board) error {
	valid := p.validMove(from, to, board)
	if !valid {
		return fmt.Errorf("invalid move for pawn")
	}

	var capturedPiece piece
	var capturedPos *position
	if abs(to.file-from.file) == 1 && to.piece == nil {
		capturedRank := from.rank
		capturedFile := to.file
		capturedPos = board.spots[capturedRank][capturedFile]
		capturedPiece = capturedPos.piece
		if capturedPiece == nil {
			return fmt.Errorf("no piece to capture en passant")
		}
		capturedPos.piece = nil
	}

	backupPiece := to.piece
	to.piece = p
	from.piece = nil

	var kingPos *position
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
			capturedPos.piece = capturedPiece
		}
		return err
	}

	p.hasMoved = true
	if abs(to.rank-from.rank) == 2 {
		board.enPassantTarget = to
	}
	board.staleTurns = 0

	return nil
}
