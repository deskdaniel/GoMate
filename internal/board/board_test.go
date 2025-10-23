package board

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/dragoo23/Go-chess/internal/app"
)

func TestInsufficientMaterialDraw(t *testing.T) {
	b := board{}

	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	// Case 1: king vs king
	b.spots[0][1].piece = &king{
		color: "white",
	}
	b.spots[6][2].piece = &king{
		color: "black",
	}
	t.Run("Only kings case", func(t *testing.T) {
		if haveSufficientMaterial(&b) {
			t.Error("Expecded draw by insufficient material (only kings)")
		}
	})

	// Case 2: king vs king + bishop
	b.spots[2][3].piece = &bishop{
		color: "white",
	}

	t.Run("Kings + single bishop case", func(t *testing.T) {
		if haveSufficientMaterial(&b) {
			t.Error("Expected draw by insufficient material (kings + single bishop)")
		}
	})

	// Case 3: king vs king + knight
	b.spots[2][3].piece = &knight{
		color: "white",
	}

	t.Run("Kings + single knight case", func(t *testing.T) {
		if haveSufficientMaterial(&b) {
			t.Error("Expected draw by insufficient material (kings + single knight)")
		}
	})

	// Case 4: king vs king + 2 bishops on same color squares
	b.spots[2][3].piece = &bishop{
		color: "white",
	}

	b.spots[4][5].piece = &bishop{
		color: "white",
	}

	t.Run("Kings + 2 bishops on same color", func(t *testing.T) {
		if haveSufficientMaterial(&b) {
			t.Error("Expected draw by insufficient material (kings + 2 bishops on same color squares)")
		}
	})

	// Case 5: king vs king + 2 bishops on different color squares
	b.spots[4][5].piece = nil
	b.spots[5][5].piece = &bishop{
		color: "white",
	}

	t.Run("Kings + 2 bishops on different color", func(t *testing.T) {
		if !haveSufficientMaterial(&b) {
			t.Error("Expected game to have sufficient material for checkmate (bishops on different colored squares)")
		}
	})

	// Case 6: king vs king + rook
	b.spots[5][5].piece = nil
	b.spots[2][3].piece = &rook{
		color: "white",
	}

	t.Run("Kings + rook", func(t *testing.T) {
		if !haveSufficientMaterial(&b) {
			t.Error("Expected game to have sufficient material for checkmate (rook)")
		}
	})

	// Case 7: king vs king + pawn
	b.spots[2][3].piece = &pawn{
		color:     "white",
		direction: 1,
	}

	t.Run("Kings + pawn", func(t *testing.T) {
		if !haveSufficientMaterial(&b) {
			t.Error("Expected game to have sufficient material for checkmate (pawn)")
		}
	})

	// Case 8: king vs king + queen
	b.spots[2][3].piece = &queen{
		color: "white",
	}

	t.Run("Kings + queen", func(t *testing.T) {
		if !haveSufficientMaterial(&b) {
			t.Error("Expected game to have sufficient material for checkmate (queen)")
		}
	})

	// Case 9: king vs king + 2 knights
	b.spots[5][5].piece = &knight{
		color: "white",
	}
	b.spots[2][3].piece = &knight{
		color: "white",
	}

	t.Run("Kings + 2 knights", func(t *testing.T) {
		if !haveSufficientMaterial(&b) {
			t.Error("Expected game to have sufficient material for checkmate (2 knights)")
		}
	})
}

func TestCheckmate(t *testing.T) {
	b := board{}

	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	// Case 1: simple checkmate
	b.spots[7][4].piece = &king{
		color: "black",
	}
	b.blackKingPosition = b.spots[7][4]
	b.spots[6][4].piece = &queen{
		color: "white",
	}
	b.spots[5][4].piece = &king{
		color: "white",
	}
	b.whiteKingPosition = b.spots[5][4]

	t.Run("Test checkmate", func(t *testing.T) {
		if hasLegalMove(&b, "black") {
			t.Error("Expected game to end with checkmate")
		}
	})

	// Case 2: simple check
	b.spots[5][4].piece = nil
	b.spots[4][4].piece = &king{
		color: "white",
	}
	b.whiteKingPosition = b.spots[4][4]

	t.Run("Test checkmate", func(t *testing.T) {
		if !hasLegalMove(&b, "black") {
			t.Error("Expected game not to end with checkmate")
		}
	})

	// Case 3: checkmate with pin
	b.spots[6][3].piece = &queen{
		color: "black",
	}
	b.spots[7][3].piece = &rook{
		color: "black",
	}
	b.spots[7][5].piece = &rook{
		color: "black",
	}
	b.spots[6][5].piece = &pawn{
		color:     "black",
		direction: -1,
	}

	b.spots[6][4].piece = nil
	b.spots[4][1].piece = &bishop{
		color: "white",
	}
	b.spots[4][4].piece = &queen{
		color: "white",
	}
	b.spots[3][4].piece = &king{
		color: "white",
	}
	b.whiteKingPosition = b.spots[3][4]

	t.Run("Test checkmate with pin", func(t *testing.T) {
		if hasLegalMove(&b, "black") {
			t.Error("Expected game to end with checkmate")
		}
	})
}

