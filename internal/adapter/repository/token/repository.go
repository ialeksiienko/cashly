package token

import (
	"cashly/internal/adapter/database"
	"cashly/pkg/slogx"
	"log/slog"
)

type Repository struct {
	db     database.Datastore
	logger slogx.Logger
}

func New(db database.Datastore, l slogx.Logger) Repository {
	return Repository{
		db:     db,
		logger: l.With(slog.String("repository", "token")),
	}
}
