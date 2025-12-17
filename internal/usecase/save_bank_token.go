package usecase

import (
	"cashly/internal/entity"
	"context"
)

func (uc *UseCase) SaveBankToken(ctx context.Context, familyID int, userID int64, token string) (*entity.UserBankToken, error) {
	ubt, err := uc.tokenService.Save(ctx, familyID, userID, token)
	if err != nil {
		return nil, err
	}

	return ubt, nil
}
