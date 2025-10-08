package app

import (
	"github.com/dragoo23/Go-chess/internal/database"
)

type Context struct {
	Queries  *database.Queries
	Username string
	Password string
	User1    *User
	User2    *User
}

type User struct {
	ID       string
	Username string
	Slot     int
}
