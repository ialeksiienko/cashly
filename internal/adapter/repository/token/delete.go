package token

import (
	"context"
	"log/slog"
)

func (tr Repository) Delete(ctx context.Context, familyID int, userID int64) error {
	q := `DELETE FROM user_bank_tokens 
    WHERE user_id = $1 AND family_id = $2`

	_, err := tr.db.Pool().Exec(ctx, q, userID, familyID)
	if err != nil {
		tr.logger.Error("unable to delete token from user_bank_tokens", slog.Int("family_id", familyID), slog.Int("user_id", int(userID)), slog.Any("err", err))
	}

	return err
}
