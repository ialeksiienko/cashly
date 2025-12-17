package familyservice

import (
	"cashly/internal/entity"
	"cashly/internal/pkg/errorsx"
	"context"
	"log/slog"
	"time"
)

type inviteCodeSaver interface {
	SaveInviteCode(ctx context.Context, uid int64, fid int, code string) (time.Time, error)
}

type inviteCodeProvider interface {
	GetInviteCode(ctx context.Context, fid int) (string, time.Time, error)
}

func (s Service) CreateNewInviteCode(ctx context.Context, f *entity.Family, uid int64) (string, time.Time, error) {
	code, expiresAt, err := s.inviteCodeProvider.GetInviteCode(ctx, f.ID)
	if err != nil {
		s.logger.Error("unable to get invite code", slog.Int("family_id", f.ID), slog.String("err", err.Error()))
		return "", time.Time{}, err
	}

	if code == "" {
		code, err = s.GenerateInviteCode()
		if err != nil {
			s.logger.Error("failed to generate family invite code", slog.Int("family_id", f.ID), slog.String("err", err.Error()))
			return "", time.Time{}, errorsx.New("unable to generate invite code", errorsx.ErrCodeFailedToGenerateInviteCode, struct{}{})
		}

		expiresAt, err = s.inviteCodeSaver.SaveInviteCode(ctx, uid, f.ID, code)
		if err != nil {
			s.logger.Error("failed to save family invite code", slog.Int64("created_by", uid), slog.Int("family_id", f.ID), slog.String("code", code), slog.String("error", err.Error()))
			return "", time.Time{}, err
		}
	}

	return code, expiresAt, nil
}
