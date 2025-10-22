package usecase

import (
	"cashly/internal/entity"
	"context"
)

func (s *UseCase) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	return s.userService.GetUserByID(ctx, id)
}
