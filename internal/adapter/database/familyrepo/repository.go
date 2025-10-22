package familyrepo

import (
	"cashly/internal/adapter/database"
	"cashly/internal/pkg/sl"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v4"
)

type FamilyRepository struct {
	db database.PgxIface
	sl sl.Logger
}

func New(db database.PgxIface, sl sl.Logger) *FamilyRepository {
	return &FamilyRepository{
		db: db,
		sl: sl,
	}
}

func (fr *FamilyRepository) WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := fr.db.Begin(ctx)
	if err != nil {
		fr.sl.Error("unable to begin transaction", slog.String("error", err.Error()))
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
