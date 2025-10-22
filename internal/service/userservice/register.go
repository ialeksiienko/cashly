package userservice

import (
	"cashly/internal/entity"
	"context"
)

func (s *UserService) Register(ctx context.Context, user *entity.User) (*entity.User, error) {
	return s.userSaver.Save(ctx, user)
}
