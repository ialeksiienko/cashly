package tokenrepo

import (
	"cashly/internal/entity"
	"context"
)

func (tr *TokenRepository) Save(ctx context.Context, familyID int, userID int64, token string) (*entity.UserBankToken, error) {
	q := `INSERT INTO user_bank_tokens (user_id, family_id, token)
        VALUES ($1, $2, $3) RETURNING id, user_id, family_id, token, created_at`

	ubt := new(entity.UserBankToken)

	err := tr.db.QueryRow(ctx, q, userID, familyID, token).Scan(&ubt.ID, &ubt.UserID, &ubt.FamilyID, &ubt.Token, &ubt.CreatedAt)
	if err != nil {
		tr.sl.Error(err.Error())
		return nil, err
	}

	return ubt, nil
}
