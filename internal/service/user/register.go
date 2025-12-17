package userservice

import (
	"cashly/internal/entity"
	"context"
)

func (s Service) Register(ctx context.Context, user *entity.User) (*entity.User, error) {
	return s.saver.Save(ctx, user)
}
