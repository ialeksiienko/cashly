package usecase

import (
	"cashly/internal/entity"
	"cashly/internal/errorsx"
	"context"
	"time"
)

func (uc *UseCase) JoinFamily(ctx context.Context, code string, userID int64) (*entity.Family, error) {
	family, expiresAt, err := uc.familyService.GetFamilyByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if time.Now().After(expiresAt) {
		return nil, errorsx.New("family invite code expired", errorsx.ErrCodeFamilyCodeExpired, expiresAt)
	}

	if err := uc.userService.SaveUserToFamily(ctx, family.ID, userID); err != nil {
		return nil, err
	}

	return family, nil
}
