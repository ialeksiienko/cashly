package tokenservice

import (
	"cashly/internal/entity"
	"context"
	"log/slog"
)

type tokenSaver interface {
	Save(ctx context.Context, familyID int, userID int64, token string) (*entity.UserBankToken, error)
}

func (ts *TokenService) Save(ctx context.Context, familyID int, userID int64, token string) (*entity.UserBankToken, error) {
	ubt, err := ts.tokenSaver.Save(ctx, familyID, userID, token)
	if err != nil {
		ts.sl.Error("unable to save user bank token", slog.Int("family_id", familyID), slog.Int("user_id", int(userID)), slog.String("err", err.Error()))
		return nil, err
	}
	return ubt, nil
}
