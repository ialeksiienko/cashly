package userrepo

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v4"
)

func (ur Repository) DeleteFromFamily(ctx context.Context, tx pgx.Tx, familyID int, userID int64) error {
	_, err := tx.Exec(ctx, `DELETE FROM user_bank_tokens WHERE user_id = $1 AND family_id = $2`, userID, familyID)
	if err != nil {
		ur.logger.Error("unable to delete token from user_bank_tokens", slog.Int("user_id", int(userID)), slog.Int("family_id", familyID), slog.String("err", err.Error()))
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM users_to_families WHERE user_id = $1 AND family_id = $2`, userID, familyID)
	if err != nil {
		ur.logger.Error("failed to leave family", slog.Int("user_id", int(userID)), slog.Int("family_id", familyID), slog.String("err", err.Error()))
	}

	return err
}
