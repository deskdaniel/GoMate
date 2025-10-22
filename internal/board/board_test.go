package board

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/dragoo23/Go-chess/internal/app"
)

func TestInsufficientMaterialDraw(t *testing.T) {
	board := Board{}

	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &Position{
				Rank:  i,
				File:  j,
				Piece: nil,
			}
		}
	}

	// Case 1: king vs king
	board.spots[0][1].Piece = &King{
		color: "white",
	}
	board.spots[6][2].Piece = &King{
		color: "black",
	}
	t.Run("Only kings case", func(t *testing.T) {
		if haveSufficientMaterial(&board) {
			t.Error("Expecded draw by insufficient material (only kings)")
		}
	})

	// Case 2: king vs king + bishop
	board.spots[2][3].Piece = &Bishop{
		color: "white",
	}

	t.Run("Kings + single bishop case", func(t *testing.T) {
		if haveSufficientMaterial(&board) {
			t.Error("Expected draw by insufficient material (kings + single bishop)")
		}
	})

	// Case 3: king vs king + knight
	board.spots[2][3].Piece = &Knight{
		color: "white",
	}

	t.Run("Kings + single knight case", func(t *testing.T) {
		if haveSufficientMaterial(&board) {
			t.Error("Expected draw by insufficient material (kings + single knight)")
		}
	})

	// Case 4: king vs king + 2 bishops on same color squares
	board.spots[2][3].Piece = &Bishop{
		color: "white",
	}

	board.spots[4][5].Piece = &Bishop{
		color: "white",
	}

	t.Run("Kings + 2 bishops on same color", func(t *testing.T) {
		if haveSufficientMaterial(&board) {
			t.Error("Expected draw by insufficient material (kings + 2 bishops on same color squares)")
		}
	})

	// Case 5: king vs king + 2 bishops on different color squares
	board.spots[4][5].Piece = nil
	board.spots[5][5].Piece = &Bishop{
		color: "white",
	}

	t.Run("Kings + 2 bishops on different color", func(t *testing.T) {
		if !haveSufficientMaterial(&board) {
			t.Error("Expected game to have sufficient material for checkmate (bishops on different colored squares)")
		}
	})

	// Case 6: king vs king + rook
	board.spots[5][5].Piece = nil
	board.spots[2][3].Piece = &Rook{
		color: "white",
	}

	t.Run("Kings + rook", func(t *testing.T) {
		if !haveSufficientMaterial(&board) {
			t.Error("Expected game to have sufficient material for checkmate (rook)")
		}
	})

	// Case 7: king vs king + pawn
	board.spots[2][3].Piece = &Pawn{
		color:     "white",
		direction: 1,
	}

	t.Run("Kings + pawn", func(t *testing.T) {
		if !haveSufficientMaterial(&board) {
			t.Error("Expected game to have sufficient material for checkmate (pawn)")
		}
	})

	// Case 8: king vs king + queen
	board.spots[2][3].Piece = &Queen{
		color: "white",
	}

	t.Run("Kings + queen", func(t *testing.T) {
		if !haveSufficientMaterial(&board) {
			t.Error("Expected game to have sufficient material for checkmate (queen)")
		}
	})

	// Case 9: king vs king + 2 knights
	board.spots[5][5].Piece = &Knight{
		color: "white",
	}
	board.spots[2][3].Piece = &Knight{
		color: "white",
	}

	t.Run("Kings + 2 knights", func(t *testing.T) {
		if !haveSufficientMaterial(&board) {
			t.Error("Expected game to have sufficient material for checkmate (2 knights)")
		}
	})
}

func TestCheckmate(t *testing.T) {
	board := Board{}

	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &Position{
				Rank:  i,
				File:  j,
				Piece: nil,
			}
		}
	}

	// Case 1: simple checkmate
	board.spots[7][4].Piece = &King{
		color: "black",
	}
	board.blackKingPosition = board.spots[7][4]
	board.spots[6][4].Piece = &Queen{
		color: "white",
	}
	board.spots[5][4].Piece = &King{
		color: "white",
	}
	board.whiteKingPosition = board.spots[5][4]

	t.Run("Test checkmate", func(t *testing.T) {
		if hasLegalMove(&board, "black") {
			t.Error("Expected game to end with checkmate")
		}
	})

	// Case 2: simple check
	board.spots[5][4].Piece = nil
	board.spots[4][4].Piece = &King{
		color: "white",
	}
	board.whiteKingPosition = board.spots[4][4]

	t.Run("Test checkmate", func(t *testing.T) {
		if !hasLegalMove(&board, "black") {
			t.Error("Expected game not to end with checkmate")
		}
	})

	// Case 3: checkmate with pin
	board.spots[6][3].Piece = &Queen{
		color: "black",
	}
	board.spots[7][3].Piece = &Rook{
		color: "black",
	}
	board.spots[7][5].Piece = &Rook{
		color: "black",
	}
	board.spots[6][5].Piece = &Pawn{
		color:     "black",
		direction: -1,
	}

	board.spots[6][4].Piece = nil
	board.spots[4][1].Piece = &Bishop{
		color: "white",
	}
	board.spots[4][4].Piece = &Queen{
		color: "white",
	}
	board.spots[3][4].Piece = &King{
		color: "white",
	}
	board.whiteKingPosition = board.spots[3][4]

	t.Run("Test checkmate with pin", func(t *testing.T) {
		if hasLegalMove(&board, "black") {
			t.Error("Expected game to end with checkmate")
		}
	})
}

