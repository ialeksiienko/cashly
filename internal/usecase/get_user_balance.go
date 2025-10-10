package usecase

import (
	"context"
)

func (s *UseCase) GetBalance(ctx context.Context, familyID int, userID int64, cardType string, currency string) (float64, error) {
	return s.userService.GetBalance(ctx, familyID, userID, cardType, currency)
}
