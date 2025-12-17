package usecase

import (
	"cashly/internal/entity"
	"cashly/internal/pkg/errorsx"
	"cashly/internal/validate"
	"context"
)

func (uc *UseCase) DeleteFamily(ctx context.Context, family *entity.Family, userID int64) error {
	if !validate.AdminPermission(userID, family.CreatedBy) {
		return errorsx.New("no permission", errorsx.ErrCodeNoPermission, struct{}{})
	}

	return uc.familyService.Delete(ctx, family.ID)
}
