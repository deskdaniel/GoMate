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
	if from == to {
		return false
	}

	rankDiff := to.Rank - from.Rank
	fileDiff := to.File - from.File

	switch k.color {
	case "white":
		if from.File == 4 && from.Rank == 0 && !k.hasMoved {
			if to.File == 6 && to.Rank == 0 {
				rookPos := board.spots[0][7]
				rook, ok := rookPos.Piece.(*Rook)
				if ok && rook != nil && !rook.hasMoved {
					if board.spots[0][5].Piece == nil && board.spots[0][6].Piece == nil {
						if !isUnderAttack(from, k.color, board) &&
							!isUnderAttack(board.spots[0][5], k.color, board) &&
							!isUnderAttack(to, k.color, board) {
							return true
						}
					}
				}
			} else if to.File == 2 && to.Rank == 0 {
				rookPos := board.spots[0][0]
				rook, ok := rookPos.Piece.(*Rook)
				if ok && rook != nil && !rook.hasMoved {
					if board.spots[0][1].Piece == nil && board.spots[0][2].Piece == nil && board.spots[0][3].Piece == nil {
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
		if from.File == 4 && from.Rank == 7 && !k.hasMoved {
			if to.File == 6 && to.Rank == 7 {
				rookPos := board.spots[7][7]
				rook, ok := rookPos.Piece.(*Rook)
				if ok && rook != nil && !rook.hasMoved {
					if board.spots[7][5].Piece == nil && board.spots[7][6].Piece == nil {
						if !isUnderAttack(from, k.color, board) &&
							!isUnderAttack(board.spots[7][5], k.color, board) &&
							!isUnderAttack(to, k.color, board) {
							return true
						}
					}
				}
			} else if to.File == 2 && to.Rank == 7 {
				rookPos := board.spots[7][0]
				rook, ok := rookPos.Piece.(*Rook)
				if ok && rook != nil && !rook.hasMoved {
					if board.spots[7][1].Piece == nil && board.spots[7][2].Piece == nil && board.spots[7][3].Piece == nil {
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

	// if isUnderAttack(to, k.color, board) {
	// 	return fmt.Errorf("king cannot move into check")
	// }

	if abs(to.File-from.File) == 2 {
		switch to.File {
		case 6:
			rookFrom := board.spots[from.Rank][7]
			rookTo := board.spots[from.Rank][5]
			rook, ok := rookFrom.Piece.(*Rook)
			if !ok {
				return fmt.Errorf("no rook to castle with")
			}
			rookFrom.Piece = nil
			rookTo.Piece = rook
			rook.hasMoved = true
		case 2:
			rookFrom := board.spots[from.Rank][0]
			rookTo := board.spots[from.Rank][3]
			rook, ok := rookFrom.Piece.(*Rook)
			if !ok {
				return fmt.Errorf("no rook to castle with")
			}
			rookFrom.Piece = nil
			rookTo.Piece = rook
			rook.hasMoved = true
		default:
			return fmt.Errorf("invalid castling move")
		}
	}

	backupPiece := to.Piece
	to.Piece = k
	from.Piece = nil
	k.hasMoved = true

	if isUnderAttack(to, k.color, board) {
		from.Piece = k
		to.Piece = backupPiece
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
