package userservice

import (
	"context"
	"log/slog"
)

func (s Service) SaveToFamily(ctx context.Context, familyID int, userID int64) error {
	saveErr := s.saver.SaveToFamily(ctx, familyID, userID)
	if saveErr != nil {
		s.logger.Error("unable to save user to family", slog.Int("user_id", int(userID)), slog.Int("family_id", familyID), slog.String("err", saveErr.Error()))
	}

	return saveErr
}
