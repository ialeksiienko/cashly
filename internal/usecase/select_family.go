package usecase

import (
	"context"
	"monofamily/internal/entity"
)

func (uc *UseCase) SelectFamily(ctx context.Context, familyID int, userID int64) (bool, bool, *entity.Family, error) {
	f, err := uc.familyService.GetFamilyByID(ctx, familyID)
	if err != nil {
		return false, false, nil, err
	}

	isAdmin := f.CreatedBy == userID

	hasToken, _, err := uc.tokenService.Get(ctx, familyID, userID)
	if err != nil {
		return false, false, nil, err
	}

	return isAdmin, hasToken, f, nil
}
