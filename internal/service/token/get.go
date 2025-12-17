package tokenservice

import (
	"cashly/internal/entity"
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v4"
)

type provider interface {
	Get(ctx context.Context, fid int, uid int64) (*entity.UserBankToken, error)
}

func (s Service) Get(ctx context.Context, fid int, uid int64) (bool, *entity.UserBankToken, error) {

	logger := s.logger.With(
		slog.Int("family_id", fid),
		slog.Int64("user_id", uid),
	)

	ubt, err := s.provider.Get(ctx, fid, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Debug("token not found")
			return false, nil, nil
		}
		logger.Error("unable to get token from db", slog.String("err", err.Error()))
		return false, nil, err
	}

	decryptedToken, decErr := s.encryptor.Decrypt(ubt.Token)
	if decErr != nil {
		logger.Error("unable to decrypt token", slog.String("err", decErr.Error()))
		return false, nil, decErr
	}

	ubt.Token = decryptedToken

	return true, ubt, nil
}
