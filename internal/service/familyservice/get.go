package familyservice

import (
	"cashly/internal/entity"
	"cashly/internal/errorsx"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v4"
)

type FamilyProvider interface {
	GetFamiliesByUserID(ctx context.Context, userID int64) ([]entity.Family, error)
	GetFamilyByCode(ctx context.Context, code string) (*entity.Family, time.Time, error)
	GetFamilyByID(ctx context.Context, id int) (*entity.Family, error)
}

func (s *FamilyService) GetFamiliesByUserID(ctx context.Context, userID int64) ([]entity.Family, error) {
	families, err := s.familyProvider.GetFamiliesByUserID(ctx, userID)
	if err != nil {
		s.sl.Error("failed to get family by user id", slog.Int("user_id", int(userID)), slog.String("err", err.Error()))
		return nil, err
	}

	if len(families) == 0 {
		return nil, errorsx.New("user has no family", errorsx.ErrCodeUserHasNoFamily, struct{}{})
	}

	return families, nil
}

func (s *FamilyService) GetFamilyByCode(ctx context.Context, code string) (*entity.Family, time.Time, error) {
	f, expiresAt, err := s.familyProvider.GetFamilyByCode(ctx, code)
	if err != nil {
		s.sl.Error("failed to get family by code", slog.String("code", code), slog.String("err", err.Error()))
		if errors.Is(err, pgx.ErrNoRows) {
			s.sl.Debug("family not found with code")
			return nil, time.Time{}, errorsx.New("family not found by invite code", errorsx.ErrCodeFamilyNotFound, struct{}{})
		}
		return nil, time.Time{}, err
	}

	return f, expiresAt, nil
}

func (s *FamilyService) GetFamilyByID(ctx context.Context, id int) (*entity.Family, error) {
	f, err := s.familyProvider.GetFamilyByID(ctx, id)
	if err != nil {
		s.sl.Error("failed to get family by id", slog.Int("id", id), slog.String("err", err.Error()))
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorsx.New("family not found by id", errorsx.ErrCodeFamilyNotFound, struct{}{})
		}
		return nil, err
	}

	return f, nil
}
