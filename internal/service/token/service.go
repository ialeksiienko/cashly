package tokenservice

import (
	"cashly/internal/entity"
	"cashly/pkg/slogx"

	"context"
)

type TokenIface interface {
	Save(ctx context.Context, fid int, uid int64, token string) (*entity.UserBankToken, error)
	Get(ctx context.Context, fid int, uid int64) (*entity.UserBankToken, error)
	Delete(ctx context.Context, fid int, uid int64) error
}

type Encryptor interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(encrypted string) (string, error)
}

type Service struct {
	encryptor Encryptor
	saver     saver
	provider  provider
	deletor   deletor
	logger    slogx.Logger
}

func New(
	enc Encryptor,
	ti TokenIface,
	l slogx.Logger,
) Service {
	return Service{
		encryptor: enc,
		saver:     ti,
		provider:  ti,
		deletor:   ti,
		logger:    l,
	}
}
