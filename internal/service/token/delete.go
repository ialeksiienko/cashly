package tokenservice

import "context"

type deletor interface {
	Delete(ctx context.Context, familyID int, userID int64) error
}

func (s Service) Delete(ctx context.Context, familyID int, userID int64) error {
	return s.deletor.Delete(ctx, familyID, userID)
}
