package usecase

import (
	"context"
)

func (s *UseCase) GetBalance(ctx context.Context, familyID int, checkedUserID int64, cardType string, currency string) (float64, error) {
	return s.userService.GetBalance(ctx, familyID, checkedUserID, cardType, currency)
}
