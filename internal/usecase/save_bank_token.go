package usecase

import (
	"cashly/internal/entity"
	"context"
)

func (uc *UseCase) SaveBankToken(ctx context.Context, familyID int, userID int64, token string) (*entity.UserBankToken, error) {
	encryptedToken, err := uc.tokenService.Encrypt(token)
	if err != nil {
		return nil, err
	}

	ubt, err := uc.tokenService.Save(ctx, familyID, userID, encryptedToken)
	if err != nil {
		return nil, err
	}

	return ubt, nil
}
