package familyservice

import (
	"cashly/internal/entity"
	"context"
	"log/slog"
)

type FamilyCreator interface {
	Create(ctx context.Context, inp *entity.Family) (*entity.Family, error)
}

func (s *FamilyService) Create(ctx context.Context, familyName string, userID int64) (*entity.Family, error) {
	f, err := s.familyCreator.Create(ctx, &entity.Family{
		Name:      familyName,
		CreatedBy: userID,
	})
	if err != nil {
		s.sl.Error("failed to create family", slog.Int("user_id", int(userID)), slog.String("err", err.Error()))
		return nil, err
	}

	s.sl.Debug("family created", slog.Int("familyID", f.ID))

	return f, nil
}
