package userservice

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v4"
)

type deletor interface {
	DeleteFromFamily(context.Context, pgx.Tx, int, int64) error
}

func (s Service) DeleteFromFamily(ctx context.Context, familyID int, userID int64) error {
	err := s.withTransaction.WithTransaction(ctx, func(tx pgx.Tx) error {
		return s.deletor.DeleteFromFamily(ctx, tx, familyID, userID)
	})
	if err != nil {
		s.logger.Error("failed to delete family", slog.Int("family_id", familyID), slog.String("error", err.Error()))
	}

	return err
}
