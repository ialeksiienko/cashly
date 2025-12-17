package familyservice

import (
	"cashly/internal/entity"
	"cashly/internal/pkg/errorsx"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v4"
)

type provider interface {
	GetFamiliesByUserID(ctx context.Context, uid int64) ([]entity.Family, error)
	GetByCode(ctx context.Context, code string) (*entity.Family, time.Time, error)
	GetByID(ctx context.Context, id int) (*entity.Family, error)
}

func (s Service) GetFamiliesByUserID(ctx context.Context, uid int64) ([]entity.Family, error) {
	families, err := s.provider.GetFamiliesByUserID(ctx, uid)
	if err != nil {
		s.logger.Error("failed to get family by user id", slog.Int("user_id", int(uid)), slog.String("err", err.Error()))
		return nil, err
	}

	if len(families) == 0 {
		return nil, errorsx.New("user has no family", errorsx.ErrCodeUserHasNoFamily, struct{}{})
	}

	return families, nil
}

func (s Service) GetByCode(ctx context.Context, code string) (*entity.Family, time.Time, error) {
	f, expiresAt, err := s.provider.GetByCode(ctx, code)
	if err != nil {
		s.logger.Error("failed to get family by code", slog.String("code", code), slog.String("err", err.Error()))
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Debug("family not found with code")
			return nil, time.Time{}, errorsx.New("family not found by invite code", errorsx.ErrCodeFamilyNotFound, struct{}{})
		}
	}

	return f, expiresAt, err
}

func (s Service) GetByID(ctx context.Context, id int) (*entity.Family, error) {
	f, err := s.provider.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get family by id", slog.Int("id", id), slog.String("err", err.Error()))
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorsx.New("family not found by id", errorsx.ErrCodeFamilyNotFound, struct{}{})
		}
	}

	return f, err
}
