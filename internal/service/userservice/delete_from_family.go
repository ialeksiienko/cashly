package userservice

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v4"
)

type userDeletor interface {
	DeleteUserFromFamily(ctx context.Context, tx pgx.Tx, familyID int, userID int64) error
}

func (s *UserService) DeleteUserFromFamily(ctx context.Context, familyID int, userID int64) error {

	err := s.withTransaction.WithTransaction(ctx, func(tx pgx.Tx) error {
		return s.userDeletor.DeleteUserFromFamily(ctx, tx, familyID, userID)
	})
	if err != nil {
		s.sl.Error("failed to delete family", slog.Int("family_id", familyID), slog.String("error", err.Error()))
		return err
	}

	return nil
}
