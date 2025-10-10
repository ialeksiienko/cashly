package userservice

import (
	"context"
	"monofamily/internal/entity"
)

func (s *UserService) Register(ctx context.Context, user *entity.User) (*entity.User, error) {
	return s.userSaver.Save(ctx, user)
}
