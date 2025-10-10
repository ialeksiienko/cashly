package userrepo

import (
	"context"
	"log/slog"
	"monofamily/internal/adapter/database"
	"monofamily/internal/pkg/sl"

	"github.com/jackc/pgx/v4"
)

type databaseInface interface {
	database.PgxIface
}

type UserRepository struct {
	db databaseInface
	sl sl.Logger
}

func New(db databaseInface, sl sl.Logger) *UserRepository {
	return &UserRepository{
		db: db,
		sl: sl,
	}
}

func (ur *UserRepository) WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := ur.db.Begin(ctx)
	if err != nil {
		ur.sl.Error("unable to begin transaction", slog.String("error", err.Error()))
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	err = fn(tx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
