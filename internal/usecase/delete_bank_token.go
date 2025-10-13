package usecase

import "context"

func (uc *UseCase) DeleteUserBankToken(ctx context.Context, familyID int, userID int64) error {
	return uc.tokenService.Delete(ctx, familyID, userID)
}
