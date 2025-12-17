package usecase

import (
	"cashly/internal/entity"
	"context"
)

func (uc *UseCase) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	return uc.userService.GetByID(ctx, id)
}