func TestIsUnderAttack(t *testing.T) {
	board := Board{}

	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &Position{
				Rank:  i,
				File:  j,
				Piece: nil,
			}
		}
	}

	// Case 1: simple check
	board.spots[7][4].Piece = &King{
		color: "black",
	}
	board.blackKingPosition = board.spots[7][4]

	board.spots[5][4].Piece = &Queen{
		color: "white",
	}

	t.Run("Check by queen", func(t *testing.T) {
		if !isUnderAttack(board.blackKingPosition, "black", &board) {
			t.Error("Expected black king to be under check by queen")
		}
	})

	// Case 2: not a check (blocked by piece)
	board.spots[6][4].Piece = &Pawn{
		color:     "black",
		direction: -1,
	}

	t.Run("Not a check", func(t *testing.T) {
		if isUnderAttack(board.blackKingPosition, "black", &board) {
			t.Error("Expected black king not to be under check")
		}
	})

	// Case 3: attacked by pawn
	board.spots[6][5].Piece = &Pawn{
		color:     "white",
		direction: 1,
	}

	t.Run("Check by pawn", func(t *testing.T) {
		if !isUnderAttack(board.blackKingPosition, "black", &board) {
			t.Error("Expected black king to be under attack by pawn")
		}
	})
}

func TestStalemate(t *testing.T) {
	model := boardModel{}

	board := Board{}

	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &Position{
				Rank:  i,
				File:  j,
				Piece: nil,
			}
		}
	}

	board.spots[7][0].Piece = &King{
		color: "black",
	}
	board.blackKingPosition = board.spots[7][0]

	board.spots[6][2].Piece = &King{
		color: "white",
	}
	board.whiteKingPosition = board.spots[6][2]

	board.spots[5][1].Piece = &Queen{
		color: "white",
	}

	model.board = &board
	model.whiteTurn = false

	t.Run("Stalemate check", func(t *testing.T) {
		if !stalemateCheck(&model) {
			t.Error("Expected stalemate")
		}
	})

	// Reset board for case 2
	board2 := Board{}
	for i := range board2.spots {
		for j := range board2.spots[i] {
			board2.spots[i][j] = &Position{
				Rank:  i,
				File:  j,
				Piece: nil,
			}
		}
	}

	board2.spots[7][0].Piece = &King{
		color: "black",
	}
	board2.blackKingPosition = board2.spots[7][0]
	board2.spots[6][1].Piece = &Rook{
		color: "black",
	}

	board2.spots[5][0].Piece = &King{
		color: "white",
	}
	board2.whiteKingPosition = board2.spots[5][0]
	board2.spots[5][2].Piece = &Bishop{
		color: "white",
	}
	board2.spots[5][3].Piece = &Bishop{
		color: "white",
	}

	model.board = &board2
	model.check = ""

	t.Run("Stalemate check with pin", func(t *testing.T) {
		if !stalemateCheck(&model) {
			t.Errorf("Expected stalemate")
		}
	})
}

