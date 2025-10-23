package board

import (
	"testing"
)

func TestIsValidPosition(t *testing.T) {
	tests := []struct {
		name           string
		positionString string
		wantErr        bool
		position       position
	}{
		{"valid position a1", "a1", false, position{rank: 0, file: 0}},
		{"valid position h8", "h8", false, position{rank: 7, file: 7}},
		{"invalid position i5", "i5", true, position{}},
		{"invalid position a0", "a0", true, position{}},
		{"invalid position empty", "", true, position{}},
		{"invalid position too long", "a10", true, position{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			model := &boardModel{}
			pos, err := positionFromString(test.positionString, model)
			if (err != nil) != test.wantErr {
				t.Errorf("Position from string(%q) error = %v, wantErr %v", test.position, err, test.wantErr)
			}
			if err == nil && test.wantErr == false {
				if pos.rank != test.position.rank || pos.file != test.position.file {
					t.Errorf("Position from string(%q) = %v, want %v", test.positionString, pos, test.position)
				}
			}
		})
	}
}

func TestBishopMoves(t *testing.T) {
	board := board{}

	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	board.spots[0][4].piece = &king{
		color: "white",
	}

	board.whiteKingPosition = board.spots[0][4] // e1
	from := board.spots[0][2]                   // c1
	toValid := board.spots[3][5]                // f4
	toOccupied := board.spots[2][4]             // e3
	toInvalid := board.spots[2][5]              // f3
	toInvalid2 := board.spots[0][5]             // f1

	board.spots[from.rank][from.file].piece = &bishop{
		color: "white",
	}

	t.Run("valid diagonal move", func(t *testing.T) {
		if !from.piece.validMove(from, toValid, &board) {
			bishopStr, _ := from.string()
			destStr, _ := toValid.string()
			t.Errorf("Expected valid move from %v to %v", bishopStr, destStr)
		}
	})

	t.Run("invalid movee, almost diagonal", func(t *testing.T) {
		if from.piece.validMove(from, toInvalid, &board) {
			bishopStr, _ := from.string()
			destStr, _ := toInvalid.string()
			t.Errorf("Expected invalid move from %v to %v", bishopStr, destStr)
		}
	})

	t.Run("invalid move, ortogonal", func(t *testing.T) {
		if from.piece.validMove(from, toInvalid2, &board) {
			bishopStr, _ := from.string()
			destStr, _ := toInvalid2.string()
			t.Errorf("Expected invalid move from %v to %v", bishopStr, destStr)
		}
	})

	// Insert same color piece to block path
	board.spots[2][4].piece = &bishop{
		color: "white",
	}

	t.Run("move occupied by same color", func(t *testing.T) {
		if from.piece.validMove(from, toOccupied, &board) {
			bishopStr, _ := from.string()
			destStr, _ := toOccupied.string()
			t.Errorf("Expected invalid move from %v to %v due to occupied by same color", bishopStr, destStr)
		}
	})

	t.Run("move obstructed by same color", func(t *testing.T) {
		if from.piece.validMove(from, toValid, &board) {
			bishopStr, _ := from.string()
			destStr, _ := toValid.string()
			t.Errorf("Expected invalid move from %v to %v due to obstruction by same color", bishopStr, destStr)
		}
	})

	// New scenario with opponent piece
	board.spots[2][4].piece = nil

	newFrom := board.spots[1][3]
	board.spots[newFrom.rank][newFrom.file].piece = &bishop{
		color: "white",
	}
	opponentFrom := board.spots[2][2]
	board.spots[opponentFrom.rank][opponentFrom.file].piece = &bishop{
		color: "black",
	}
	toExpose := board.spots[3][5]
	t.Run("move esposes king", func(t *testing.T) {
		if newFrom.piece.move(newFrom, toExpose, &board) == nil {
			bishopStr, _ := newFrom.string()
			destStr, _ := toExpose.string()
			t.Errorf("Expected move to expose king from %v to %v", bishopStr, destStr)
		}
	})

	toObstruct := board.spots[3][1]
	t.Run("move blocked by oponent piece", func(t *testing.T) {
		if newFrom.piece.validMove(newFrom, toObstruct, &board) {
			bishopStr, _ := newFrom.string()
			destStr, _ := toObstruct.string()
			t.Errorf("Expected valid move blocked by opponent piece from %v to %v", bishopStr, destStr)
		}
	})

	t.Run("valid capture move", func(t *testing.T) {
		if newFrom.piece.move(newFrom, opponentFrom, &board) != nil {
			bishopStr, _ := newFrom.string()
			destStr, _ := opponentFrom.string()
			t.Errorf("Expected valid capture move from %v to %v", bishopStr, destStr)
		}
	})
}