func TestIsUnderAttack(t *testing.T) {
	b := board{}

	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	// Case 1: simple check
	b.spots[7][4].piece = &king{
		color: "black",
	}
	b.blackKingPosition = b.spots[7][4]

	b.spots[5][4].piece = &queen{
		color: "white",
	}

	t.Run("Check by queen", func(t *testing.T) {
		if !isUnderAttack(b.blackKingPosition, "black", &b) {
			t.Error("Expected black king to be under check by queen")
		}
	})

	// Case 2: not a check (blocked by piece)
	b.spots[6][4].piece = &pawn{
		color:     "black",
		direction: -1,
	}

	t.Run("Not a check", func(t *testing.T) {
		if isUnderAttack(b.blackKingPosition, "black", &b) {
			t.Error("Expected black king not to be under check")
		}
	})

	// Case 3: attacked by pawn
	b.spots[6][5].piece = &pawn{
		color:     "white",
		direction: 1,
	}

	t.Run("Check by pawn", func(t *testing.T) {
		if !isUnderAttack(b.blackKingPosition, "black", &b) {
			t.Error("Expected black king to be under attack by pawn")
		}
	})
}

func TestStalemate(t *testing.T) {
	model := boardModel{}

	b := board{}

	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	b.spots[7][0].piece = &king{
		color: "black",
	}
	b.blackKingPosition = b.spots[7][0]

	b.spots[6][2].piece = &king{
		color: "white",
	}
	b.whiteKingPosition = b.spots[6][2]

	b.spots[5][1].piece = &queen{
		color: "white",
	}

	model.board = &b
	model.whiteTurn = false

	t.Run("Stalemate check", func(t *testing.T) {
		if !stalemateCheck(&model) {
			t.Error("Expected stalemate")
		}
	})

	// Reset board for case 2
	b2 := board{}
	for i := range b2.spots {
		for j := range b2.spots[i] {
			b2.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	b2.spots[7][0].piece = &king{
		color: "black",
	}
	b2.blackKingPosition = b2.spots[7][0]
	b2.spots[6][1].piece = &rook{
		color: "black",
	}

	b2.spots[5][0].piece = &king{
		color: "white",
	}
	b2.whiteKingPosition = b2.spots[5][0]
	b2.spots[5][2].piece = &bishop{
		color: "white",
	}
	b2.spots[5][3].piece = &bishop{
		color: "white",
	}

	model.board = &b2
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
			b := board{}
			b.staleTurns = test.staleTurns
			draw, warning := check50MoveFule(b.staleTurns)
			if (draw != test.wantDraw) || (warning != test.wantWarning) {
				t.Errorf("Expected draw to be %v got %v. Expected warning to be %v got %v", test.wantDraw, draw, test.wantWarning, warning)
			}
		})
	}

	b := board{}
	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
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
		{"Test increase/reset counter logic for bishop", "a1", "b2", "d4", "c3", 30, func(color string) piece { return &bishop{color: color} }},
		{"Test increase/reset counter logic for knight", "a1", "b1", "d1", "c3", 20, func(color string) piece { return &knight{color: color} }},
		{"Test increase/reset counter logic for queen", "a1", "b2", "d4", "c3", 40, func(color string) piece { return &queen{color: color} }},
		{"Test increase/reset counter logic for rook", "a1", "b1", "b5", "b2", 10, func(color string) piece { return &rook{color: color} }},
	}

	for _, test := range tests2 {
		t.Run(test.testName, func(t *testing.T) {
			model := boardModel{}
			b := board{}
			for i := range b.spots {
				for j := range b.spots[i] {
					b.spots[i][j] = &position{
						rank:  i,
						file:  j,
						piece: nil,
					}
				}
			}
			model.board = &b
			kingPos, _ := positionFromString(test.kingPos, &model)
			b.spots[kingPos.rank][kingPos.file].piece = &king{
				color: "white",
			}
			b.whiteKingPosition = kingPos

			whitePos, _ := positionFromString(test.whitePiecePos, &model)
			whiteSquare := b.spots[whitePos.rank][whitePos.file]
			whiteSquare.piece = test.pieceCreation("white")

			blackPos, _ := positionFromString(test.blackPiecePos, &model)
			blackSquare := b.spots[blackPos.rank][blackPos.file]
			blackSquare.piece = test.pieceCreation("black")

			b.staleTurns = test.staleTurns

			movePos, _ := positionFromString(test.movePos, &model)
			moveSquare := b.spots[movePos.rank][movePos.file]

			err := whiteSquare.piece.move(whiteSquare, moveSquare, &b)
			if err != nil {
				t.Error("Expected legal move", err)
			}
			if b.staleTurns != (test.staleTurns + 1) {
				t.Errorf("Expected stale turns counter to increase to %d, got %d instead", test.staleTurns+1, b.staleTurns)
			}

			err = moveSquare.piece.move(moveSquare, blackSquare, &b)
			if err != nil {
				t.Error("Expected legal move", err)
			}
			if b.staleTurns != 0 {
				t.Errorf("Expected stale turns counter to reset, got %d instead", b.staleTurns)
			}
		})
	}

	// Test increase/reset counter logic for king
	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	b.spots[0][0].piece = &king{
		color: "white",
	}
	b.whiteKingPosition = b.spots[0][0]
	b.spots[1][2].piece = &pawn{
		color:     "black",
		direction: -1,
	}

	b.staleTurns = 20
	err := b.spots[0][0].piece.move(b.spots[0][0], b.spots[1][1], &b)
	if err != nil {
		t.Error("Expected legal move1 for king")
	}
	if b.staleTurns != 21 {
		t.Errorf("Expected stale turns counter to increase to 21, got %d instead", b.staleTurns)
	}
	err = b.spots[1][1].piece.move(b.spots[1][1], b.spots[1][2], &b)
	if err != nil {
		t.Error("Expected legal move2 for king")
	}
	if b.staleTurns != 0 {
		t.Errorf("Expected stale turns counter to reset, got %d instead", b.staleTurns)
	}

	// Test reset counter for pawn move
	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	b.spots[0][0].piece = &king{
		color: "white",
	}
	b.whiteKingPosition = b.spots[0][0]
	b.spots[2][1].piece = &pawn{
		color:     "white",
		direction: 1,
	}

	b.staleTurns = 40
	err = b.spots[2][1].piece.move(b.spots[2][1], b.spots[3][1], &b)
	if err != nil {
		t.Error("Expected legal move2 for pawn")
	}
	if b.staleTurns != 0 {
		t.Errorf("Expected stale turns counter to reset, got %d instead", b.staleTurns)
	}
}

