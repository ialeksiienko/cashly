package userservice

import (
	"cashly/internal/entity"
	"cashly/internal/pkg/errorsx"
	"context"
	"log/slog"
)

type provider interface {
	GetAllUsersInFamily(ctx context.Context, fid int) ([]entity.User, error)
	GetByID(ctx context.Context, fid int64) (*entity.User, error)
}

type Member struct {
	ID        int64
	Username  string
	Firstname string
	IsAdmin   bool
	IsCurrent bool
	HasToken  bool
}

func (s Service) GetFamilyMembers(ctx context.Context, f *entity.Family, uid int64) ([]Member, error) {
	users, err := s.provider.GetAllUsersInFamily(ctx, f.ID)
	if err != nil {
		s.logger.Error("failed to get all users in family", slog.String("family_name", f.Name), slog.String("err", err.Error()))
		return nil, err
	}

	if len(users) == 0 {
		return nil, errorsx.New("family has not members", errorsx.ErrCodeFamilyHasNoMembers, struct{}{})
	}

	members := make([]Member, len(users))
	for i, user := range users {
		hasToken, _, err := s.tokenProvider.Get(ctx, f.ID, user.ID)
		if err != nil {
			s.logger.Error("failed to get user bank token",
				slog.String("family_name", f.Name),
				slog.Int("user_id", int(user.ID)),
				slog.String("err", err.Error()),
			)
		}

		members[i] = Member{
			ID:        user.ID,
			Username:  user.Username,
			Firstname: user.Firstname,
			IsAdmin:   f.CreatedBy == user.ID,
			IsCurrent: user.ID == uid,
			HasToken:  hasToken,
		}
	}

	return members, nil
}

func (s Service) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	return s.provider.GetByID(ctx, id)
}

func (s Service) GetUsersByFamilyID(ctx context.Context, familyID int) ([]entity.User, error) {
	return s.provider.GetAllUsersInFamily(ctx, familyID)
}