func TestKingMoves(t *testing.T) {
	board := board{}

	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	kingPos := board.spots[4][4] // e5
	board.spots[kingPos.rank][kingPos.file].piece = &king{
		color: "white",
	}

	board.whiteKingPosition = kingPos

	rookPos := board.spots[3][5] // f4
	board.spots[rookPos.rank][rookPos.file].piece = &rook{
		color: "black",
	}

	validMove1 := board.spots[5][4]   // e6
	validMove2 := board.spots[4][3]   // d5
	validMove3 := board.spots[5][3]   // d6
	invalidMove := board.spots[3][4]  // e4
	invalidMove2 := board.spots[6][4] // e7
	captureMove := board.spots[3][5]  // f4

	t.Run("valid king move up", func(t *testing.T) {
		if !kingPos.piece.validMove(kingPos, validMove1, &board) {
			kingStr, _ := kingPos.string()
			destStr, _ := validMove1.string()
			t.Errorf("Expected valid move from %v to %v", kingStr, destStr)
		}
	})

	t.Run("valid king move left", func(t *testing.T) {
		if !kingPos.piece.validMove(kingPos, validMove2, &board) {
			kingStr, _ := kingPos.string()
			destStr, _ := validMove2.string()
			t.Errorf("Expected valid move from %v to %v", kingStr, destStr)
		}
	})

	t.Run("valid king move diagonally", func(t *testing.T) {
		if !kingPos.piece.validMove(kingPos, validMove3, &board) {
			kingStr, _ := kingPos.string()
			destStr, _ := validMove3.string()
			t.Errorf("Expected valid move from %v to %v", kingStr, destStr)
		}
	})

	t.Run("invalid king move into check", func(t *testing.T) {
		if kingPos.piece.move(kingPos, invalidMove, &board) == nil {
			kingStr, _ := kingPos.string()
			destStr, _ := invalidMove.string()
			t.Errorf("Expected invalid move from %v to %v into check", kingStr, destStr)
		}
	})

	t.Run("invalid king move two squares", func(t *testing.T) {
		if kingPos.piece.validMove(kingPos, invalidMove2, &board) {
			kingStr, _ := kingPos.string()
			destStr, _ := invalidMove2.string()
			t.Errorf("Expected invalid move from %v to %v two squares", kingStr, destStr)
		}
	})

	t.Run("valid king capture move", func(t *testing.T) {
		if kingPos.piece.move(kingPos, captureMove, &board) != nil {
			kingStr, _ := kingPos.string()
			destStr, _ := captureMove.string()
			t.Errorf("Expected valid capture move from %v to %v", kingStr, destStr)
		}
	})
}

