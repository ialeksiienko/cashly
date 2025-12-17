package familyservice

import (
	"cashly/internal/entity"
	"context"
	"log/slog"
)

type creator interface {
	Create(ctx context.Context, inp *entity.Family) (*entity.Family, error)
}

func (s Service) Create(ctx context.Context, fname string, uid int64) (*entity.Family, error) {
	f, err := s.creator.Create(ctx, &entity.Family{
		Name:      fname,
		CreatedBy: uid,
	})
	if err != nil {
		s.logger.Error("failed to create family", slog.Int("user_id", int(uid)), slog.String("err", err.Error()))
		return nil, err
	}

	s.logger.Debug("family created", slog.Int("familyID", f.ID))

	return f, nil
}
