package usecase

import (
	"cashly/internal/entity"
	"cashly/internal/errorsx"
	"context"
)

func (uc *UseCase) LeaveFamily(ctx context.Context, family *entity.Family, userID int64) error {
	if family.CreatedBy == userID {
		return errorsx.New("admin cannot leave family", errorsx.ErrCodeCannotRemoveSelf, struct{}{})
	}

	err := uc.userService.DeleteUserFromFamily(ctx, family.ID, userID)
	if err != nil {
		return err
	}

	return nil
}
