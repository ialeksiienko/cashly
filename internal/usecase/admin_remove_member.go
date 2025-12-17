package usecase

import (
	"cashly/internal/pkg/errorsx"
	"cashly/internal/validate"
	"context"
)

func (uc *UseCase) RemoveMember(ctx context.Context, familyID int, userID int64, memberID int64) error {
	family, err := uc.familyService.GetByID(ctx, familyID)
	if err != nil {
		return err
	}

	if !validate.AdminPermission(userID, family.CreatedBy) {
		return errorsx.New("no permission", errorsx.ErrCodeNoPermission, struct{}{})
	}

	if userID == memberID {
		return errorsx.New("cannot remove self", errorsx.ErrCodeCannotRemoveSelf, struct{}{})
	}

	return uc.adminService.DeleteFromFamily(ctx, family.ID, memberID)
}
