package usecase

import (
	"cashly/internal/entity"
	"cashly/internal/pkg/errorsx"
	"cashly/internal/validate"
	"context"
)

func (uc *UseCase) LeaveFamily(ctx context.Context, f *entity.Family, uid int64) error {
	if !validate.AdminPermission(uid, f.CreatedBy) {
		return errorsx.New("admin cannot leave family", errorsx.ErrCodeCannotRemoveSelf, struct{}{})
	}

	return uc.userService.DeleteFromFamily(ctx, f.ID, uid)
}
