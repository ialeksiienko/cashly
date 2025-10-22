package usecase

import (
	"cashly/internal/errorsx"
	"context"
	"time"
)

func (uc *UseCase) JoinFamily(ctx context.Context, code string, userID int64) (string, error) {
	family, expiresAt, err := uc.familyService.GetFamilyByCode(ctx, code)
	if err != nil {
		return "", err
	}

	if time.Now().After(expiresAt) {
		return "", errorsx.New("family invite code expired", errorsx.ErrCodeFamilyCodeExpired, expiresAt)
	}

	if err := uc.userService.SaveUserToFamily(ctx, family.ID, userID); err != nil {
		return "", err
	}

	return family.Name, nil
}
