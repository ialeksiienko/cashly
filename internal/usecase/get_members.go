package usecase

import (
	"cashly/internal/entity"
	"cashly/internal/service/userservice"
	"context"
)

func (uc *UseCase) GetFamilyMembersInfo(ctx context.Context, family *entity.Family, userID int64) ([]userservice.MemberInfo, error) {
	return uc.userService.GetFamilyMembersInfo(ctx, family, userID)
}
