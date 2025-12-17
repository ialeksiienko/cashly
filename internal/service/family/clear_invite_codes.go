package familyservice

import (
	"context"
	"log/slog"
)

type inviteCodeCleaner interface {
	ClearInviteCodes(ctx context.Context) error
}

func (s Service) ClearInviteCodes(ctx context.Context) error {
	err := s.inviteCodeCleaner.ClearInviteCodes(ctx)
	if err != nil {
		s.logger.Error("failed to clear invite codes", slog.String("error", err.Error()))
	}

	return err
}
