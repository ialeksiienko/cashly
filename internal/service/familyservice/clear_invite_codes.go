package familyservice

import (
	"context"
	"log/slog"
)

type FamilyInviteCodeCleaner interface {
	ClearInviteCodes(ctx context.Context) error
}

func (s *FamilyService) ClearInviteCodes(ctx context.Context) error {
	err := s.familyInviteCodeCleaner.ClearInviteCodes(ctx)
	if err != nil {
		s.sl.Error("failed to clear invite codes", slog.String("error", err.Error()))
		return err
	}

	return nil
}
