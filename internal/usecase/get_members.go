package usecase

import (
	"cashly/internal/entity"
	userservice "cashly/internal/service/user"
	"context"
)

func (uc *UseCase) GetFamilyMembers(ctx context.Context, family *entity.Family, userID int64) ([]userservice.Member, error) {
	return uc.userService.GetFamilyMembers(ctx, family, userID)
}