func TestDrawOffer(t *testing.T) {
	ctx := app.Context{}
	model := &boardModel{}
	b := &board{}
	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}
	model.board = b
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
		t.Fatal("Expected command signaling end of the game, got nothing instead")
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

func TestOverMsgCheckmate(t *testing.T) {
	// Case 1: proper checkmate
	ctx := app.Context{}
	model := &boardModel{}
	b := &board{}
	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	b.spots[7][4].piece = &king{
		color: "black",
	}
	b.blackKingPosition = b.spots[7][4]
	b.spots[6][1].piece = &queen{
		color: "white",
	}
	b.spots[5][4].piece = &king{
		color: "white",
	}
	b.whiteKingPosition = b.spots[5][4]

	model.whiteTurn = true
	model.board = b
	model.input = textinput.New()
	model.input.Placeholder = ""
	model.input.Prompt = ""
	model.ctx = &ctx

	msg := gameMsg{input: "b7 e7"}
	_, cmd := model.Update(msg)
	if cmd == nil {
		t.Fatal("Expected game over message (checkmate)")
	}
	msgOut := cmd()
	over, ok := msgOut.(overMsg)
	if !ok {
		t.Errorf("Expected overMsg, got %T", msgOut)
	}
	if over.draw {
		t.Error("Expected game to end in a win/loss, not in draw")
	}
	if over.message != "Black king is in checkmate! Game over." {
		t.Errorf("Expected black king to be in checkmate, got: %q", over.message)
	}

	// Case 2: only check
	model = &boardModel{}
	b = &board{}
	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	b.spots[7][4].piece = &king{
		color: "black",
	}
	b.blackKingPosition = b.spots[7][4]
	b.spots[6][1].piece = &queen{
		color: "white",
	}
	b.spots[4][4].piece = &king{
		color: "white",
	}
	b.whiteKingPosition = b.spots[4][4]

	model.whiteTurn = true
	model.board = b
	model.input = textinput.New()
	model.input.Placeholder = ""
	model.input.Prompt = ""
	model.ctx = &ctx

	msg = gameMsg{input: "b7 e7"}
	_, cmd = model.Update(msg)

	if cmd != nil {
		t.Error("Expected game to continue")
	}
}

