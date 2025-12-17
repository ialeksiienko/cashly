package usecase

import (
	"context"
)

func (uc *UseCase) GetBalance(ctx context.Context, familyID int, checkedUserID int64, cardType string, currency string) (float64, error) {
	return uc.userService.GetBalance(ctx, familyID, checkedUserID, cardType, currency)
}
