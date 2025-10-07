package player

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dragoo23/Go-chess/internal/app"
	"github.com/dragoo23/Go-chess/internal/database"
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
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	err = CheckHash(password, hashed)
	if err != nil {
		t.Errorf("CheckHash failed for correct password: %v", err)
	}

	err = CheckHash("WrongPass!", hashed)
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

	err := RegisterPlayer(ctx)
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

	err = RegisterPlayer(ctx)
	if err == nil {
		t.Fatal("RegisterPlayer did not fail for duplicate username")
	}
}
