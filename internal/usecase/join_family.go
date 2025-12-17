package usecase

import (
	"cashly/internal/entity"
	"cashly/internal/pkg/errorsx"
	"context"
	"time"
)

func (uc *UseCase) JoinFamily(ctx context.Context, code string, uid int64) (*entity.Family, error) {
	f, expiresAt, err := uc.familyService.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if time.Now().After(expiresAt) {
		return nil, errorsx.New("family invite code expired", errorsx.ErrCodeFamilyCodeExpired, expiresAt)
	}

	err = uc.userService.SaveToFamily(ctx, f.ID, uid)

	return f, err
}
