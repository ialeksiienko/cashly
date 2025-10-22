package usecase

import (
	"cashly/internal/errorsx"
	"context"
)

func (uc *UseCase) RemoveMember(ctx context.Context, familyID int, userID int64, memberID int64) error {
	family, err := uc.familyService.GetFamilyByID(ctx, familyID)
	if err != nil {
		return err
	}

	if err := uc.checkAdminPermission(family.CreatedBy, userID); err != nil {
		return err
	}

	if userID == memberID {
		return errorsx.New("cannot remove self", errorsx.ErrCodeCannotRemoveSelf, struct{}{})
	}

	return uc.adminService.DeleteUserFromFamily(ctx, family.ID, memberID)
}
