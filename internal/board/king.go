package board

import (
	"fmt"
	"strings"
)

type king struct {
	color    string
	hasMoved bool
}

func (k *king) colorString() (string, error) {
	if strings.ToLower(k.color) != "white" && strings.ToLower(k.color) != "black" {
		return "", fmt.Errorf("malformed king struct, incorrect color field")
	}

	return strings.ToLower(k.color), nil
}

func (k *king) symbol() (rune, error) {
	color, err := k.colorString()
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

func isUnderAttack(pos *position, color string, board *board) bool {
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			piece := board.spots[rank][file].piece
			if piece == nil {
				continue
			}

			pieceColor, err := piece.colorString()
			if err != nil || pieceColor == color {
				continue
			}

			from := board.spots[rank][file]
			switch p := piece.(type) {
			case *pawn:
				rankDiff := pos.rank - rank
				fileDiff := pos.file - file
				if p.color == "white" && rankDiff == 1 && abs(fileDiff) == 1 {
					return true
				}
				if p.color == "black" && rankDiff == -1 && abs(fileDiff) == 1 {
					return true
				}
			default:
				if piece.validMove(from, pos, board) {
					return true
				}
			}
		}
	}

	return false
}

func (k *king) validMove(from, to *position, board *board) bool {
	if from == to {
		return false
	}

	rankDiff := to.rank - from.rank
	fileDiff := to.file - from.file

	switch k.color {
	case "white":
		if from.file == 4 && from.rank == 0 && !k.hasMoved {
			if to.file == 6 && to.rank == 0 {
				rookPos := board.spots[0][7]
				rook, ok := rookPos.piece.(*rook)
				if ok && rook != nil && !rook.hasMoved {
					if board.spots[0][5].piece == nil && board.spots[0][6].piece == nil {
						if !isUnderAttack(from, k.color, board) &&
							!isUnderAttack(board.spots[0][5], k.color, board) &&
							!isUnderAttack(to, k.color, board) {
							return true
						}
					}
				}
			} else if to.file == 2 && to.rank == 0 {
				rookPos := board.spots[0][0]
				rook, ok := rookPos.piece.(*rook)
				if ok && rook != nil && !rook.hasMoved {
					if board.spots[0][1].piece == nil && board.spots[0][2].piece == nil && board.spots[0][3].piece == nil {
						if !isUnderAttack(from, k.color, board) &&
							!isUnderAttack(board.spots[0][3], k.color, board) &&
							!isUnderAttack(to, k.color, board) {
							return true
						}
					}
				}
			}
		}
	case "black":
		if from.file == 4 && from.rank == 7 && !k.hasMoved {
			if to.file == 6 && to.rank == 7 {
				rookPos := board.spots[7][7]
				rook, ok := rookPos.piece.(*rook)
				if ok && rook != nil && !rook.hasMoved {
					if board.spots[7][5].piece == nil && board.spots[7][6].piece == nil {
						if !isUnderAttack(from, k.color, board) &&
							!isUnderAttack(board.spots[7][5], k.color, board) &&
							!isUnderAttack(to, k.color, board) {
							return true
						}
					}
				}
			} else if to.file == 2 && to.rank == 7 {
				rookPos := board.spots[7][0]
				rook, ok := rookPos.piece.(*rook)
				if ok && rook != nil && !rook.hasMoved {
					if board.spots[7][1].piece == nil && board.spots[7][2].piece == nil && board.spots[7][3].piece == nil {
						if !isUnderAttack(from, k.color, board) &&
							!isUnderAttack(board.spots[7][3], k.color, board) &&
							!isUnderAttack(to, k.color, board) {
							return true
						}
					}
				}
			}
		}
	default:
		return false
	}

	if abs(rankDiff) > 1 || abs(fileDiff) > 1 {
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

func (k *king) move(from, to *position, board *board) error {
	if !k.validMove(from, to, board) {
		return fmt.Errorf("invalid move for king")
	}

	if abs(to.file-from.file) == 2 {
		switch to.file {
		case 6:
			rookFrom := board.spots[from.rank][7]
			rookTo := board.spots[from.rank][5]
			rook, ok := rookFrom.piece.(*rook)
			if !ok {
				return fmt.Errorf("no rook to castle with")
			}
			rookFrom.piece = nil
			rookTo.piece = rook
			rook.hasMoved = true
		case 2:
			rookFrom := board.spots[from.rank][0]
			rookTo := board.spots[from.rank][3]
			rook, ok := rookFrom.piece.(*rook)
			if !ok {
				return fmt.Errorf("no rook to castle with")
			}
			rookFrom.piece = nil
			rookTo.piece = rook
			rook.hasMoved = true
		default:
			return fmt.Errorf("invalid castling move")
		}
	}

	backupPiece := to.piece
	to.piece = k
	from.piece = nil
	k.hasMoved = true

	if isUnderAttack(to, k.color, board) {
		from.piece = k
		to.piece = backupPiece
		k.hasMoved = false
		return fmt.Errorf("king cannot move into check")
	}

	switch k.color {
	case "white":
		board.whiteKingPosition = to
	case "black":
		board.blackKingPosition = to
	default:
		return fmt.Errorf("malformed king struct, incorrect color field")
	}

	if backupPiece != nil {
		board.staleTurns = 0
	} else {
		board.staleTurns++
	}

	return nil
}