func TestCastling(t *testing.T) {
	board := board{}

	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	kingPos := board.spots[0][4] // e1
	board.spots[kingPos.rank][kingPos.file].piece = &king{
		color:    "white",
		hasMoved: false,
	}

	rookPos := board.spots[0][7] // h1
	board.spots[rookPos.rank][rookPos.file].piece = &rook{
		color:    "white",
		hasMoved: false,
	}

	rook2Pos := board.spots[0][0] // a1
	board.spots[rook2Pos.rank][rook2Pos.file].piece = &rook{
		color:    "white",
		hasMoved: false,
	}

	board.whiteKingPosition = kingPos

	castleKingSide := board.spots[0][6]  // g1
	castleQueenSide := board.spots[0][2] // c1
	invalidCastle := board.spots[0][1]   // b1

	t.Run("valid king-side castling", func(t *testing.T) {
		if kingPos.piece.move(kingPos, castleKingSide, &board) != nil {
			kingStr, _ := kingPos.string()
			destStr, _ := castleKingSide.string()
			t.Errorf("Expected valid king-side castling move from %v to %v", kingStr, destStr)
		}
	})

	// Reset positions for queen-side castling test
	board.spots[0][4].piece = &king{
		color:    "white",
		hasMoved: false,
	}
	board.whiteKingPosition = board.spots[0][4]

	t.Run("invalid queen-side castling", func(t *testing.T) {
		if kingPos.piece.move(kingPos, invalidCastle, &board) == nil {
			kingStr, _ := kingPos.string()
			destStr, _ := invalidCastle.string()
			t.Errorf("Expected invalid queen-side castling move from %v to %v due to obstruction", kingStr, destStr)
		}
	})

	t.Run("valid queen-side castling", func(t *testing.T) {
		board.spots[0][1].piece = nil
		board.spots[0][2].piece = nil
		board.spots[0][3].piece = nil

		if kingPos.piece.move(kingPos, castleQueenSide, &board) != nil {
			kingStr, _ := kingPos.string()
			destStr, _ := castleQueenSide.string()
			t.Errorf("Expected valid queen-side castling move from %v to %v", kingStr, destStr)
		}
	})

	// Invalid castling due to checked square
	board.spots[7][5].piece = &rook{
		color: "black",
	}
	board.spots[0][4].piece = &king{
		color:    "white",
		hasMoved: false,
	}
	board.whiteKingPosition = board.spots[0][4]
	t.Run("invalid castling through check", func(t *testing.T) {
		if kingPos.piece.move(kingPos, castleKingSide, &board) == nil {
			kingStr, _ := kingPos.string()
			destStr, _ := castleKingSide.string()
			t.Errorf("Expected invalid castling move from %v to %v through check", kingStr, destStr)
		}
	})

	// Invalid castling due to moved rook
	board.spots[0][0].piece = &rook{
		color:    "white",
		hasMoved: true,
	}
	board.spots[0][4].piece = &king{
		color:    "white",
		hasMoved: false,
	}
	board.whiteKingPosition = board.spots[0][4]
	t.Run("invalid castling with moved rook", func(t *testing.T) {
		if kingPos.piece.move(kingPos, castleQueenSide, &board) == nil {
			kingStr, _ := kingPos.string()
			destStr, _ := castleQueenSide.string()
			t.Errorf("Expected invalid castling move from %v to %v with moved rook", kingStr, destStr)
		}
	})

	// Invalid castling due to moved king
	board.spots[0][0].piece = &rook{
		color:    "white",
		hasMoved: false,
	}
	board.spots[0][4].piece = &king{
		color:    "white",
		hasMoved: true,
	}
	board.whiteKingPosition = board.spots[0][4]
	t.Run("invalid castling with moved king", func(t *testing.T) {
		if kingPos.piece.move(kingPos, castleQueenSide, &board) == nil {
			kingStr, _ := kingPos.string()
			destStr, _ := castleQueenSide.string()
			t.Errorf("Expected invalid castling move from %v to %v with moved king", kingStr, destStr)
		}
	})

	// Invalid casting due to obstruction
	board.spots[0][1].piece = &bishop{
		color: "white",
	}
	board.spots[0][4].piece = &king{
		color:    "white",
		hasMoved: false,
	}
	board.whiteKingPosition = board.spots[0][4]
	t.Run("invalid castling with obstruction", func(t *testing.T) {
		if kingPos.piece.move(kingPos, castleQueenSide, &board) == nil {
			kingStr, _ := kingPos.string()
			destStr, _ := castleQueenSide.string()
			t.Errorf("Expected invalid castling move from %v to %v with obstruction", kingStr, destStr)
		}
	})

	// Invalid castling due to check
	board.spots[7][4].piece = &rook{
		color: "black",
	}
	board.spots[0][4].piece = &king{
		color:    "white",
		hasMoved: false,
	}
	board.whiteKingPosition = board.spots[0][4]
	t.Run("invalid castling while in check", func(t *testing.T) {
		if kingPos.piece.move(kingPos, castleQueenSide, &board) == nil {
			kingStr, _ := kingPos.string()
			destStr, _ := castleQueenSide.string()
			t.Errorf("Expected invalid castling move from %v to %v while in check", kingStr, destStr)
		}
	})
}

