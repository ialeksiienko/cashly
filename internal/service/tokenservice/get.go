package tokenservice

import (
	"context"
	"errors"
	"log/slog"
	"monofamily/internal/entity"

	"github.com/jackc/pgx/v4"
)

func (ts *TokenService) Get(ctx context.Context, familyID int, userID int64) (bool, *entity.UserBankToken, error) {

	logger := ts.sl.With(
		slog.Int("family_id", familyID),
		slog.Int64("user_id", userID),
	)

	ubt, err := ts.tokenProvider.Get(ctx, familyID, userID)
	if err != nil || ubt == nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Debug("token not found")
			return false, nil, nil
		} else {
			logger.Error("unable to get token from db", slog.String("err", err.Error()))
			return false, nil, err
		}
	}

	decryptedToken, err := ts.encryptor.Decrypt(ubt.Token)
	if err != nil {
		logger.Error("unable to decrypt token", slog.String("err", err.Error()))
		return false, nil, err
	}

	ubt.Token = decryptedToken

	return true, ubt, nil
}
