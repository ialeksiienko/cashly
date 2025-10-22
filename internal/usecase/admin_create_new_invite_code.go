package usecase

import (
	"cashly/internal/entity"
	"context"
	"time"
)

func (uc *UseCase) CreateNewInviteCode(ctx context.Context, family *entity.Family, userID int64) (string, time.Time, error) {
	if err := uc.checkAdminPermission(family.CreatedBy, userID); err != nil {
		return "", time.Time{}, err
	}

	return uc.familyService.CreateNewInviteCode(ctx, family, userID)
}
