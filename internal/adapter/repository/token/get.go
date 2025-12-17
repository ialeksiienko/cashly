package token

import (
	"cashly/internal/entity"
	"context"
)

func (tr Repository) Get(ctx context.Context, familyID int, userID int64) (*entity.UserBankToken, error) {
	q := `SELECT id, user_id, family_id, token, created_at
        FROM user_bank_tokens WHERE user_id = $1 AND family_id = $2`

	ubt := new(entity.UserBankToken)

	err := tr.db.Pool().QueryRow(ctx, q, userID, familyID).Scan(&ubt.ID, &ubt.UserID, &ubt.FamilyID, &ubt.Token, &ubt.CreatedAt)
	if err != nil {
		tr.logger.Error(err.Error())
	}

	return ubt, err
}
