package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dragoo23/Go-chess/internal/app"

	"github.com/dragoo23/Go-chess/internal/database"
	"github.com/dragoo23/Go-chess/internal/navigation"
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
	ctx := &app.Context{
		Queries: queries,
	}

	m := navigation.SetupNavigation(ctx)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
