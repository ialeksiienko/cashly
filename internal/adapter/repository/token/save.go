package token

import (
	"cashly/internal/entity"
	"context"
	"log/slog"
)

func (tr Repository) Save(ctx context.Context, familyID int, userID int64, token string) (*entity.UserBankToken, error) {
	q := `INSERT INTO user_bank_tokens (user_id, family_id, token)
        VALUES ($1, $2, $3) RETURNING id, user_id, family_id, token, created_at`

	ubt := new(entity.UserBankToken)

	err := tr.db.Pool().QueryRow(ctx, q, userID, familyID, token).Scan(&ubt.ID, &ubt.UserID, &ubt.FamilyID, &ubt.Token, &ubt.CreatedAt)
	if err != nil {
		tr.logger.Error("unable to save token", slog.String("err", err.Error()))
	}

	return ubt, err
}