func TestKnightMoves(t *testing.T) {
	board := board{}

	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	from := board.spots[4][4] // e5
	board.spots[from.rank][from.file].piece = &knight{
		color: "white",
	}

	validMove := board.spots[6][5]   // f7
	validMove2 := board.spots[5][6]  // g6
	invalidMove := board.spots[5][5] // f6
	occupiedBySameColor := board.spots[5][6]
	occupiedByOpponent := board.spots[3][6]

	t.Run("valid knight move", func(t *testing.T) {
		if !from.piece.validMove(from, validMove, &board) {
			fromStr, _ := from.string()
			destStr, _ := validMove.string()
			t.Errorf("Expected valid move from %v to %v", fromStr, destStr)
		}
	})

	t.Run("valid knight move 2", func(t *testing.T) {
		if !from.piece.validMove(from, validMove2, &board) {
			fromStr, _ := from.string()
			destStr, _ := validMove2.string()
			t.Errorf("Expected valid move from %v to %v", fromStr, destStr)
		}
	})

	t.Run("invalid knight move", func(t *testing.T) {
		if from.piece.validMove(from, invalidMove, &board) {
			fromStr, _ := from.string()
			destStr, _ := invalidMove.string()
			t.Errorf("Expected invalid move from %v to %v", fromStr, destStr)
		}
	})

	// Insert same color piece to test occupation
	board.spots[occupiedBySameColor.rank][occupiedBySameColor.file].piece = &bishop{
		color: "white",
	}

	t.Run("knight move occupied by same color", func(t *testing.T) {
		if from.piece.validMove(from, occupiedBySameColor, &board) {
			fromStr, _ := from.string()
			destStr, _ := occupiedBySameColor.string()
			t.Errorf("Expected invalid move from %v to %v due to occupation by same color", fromStr, destStr)
		}
	})

	// Insert opponent piece to test capture
	board.spots[occupiedByOpponent.rank][occupiedByOpponent.file].piece = &bishop{
		color: "black",
	}

	t.Run("knight capture move", func(t *testing.T) {
		if from.piece.move(from, occupiedByOpponent, &board) != nil {
			fromStr, _ := from.string()
			destStr, _ := occupiedByOpponent.string()
			t.Errorf("Expected valid capture move from %v to %v", fromStr, destStr)
		}
	})

	// New scenario to test exposing king
	board.spots[0][4].piece = &king{
		color: "white",
	}
	board.whiteKingPosition = board.spots[0][4]
	newFrom := board.spots[1][5]
	board.spots[newFrom.rank][newFrom.file].piece = &knight{
		color: "white",
	}
	board.spots[2][6].piece = &bishop{
		color: "black",
	}

	toExpose := board.spots[2][7]
	t.Run("knight move exposes king", func(t *testing.T) {
		if newFrom.piece.move(newFrom, toExpose, &board) == nil {
			fromStr, _ := newFrom.string()
			destStr, _ := toExpose.string()
			t.Errorf("Expected move to expose king from %v to %v", fromStr, destStr)
		}
	})
}

