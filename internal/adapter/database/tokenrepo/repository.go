package tokenrepo

import (
	"cashly/internal/adapter/database"
	"cashly/internal/pkg/sl"
)

type databaseInface interface {
	database.PgxIface
}

type TokenRepository struct {
	db databaseInface
	sl sl.Logger
}

func New(db databaseInface, sl sl.Logger) *TokenRepository {
	return &TokenRepository{
		db: db,
		sl: sl,
	}
}