func TestOverMsgStalemate(t *testing.T) {
	ctx := app.Context{}
	model := &boardModel{}
	b := &board{}
	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	b.spots[7][0].piece = &king{
		color: "black",
	}
	b.blackKingPosition = b.spots[7][0]

	b.spots[6][2].piece = &king{
		color: "white",
	}
	b.whiteKingPosition = b.spots[6][2]

	b.spots[3][1].piece = &queen{
		color: "white",
	}

	model.whiteTurn = true
	model.board = b
	model.input = textinput.New()
	model.input.Placeholder = ""
	model.input.Prompt = ""
	model.ctx = &ctx

	msg := gameMsg{input: "b4 b6"}
	_, cmd := model.Update(msg)
	if cmd == nil {
		t.Fatal("Expected game over message (stalemate)")
	}
	msgOut := cmd()
	over, ok := msgOut.(overMsg)
	if !ok {
		t.Errorf("Expected overMsg, got %T", msgOut)
	}
	if !over.draw {
		t.Error("Expected game to end in a draw")
	}
	if over.message != "Draw due to stalemate! Game over." {
		t.Errorf("Expected stalemate, got: %q", over.message)
	}
}

func TestOverMsgInsufficientMaterial(t *testing.T) {
	ctx := app.Context{}
	model := &boardModel{}
	b := &board{}
	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	b.spots[0][1].piece = &king{
		color: "white",
	}
	b.whiteKingPosition = b.spots[0][1]
	b.spots[6][2].piece = &king{
		color: "black",
	}
	b.blackKingPosition = b.spots[6][2]
	b.spots[0][2].piece = &rook{
		color: "black",
	}

	model.whiteTurn = true
	model.board = b
	model.input = textinput.New()
	model.input.Placeholder = ""
	model.input.Prompt = ""
	model.ctx = &ctx

	msg := gameMsg{input: "b1 c1"}
	_, cmd := model.Update(msg)
	if cmd == nil {
		t.Fatal("Expected game over message (insufficient material)")
	}
	msgOut := cmd()
	over, ok := msgOut.(overMsg)
	if !ok {
		t.Errorf("Expected overMsg, got %T", msgOut)
	}
	if !over.draw {
		t.Error("Expected game to end in a draw")
	}
	if over.message != "Draw due to insufficient material! Game over." {
		t.Errorf("Expected insufficient material, got: %q", over.message)
	}
}

func TestOverMsg50MoveRule(t *testing.T) {
	ctx := app.Context{}
	model := &boardModel{}
	b := &board{}
	for i := range b.spots {
		for j := range b.spots[i] {
			b.spots[i][j] = &position{
				rank:  i,
				file:  j,
				piece: nil,
			}
		}
	}

	b.spots[0][1].piece = &king{
		color: "white",
	}
	b.whiteKingPosition = b.spots[0][1]
	b.spots[6][2].piece = &king{
		color: "black",
	}
	b.blackKingPosition = b.spots[6][2]
	b.spots[0][5].piece = &rook{
		color: "black",
	}

	b.staleTurns = 99
	model.whiteTurn = true
	model.board = b
	model.input = textinput.New()
	model.input.Placeholder = ""
	model.input.Prompt = ""
	model.ctx = &ctx

	msg := gameMsg{input: "b1 b2"}
	_, cmd := model.Update(msg)
	if cmd == nil {
		t.Fatal("Expected game over message (50 move rule)")
	}
	msgOut := cmd()
	over, ok := msgOut.(overMsg)
	if !ok {
		t.Errorf("Expected overMsg, got %T", msgOut)
	}
	if !over.draw {
		t.Error("Expected game to end in a draw")
	}
	if over.message != "Draw due to fifty-move rule! Game over." {
		t.Errorf("Expected fifty-move rule draw, got: %q", over.message)
	}
}