func TestQueenMoves(t *testing.T) {
	board := board{}

	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	from := board.spots[4][4] // e5
	board.spots[from.rank][from.file].piece = &queen{
		color: "white",
	}

	validMoveDiagonal := board.spots[6][6] // g7
	validMoveStraight := board.spots[4][7] // h5
	invalidMove := board.spots[5][6]       // f6

	t.Run("valid queen diagonal move", func(t *testing.T) {
		if !from.piece.validMove(from, validMoveDiagonal, &board) {
			fromStr, _ := from.string()
			destStr, _ := validMoveDiagonal.string()
			t.Errorf("Expected valid diagonal move from %v to %v", fromStr, destStr)
		}
	})

	t.Run("valid queen straight move", func(t *testing.T) {
		if !from.piece.validMove(from, validMoveStraight, &board) {
			fromStr, _ := from.string()
			destStr, _ := validMoveStraight.string()
			t.Errorf("Expected valid straight move from %v to %v", fromStr, destStr)
		}
	})

	t.Run("invalid queen move", func(t *testing.T) {
		if from.piece.validMove(from, invalidMove, &board) {
			fromStr, _ := from.string()
			destStr, _ := invalidMove.string()
			t.Errorf("Expected invalid move from %v to %v", fromStr, destStr)
		}
	})

	// Insert same color piece to test occupation and obstruction
	board.spots[5][5].piece = &bishop{
		color: "white",
	}

	t.Run("queen move obstruced", func(t *testing.T) {
		if from.piece.validMove(from, validMoveDiagonal, &board) {
			fromStr, _ := from.string()
			destStr, _ := invalidMove.string()
			t.Errorf("Expected invalid move from %v to %v due to obstruction by same color", fromStr, destStr)
		}
	})

	invalidMoveOccupied := board.spots[5][5]
	t.Run("queen move to occupied by same color", func(t *testing.T) {
		if from.piece.validMove(from, invalidMoveOccupied, &board) {
			fromStr, _ := from.string()
			destStr, _ := invalidMoveOccupied.string()
			t.Errorf("Expected invalid move from %v to %v due to occupation by same color", fromStr, destStr)
		}
	})

	// Insert opponent piece to test exposing king and capture
	board.spots[2][2].piece = &bishop{
		color: "black",
	}
	board.spots[5][5].piece = &king{
		color: "white",
	}
	board.whiteKingPosition = board.spots[5][5]
	toExpose := board.spots[4][6]
	t.Run("queen move exposes king", func(t *testing.T) {
		if from.piece.move(from, toExpose, &board) == nil {
			fromStr, _ := from.string()
			destStr, _ := toExpose.string()
			t.Errorf("Expected move to expose king from %v to %v", fromStr, destStr)
		}
	})

	t.Run("queen capture move", func(t *testing.T) {
		if from.piece.move(from, board.spots[2][2], &board) != nil {
			fromStr, _ := from.string()
			destStr, _ := board.spots[2][2].string()
			t.Errorf("Expected valid capture move from %v to %v", fromStr, destStr)
		}
	})
}

func TestRookMoves(t *testing.T) {
	board := board{}

	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	from := board.spots[4][4] // e5
	board.spots[from.rank][from.file].piece = &rook{
		color: "white",
	}

	validMove := board.spots[4][7]   // h5
	validMove2 := board.spots[1][4]  // e2
	invalidMove := board.spots[5][5] // f6

	t.Run("valid rook horizontal move", func(t *testing.T) {
		if !from.piece.validMove(from, validMove, &board) {
			fromStr, _ := from.string()
			destStr, _ := validMove.string()
			t.Errorf("Expected valid horizontal move from %v to %v", fromStr, destStr)
		}
	})

	t.Run("valid rook vertical move", func(t *testing.T) {
		if !from.piece.validMove(from, validMove2, &board) {
			fromStr, _ := from.string()
			destStr, _ := validMove2.string()
			t.Errorf("Expected valid vertical move from %v to %v", fromStr, destStr)
		}
	})

	t.Run("invalid rook move", func(t *testing.T) {
		if from.piece.validMove(from, invalidMove, &board) {
			fromStr, _ := from.string()
			destStr, _ := invalidMove.string()
			t.Errorf("Expected invalid move from %v to %v", fromStr, destStr)
		}
	})

	// Insert same color piece to test occupation and obstruction
	board.spots[4][5].piece = &bishop{
		color: "white",
	}

	t.Run("rook move obstruced", func(t *testing.T) {
		if from.piece.validMove(from, validMove, &board) {
			fromStr, _ := from.string()
			destStr, _ := invalidMove.string()
			t.Errorf("Expected invalid move from %v to %v due to obstruction by same color", fromStr, destStr)
		}
	})

	invalidMoveOccupied := board.spots[4][5]
	t.Run("rook move to occupied by same color", func(t *testing.T) {
		if from.piece.validMove(from, invalidMoveOccupied, &board) {
			fromStr, _ := from.string()
			destStr, _ := invalidMoveOccupied.string()
			t.Errorf("Expected invalid move from %v to %v due to occupation by same color", fromStr, destStr)
		}
	})

	// Insert opponent piece to test exposing king and capture
	board.spots[4][2].piece = &rook{
		color: "black",
	}
	board.spots[4][5].piece = &king{
		color: "white",
	}
	board.whiteKingPosition = board.spots[4][5]
	toExpose := board.spots[3][4]

	t.Run("rook move exposes king", func(t *testing.T) {
		if from.piece.move(from, toExpose, &board) == nil {
			fromStr, _ := from.string()
			destStr, _ := toExpose.string()
			t.Errorf("Expected move to expose king from %v to %v", fromStr, destStr)
		}
	})

	t.Run("rook capture move", func(t *testing.T) {
		if from.piece.move(from, board.spots[4][2], &board) != nil {
			fromStr, _ := from.string()
			destStr, _ := board.spots[4][2].string()
			t.Errorf("Expected valid capture move from %v to %v", fromStr, destStr)
		}
	})
}

