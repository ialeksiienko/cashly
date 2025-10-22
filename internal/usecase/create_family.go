package usecase

import (
	"cashly/internal/entity"
	"context"
	"time"
)

func (uc *UseCase) CreateFamily(ctx context.Context, familyName string, userID int64) (*entity.Family, string, time.Time, error) {
	family, err := uc.familyService.Create(ctx, familyName, userID)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	saveErr := uc.userService.SaveUserToFamily(ctx, family.ID, userID)
	if saveErr != nil {
		return nil, "", time.Time{}, saveErr
	}

	code, expiresAt, createErr := uc.familyService.CreateNewInviteCode(ctx, family, userID)
	if createErr != nil {
		return nil, "", time.Time{}, createErr
	}

	return family, code, expiresAt, nil
}
