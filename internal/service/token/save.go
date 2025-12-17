package tokenservice

import (
	"cashly/internal/entity"
	"context"
	"log/slog"
)

type saver interface {
	Save(ctx context.Context, fid int, uid int64, token string) (*entity.UserBankToken, error)
}

func (s Service) Save(ctx context.Context, fid int, uid int64, token string) (*entity.UserBankToken, error) {
	l := s.logger.With(slog.Int("family_id", fid), slog.Int("user_id", int(uid)))

	encryptedToken, err := s.encryptor.Encrypt(token)
	if err != nil {
		return nil, err
	}

	ubt, err := s.saver.Save(ctx, fid, uid, encryptedToken)
	if err != nil {
		l.Error("unable to save user bank token", slog.String("err", err.Error()))
	}

	//decryptedToken, err := s.encryptor.Decrypt(ubt.Token)
	//if err != nil {
	//	l.Error("unable to decrypt token", slog.String("err", err.Error()))
	//	return nil, err
	//}
	//
	//ubt.Token = decryptedToken

	return ubt, err
}