func TestPawnMove(t *testing.T) {
	board := board{}

	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	from := board.spots[1][4] // e2
	board.spots[from.rank][from.file].piece = &pawn{
		color:     "white",
		direction: 1,
		hasMoved:  false,
	}

	validMove := board.spots[3][4]    // e4
	validMove2 := board.spots[2][4]   // e3
	invalidMove := board.spots[4][4]  // e5
	invalidMove2 := board.spots[2][5] // f3

	t.Run("valid pawn two-square move", func(t *testing.T) {
		if !from.piece.validMove(from, validMove, &board) {
			fromStr, _ := from.string()
			destStr, _ := validMove.string()
			t.Errorf("Expected valid two-square move from %v to %v", fromStr, destStr)
		}
	})

	t.Run("valid pawn one-square move", func(t *testing.T) {
		if !from.piece.validMove(from, validMove2, &board) {
			fromStr, _ := from.string()
			destStr, _ := validMove2.string()
			t.Errorf("Expected valid one-square move from %v to %v", fromStr, destStr)
		}
	})

	t.Run("invalid pawn three-square move", func(t *testing.T) {
		if from.piece.validMove(from, invalidMove, &board) {
			fromStr, _ := from.string()
			destStr, _ := invalidMove.string()
			t.Errorf("Expected invalid three-square move from %v to %v", fromStr, destStr)
		}
	})

	t.Run("invalid pawn diagonal move without capture", func(t *testing.T) {
		if from.piece.validMove(from, invalidMove2, &board) {
			fromStr, _ := from.string()
			destStr, _ := invalidMove2.string()
			t.Errorf("Expected invalid diagonal move without capture from %v to %v", fromStr, destStr)
		}
	})

	// Insert opponent piece to test capture and exposing
	board.spots[2][5].piece = &bishop{
		color: "black",
	}
	board.spots[0][3].piece = &king{
		color: "white",
	}
	board.whiteKingPosition = board.spots[0][3]

	t.Run("pawn move exposes king", func(t *testing.T) {
		if from.piece.move(from, validMove, &board) == nil {
			fromStr, _ := from.string()
			destStr, _ := validMove.string()
			t.Errorf("Expected move to expose king from %v to %v", fromStr, destStr)
		}
	})

	t.Run("valid pawn capture move", func(t *testing.T) {
		if !from.piece.validMove(from, invalidMove2, &board) {
			fromStr, _ := from.string()
			destStr, _ := invalidMove2.string()
			t.Errorf("Expected valid capture move from %v to %v", fromStr, destStr)
		}
	})

	// New scenario to test en passant
	board.spots[2][5].piece = nil

	newFrom := board.spots[4][3]          // d5
	enPassant := board.spots[4][4]        // e5
	enPassantCapture := board.spots[5][4] // e6
	board.spots[newFrom.rank][newFrom.file].piece = &pawn{
		color:     "white",
		direction: 1,
		hasMoved:  true,
	}
	board.spots[enPassant.rank][enPassant.file].piece = &pawn{
		color:     "black",
		direction: -1,
		hasMoved:  true,
	}
	board.enPassantTarget = enPassant

	t.Run("valid en passant move", func(t *testing.T) {
		if !newFrom.piece.validMove(newFrom, enPassantCapture, &board) {
			fromStr, _ := newFrom.string()
			destStr, _ := enPassantCapture.string()
			t.Errorf("Expected valid en passant move from %v to %v", fromStr, destStr)
		}
	})
}
