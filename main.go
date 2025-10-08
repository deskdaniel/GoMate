package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dragoo23/Go-chess/internal/app"

	// "github.com/dragoo23/Go-chess/internal/board"
	"github.com/dragoo23/Go-chess/internal/database"
	// "github.com/dragoo23/Go-chess/internal/display"
	"github.com/dragoo23/Go-chess/internal/player"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := database.OpenDb()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	queries := database.New(db)
	context := &app.Context{
		Queries: queries,
	}

	// gameboard := board.InitializeBoard()
	m := player.SetupRegister(context)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
