package app

import (
	"github.com/dragoo23/Go-chess/internal/database"
)

type Context struct {
	Queries  *database.Queries
	Username string
	Password string
}
