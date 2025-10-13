package familyservice

import (
	"context"
	"log/slog"
	"monofamily/internal/entity"
	"monofamily/internal/errorsx"
	"time"
)

type familyInviteCodeSaver interface {
	SaveFamilyInviteCode(ctx context.Context, userId int64, familyId int, code string) (time.Time, error)
}

type familyInviteCodeProvider interface {
	GetInviteCode(ctx context.Context, familyID int) (string, time.Time, error)
}

func (s *FamilyService) CreateNewInviteCode(ctx context.Context, family *entity.Family, userID int64) (string, time.Time, error) {
	code, expiresAt, err := s.familyInviteCodeProvider.GetInviteCode(ctx, family.ID)
	if err != nil {
		s.sl.Error("unable to get invite code", slog.Int("family_id", family.ID), slog.String("err", err.Error()))
		return "", time.Time{}, err
	}

	if code == "" {
		code, err = s.GenerateInviteCode()
		if err != nil {
			s.sl.Error("failed to generate family invite code", slog.Int("family_id", family.ID), slog.String("err", err.Error()))
			return "", time.Time{}, errorsx.New("unable to generate invite code", errorsx.ErrCodeFailedToGenerateInviteCode, struct{}{})
		}

		expiresAt, err = s.familyInviteCodeSaver.SaveFamilyInviteCode(ctx, userID, family.ID, code)
		if err != nil {
			s.sl.Error("failed to save family invite code", slog.Int64("created_by", userID), slog.Int("family_id", family.ID), slog.String("code", code), slog.String("error", err.Error()))
			return "", time.Time{}, err
		}

	}

	return code, expiresAt, nil
}