func Test50MoveRule(t *testing.T) {
	tests := []struct {
		testName    string
		staleTurns  int
		wantDraw    bool
		wantWarning bool
	}{
		{"Low stale turns count, not a warning", 10, false, false},
		{"Medium stale turns count,  warning", 70, false, true},
		{"Reach 50 move rule treshold, draw", 100, true, false},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			board := Board{}
			board.staleTurns = test.staleTurns
			draw, warning := check50MoveFule(board.staleTurns)
			if (draw != test.wantDraw) || (warning != test.wantWarning) {
				t.Errorf("Expected draw to be %v got %v. Expected warning to be %v got %v", test.wantDraw, draw, test.wantWarning, warning)
			}
		})
	}

	board := Board{}
	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &Position{
				Rank:  i,
				File:  j,
				Piece: nil,
			}
		}
	}

	tests2 := []struct {
		testName      string
		kingPos       string
		whitePiecePos string
		blackPiecePos string
		movePos       string
		staleTurns    int
		pieceCreation func(color string) piece
	}{
		{"Test increase/reset counter logic for bishop", "a1", "b2", "d4", "c3", 30, func(color string) piece { return &Bishop{color: color} }},
		{"Test increase/reset counter logic for knight", "a1", "b1", "d1", "c3", 20, func(color string) piece { return &Knight{color: color} }},
		{"Test increase/reset counter logic for queen", "a1", "b2", "d4", "c3", 40, func(color string) piece { return &Queen{color: color} }},
		{"Test increase/reset counter logic for rook", "a1", "b1", "b5", "b2", 10, func(color string) piece { return &Rook{color: color} }},
	}

	for _, test := range tests2 {
		t.Run(test.testName, func(t *testing.T) {
			model := boardModel{}
			board := Board{}
			for i := range board.spots {
				for j := range board.spots[i] {
					board.spots[i][j] = &Position{
						Rank:  i,
						File:  j,
						Piece: nil,
					}
				}
			}
			model.board = &board
			kingPos, _ := PositionFromString(test.kingPos, &model)
			board.spots[kingPos.Rank][kingPos.File].Piece = &King{
				color: "white",
			}
			board.whiteKingPosition = kingPos

			whitePos, _ := PositionFromString(test.whitePiecePos, &model)
			whiteSquare := board.spots[whitePos.Rank][whitePos.File]
			whiteSquare.Piece = test.pieceCreation("white")

			blackPos, _ := PositionFromString(test.blackPiecePos, &model)
			blackSquare := board.spots[blackPos.Rank][blackPos.File]
			blackSquare.Piece = test.pieceCreation("black")

			board.staleTurns = test.staleTurns

			movePos, _ := PositionFromString(test.movePos, &model)
			moveSquare := board.spots[movePos.Rank][movePos.File]

			err := whiteSquare.Piece.Move(whiteSquare, moveSquare, &board)
			if err != nil {
				t.Error("Expected legal move", err)
			}
			if board.staleTurns != (test.staleTurns + 1) {
				t.Errorf("Expected stale turns counter to increase to %d, got %d instead", test.staleTurns+1, board.staleTurns)
			}

			err = moveSquare.Piece.Move(moveSquare, blackSquare, &board)
			if err != nil {
				t.Error("Expected legal move", err)
			}
			if board.staleTurns != 0 {
				t.Errorf("Expected stale turns counter to reset, got %d instead", board.staleTurns)
			}
		})
	}

	// Test increase/reset counter logic for king
	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &Position{
				Rank:  i,
				File:  j,
				Piece: nil,
			}
		}
	}

	board.spots[0][0].Piece = &King{
		color: "white",
	}
	board.whiteKingPosition = board.spots[0][0]
	board.spots[1][2].Piece = &Pawn{
		color:     "black",
		direction: -1,
	}

	board.staleTurns = 20
	err := board.spots[0][0].Piece.Move(board.spots[0][0], board.spots[1][1], &board)
	if err != nil {
		t.Error("Expected legal move1 for king")
	}
	if board.staleTurns != 21 {
		t.Errorf("Expected stale turns counter to increase to 21, got %d instead", board.staleTurns)
	}
	err = board.spots[1][1].Piece.Move(board.spots[1][1], board.spots[1][2], &board)
	if err != nil {
		t.Error("Expected legal move2 for king")
	}
	if board.staleTurns != 0 {
		t.Errorf("Expected stale turns counter to reset, got %d instead", board.staleTurns)
	}

	// Test reset counter for pawn move
	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &Position{
				Rank:  i,
				File:  j,
				Piece: nil,
			}
		}
	}

	board.spots[0][0].Piece = &King{
		color: "white",
	}
	board.whiteKingPosition = board.spots[0][0]
	board.spots[2][1].Piece = &Pawn{
		color:     "white",
		direction: 1,
	}

	board.staleTurns = 40
	err = board.spots[2][1].Piece.Move(board.spots[2][1], board.spots[3][1], &board)
	if err != nil {
		t.Error("Expected legal move2 for pawn")
	}
	if board.staleTurns != 0 {
		t.Errorf("Expected stale turns counter to reset, got %d instead", board.staleTurns)
	}
}

func TestDrawOffer(t *testing.T) {
	ctx := app.Context{}
	model := &boardModel{}
	board := &Board{}
	for i := range board.spots {
		for j := range board.spots[i] {
			board.spots[i][j] = &Position{
				Rank:  i,
				File:  j,
				Piece: nil,
			}
		}
	}
	model.board = board
	model.input = textinput.New()
	model.input.Placeholder = ""
	model.input.Prompt = ""
	model.ctx = &ctx

	msg := gameMsg{input: "draw"}
	model.Update(msg)
	if !model.offeredDraw {
		t.Error("Expected offeredDraw flag to be raised")
	}

	_, cmd := model.Update(msg)
	if cmd == nil {
		t.Error("Expected command signaling end of the game, got nothing instead")
	}
	msgOut := cmd()
	over, ok := msgOut.(overMsg)
	if !ok {
		t.Errorf("Expected overMsg, got %T", msgOut)
	}
	if !over.draw {
		t.Error("Expected draw flag to be set to true")
	}
	if over.message != "Game ended in a draw by agreement." {
		t.Errorf("Unexpected game over message: %s", over.message)
	}

	model.offeredDraw = true
	msg = gameMsg{input: "a2 a3"}
	model.Update(msg)
	if model.offeredDraw {
		t.Error("Expected draw to be refused and flag to be cleared")
	}
}
