package tokenservice

import "context"

type tokenDeletor interface {
	Delete(ctx context.Context, familyID int, userID int64) error
}

func (s *TokenService) Delete(ctx context.Context, familyID int, userID int64) error {
	return s.tokenDeletor.Delete(ctx, familyID, userID)
}
