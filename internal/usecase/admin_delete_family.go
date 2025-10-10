package usecase

import (
	"context"
	"monofamily/internal/entity"
)

func (uc *UseCase) DeleteFamily(ctx context.Context, family *entity.Family, userID int64) error {
	if err := uc.checkAdminPermission(family.CreatedBy, userID); err != nil {
		return err
	}

	err := uc.familyService.Delete(ctx, family.ID)
	if err != nil {
		return err
	}

	return nil
}
