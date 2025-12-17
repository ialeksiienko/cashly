package usecase

import (
	"cashly/internal/entity"
	"cashly/internal/pkg/errorsx"
	"cashly/internal/validate"
	"context"
	"time"
)

func (uc *UseCase) CreateNewInviteCode(ctx context.Context, family *entity.Family, userID int64) (string, time.Time, error) {
	if !validate.AdminPermission(userID, family.CreatedBy) {
		return "", time.Time{}, errorsx.New("no permission", errorsx.ErrCodeNoPermission, struct{}{})
	}

	return uc.familyService.CreateNewInviteCode(ctx, family, userID)
}
