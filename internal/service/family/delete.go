package familyservice

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v4"
)

type deletor interface {
	Delete(ctx context.Context, fn pgx.Tx, familyID int) error
}

func (s Service) Delete(ctx context.Context, fid int) error {
	err := s.withTransaction.WithTransaction(ctx, func(tx pgx.Tx) error {
		return s.deletor.Delete(ctx, tx, fid)
	})
	if err != nil {
		s.logger.Error("failed to delete family", slog.Int("family_id", fid), slog.String("error", err.Error()))
	}

	return err
}
