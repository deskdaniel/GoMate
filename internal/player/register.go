package player

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/dragoo23/Go-chess/internal/app"
	"github.com/dragoo23/Go-chess/internal/database"
	"github.com/google/uuid"
)

func RegisterPlayer(ctx *app.Context) error {
	userName := ctx.Username
	err := checkUsername(userName)
	if err != nil {
		return err
	}

	_, err = ctx.Queries.GetUserByName(context.Background(), userName)
	fmt.Print(err)
	if err == nil || err != sql.ErrNoRows {
		return fmt.Errorf("username already exists")
	}

	err = checkPassword(ctx.Password)
	if err != nil {
		return err
	}

	hashedPassword, err := HashPassword(ctx.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("failed to generate user ID: %w", err)
	}

	userParams := database.RegisterUserParams{
		ID:             id.String(),
		Username:       userName,
		CreatedAt:      sql.NullString{String: time.Now().Format(time.RFC3339), Valid: true},
		UpdatedAt:      sql.NullString{String: time.Now().Format(time.RFC3339), Valid: true},
		HashedPassword: hashedPassword,
	}

	_, err = ctx.Queries.RegisterUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

func checkUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if len(username) > 20 {
		return fmt.Errorf("username cannot exceed 20 characters")
	}
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	valid := regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username)
	if !valid {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}

	return nil
}

func checkPassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if len(password) < 9 {
		return fmt.Errorf("password must be at least 9 characters")
	}
	if len(password) > 50 {
		return fmt.Errorf("password cannot exceed 50 characters")
	}
	valid := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_\-+=\[{\]}\|\\:;"'<,>.?/]+$`).MatchString(password)
	if !valid {
		return fmt.Errorf("password contains invalid characters")
	}

	return nil
}
