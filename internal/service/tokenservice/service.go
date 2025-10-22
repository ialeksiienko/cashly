package tokenservice

import (
	"cashly/internal/entity"
	"cashly/internal/pkg/sl"
	"context"
)

type TokenServiceIface interface {
	Save(ctx context.Context, familyID int, userID int64, token string) (*entity.UserBankToken, error)
	Get(ctx context.Context, familyID int, userID int64) (*entity.UserBankToken, error)
	Delete(ctx context.Context, familyID int, userID int64) error
}

type encryptor interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(encrypted string) (string, error)
}

type TokenService struct {
	encryptor
	tokenSaver    tokenSaver
	tokenProvider tokenProvider
	tokenDeletor  tokenDeletor
	sl            sl.Logger
}

func New(
	key [32]byte,
	tokenIface TokenServiceIface,
	sl sl.Logger,
) *TokenService {
	return &TokenService{
		encryptor:     NewEncrypt(key, sl),
		tokenSaver:    tokenIface,
		tokenProvider: tokenIface,
		tokenDeletor:  tokenIface,
		sl:            sl,
	}
}
