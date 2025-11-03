package player

import (
	"context"
	"database/sql"
	"testing"

	"github.com/deskdaniel/GoMate/internal/app"
	"github.com/deskdaniel/GoMate/internal/database"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"empty password", "", true},
		{"too short", "short", true},
		{"too long", "thispasswordiswaytoolongtobeacceptedbythesystemandshouldfail", true},
		{"invalid characters", "ValidLengthButInvalidCháracter", true},
		{"valid password", "ValidPass123!", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkPassword(test.password)
			if (err != nil) != test.wantErr {
				t.Errorf("checkPassword(%q) error = %v, wantErr %v", test.password, err, test.wantErr)
			}
		})
	}
}

func TestCheckUserName(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"empty username", "", true},
		{"too short", "ab", true},
		{"too long", "thisusernameiswaytoolongtobeaccepted", true},
		{"invalid characters", "Invalid*Name!", true},
		{"valid username", "Valid_Name123", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkUsername(test.username)
			if (err != nil) != test.wantErr {
				t.Errorf("checkUserName(%q) error = %v, wantErr %v", test.username, err, test.wantErr)
			}
		})
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	password := "SecurePass123!"
	hashed, err := hashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	err = checkPasswordHash(password, hashed)
	if err != nil {
		t.Errorf("CheckHash failed for correct password: %v", err)
	}

	err = checkPasswordHash("WrongPass!", hashed)
	if err == nil {
		t.Error("CheckHash did not fail for incorrect password")
	}
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatalf("Failed to set goose dialect: %v", err)
	}

	if err := goose.Up(db, "../../sql/schema"); err != nil {
		t.Fatalf("Failed to apply migrations: %v", err)
	}

	return db
}

func TestRegisterUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	queries := database.New(db)
	ctx := &app.Context{
		Queries:  queries,
		Username: "TestUser",
		Password: "TestPass123!",
	}

	err := registerPlayer(ctx)
	if err != nil {
		t.Fatalf("RegisterPlayer failed: %v", err)
	}

	user, err := queries.GetUserByName(context.Background(), "TestUser")
	if err != nil {
		t.Fatalf("GetUserByName failed: %v", err)
	}

	if user.Username != "TestUser" {
		t.Errorf("Expected username 'TestUser', got %q", user.Username)
	}

	if user.HashedPassword == "" {
		t.Error("Expected non-empty hashed password")
	}

	err = registerPlayer(ctx)
	if err == nil {
		t.Fatal("RegisterPlayer did not fail for duplicate username")
	}
}

func TestLoginUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	queries := database.New(db)
	ctx := &app.Context{
		Queries:  queries,
		Username: "LoginUser",
		Password: "LoginPass123!",
	}

	err := registerPlayer(ctx)
	if err != nil {
		t.Fatalf("RegisterPlayer failed: %v", err)
	}

	ctx.Password = "WrongPass!"
	err = loginPlayer(ctx, 1)
	if err == nil {
		t.Fatal("LoginPlayer did not fail for incorrect password")
	}

	ctx.Password = "LoginPass123!"
	err = loginPlayer(ctx, 1)
	if err != nil {
		t.Fatalf("LoginPlayer failed: %v", err)
	}

	if ctx.User1 == nil || ctx.User1.Username != "LoginUser" {
		t.Errorf("Expected User1 to be set with username 'LoginUser', got %+v", ctx.User1)
	}

	err = loginPlayer(ctx, 2)
	if err == nil {
		t.Fatalf("LoginPlayer did not fail for already logged in username")
	}

	ctx2 := &app.Context{
		Queries:  queries,
		Username: "LoginUser2",
		Password: "LoginPass123!",
	}
	err = registerPlayer(ctx2)
	if err != nil {
		t.Fatalf("RegisterPlayer for second user failed: %v", err)
	}

	err = loginPlayer(ctx2, 2)
	if err != nil {
		t.Fatalf("LoginPlayer for second user failed: %v", err)
	}

	if ctx2.User2 == nil || ctx2.User2.Username != "LoginUser2" {
		t.Errorf("Expected User2 to be set with username 'LoginUser2', got %+v", ctx2.User2)
	}
}

func TestUpdateUserRecord(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	queries := database.New(db)
	ctx := &app.Context{
		Queries:  queries,
		Username: "RecordUser",
		Password: "RecordPass123!",
	}

	err := registerPlayer(ctx)
	if err != nil {
		t.Fatalf("RegisterPlayer failed: %v", err)
	}

	user, err := queries.GetUserByName(context.Background(), "RecordUser")
	if err != nil {
		t.Fatalf("GetUserByName failed: %v", err)
	}

	err = updateUserRecord(user.Username, ctx, true, false, false)
	if err != nil {
		t.Fatalf("UpdateUserRecord failed: %v", err)
	}

	stats, err := checkStats(user.Username, ctx)
	if err != nil {
		t.Fatalf("CheckStats failed: %v", err)
	}

	if stats.Wins != 1 || stats.Losses != 0 || stats.Draws != 0 {
		t.Errorf("Unexpected stats after win: %+v", stats)
	}

	err = updateUserRecord(user.Username, ctx, false, true, false)
	if err != nil {
		t.Fatalf("UpdateUserRecord failed: %v", err)
	}

	stats, err = checkStats(user.Username, ctx)
	if err != nil {
		t.Fatalf("CheckStats failed: %v", err)
	}

	if stats.Wins != 1 || stats.Losses != 1 || stats.Draws != 0 {
		t.Errorf("Unexpected stats after loss: %+v", stats)
	}

	err = updateUserRecord(user.Username, ctx, false, false, true)
	if err != nil {
		t.Fatalf("UpdateUserRecord failed: %v", err)
	}

	stats, err = checkStats(user.Username, ctx)
	if err != nil {
		t.Fatalf("CheckStats failed: %v", err)
	}

	if stats.Wins != 1 || stats.Losses != 1 || stats.Draws != 1 {
		t.Errorf("Unexpected stats after draw: %+v", stats)
	}

	err = updateUserRecord(user.Username, ctx, true, true, false)
	if err == nil {
		t.Fatal("UpdateUserRecord did not fail for multiple outcomes")
	}

	err = updateUserRecord("NonExistentUser", ctx, true, false, false)
	if err == nil {
		t.Fatal("UpdateUserRecord did not fail for non-existent user")
	}

	err = updateUserRecord(user.Username, ctx, true, false, false)
	if err != nil {
		t.Fatalf("UpdateUserRecord failed: %v", err)
	}

	stats, err = checkStats(user.Username, ctx)
	if err != nil {
		t.Fatalf("CheckStats failed: %v", err)
	}

	if stats.Wins != 2 || stats.Losses != 1 || stats.Draws != 1 {
		t.Errorf("Unexpected stats after win: %+v", stats)
	}
}
