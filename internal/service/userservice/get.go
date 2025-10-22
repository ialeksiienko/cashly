package userservice

import (
	"cashly/internal/entity"
	"cashly/internal/errorsx"
	"context"
	"log/slog"
)

type userProvider interface {
	GetAllUsersInFamily(ctx context.Context, familyID int) ([]entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
}

type MemberInfo struct {
	ID        int64
	Username  string
	Firstname string
	IsAdmin   bool
	IsCurrent bool
	HasToken  bool
}

func (s *UserService) GetFamilyMembersInfo(ctx context.Context, family *entity.Family, userID int64) ([]MemberInfo, error) {
	users, err := s.userProvider.GetAllUsersInFamily(ctx, family.ID)
	if err != nil {
		s.sl.Error("failed to get all users in family", slog.String("family_name", family.Name), slog.String("err", err.Error()))
		return nil, err
	}

	if len(users) == 0 {
		return nil, errorsx.New("family has not members", errorsx.ErrCodeFamilyHasNoMembers, struct{}{})
	}

	members := make([]MemberInfo, len(users))
	for i, user := range users {
		hasToken, _, err := s.tokenProvider.Get(ctx, family.ID, user.ID)
		if err != nil {
			s.sl.Warn("failed to get user bank token",
				slog.String("family_name", family.Name),
				slog.Int("user_id", int(user.ID)),
				slog.String("err", err.Error()),
			)
			continue
		}

		members[i] = MemberInfo{
			ID:        user.ID,
			Username:  user.Username,
			Firstname: user.Firstname,
			IsAdmin:   family.CreatedBy == user.ID,
			IsCurrent: user.ID == userID,
			HasToken:  hasToken,
		}
	}

	return members, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	return s.userProvider.GetByID(ctx, id)
}

func (s *UserService) GetUsersByFamilyID(ctx context.Context, familyID int) ([]entity.User, error) {
	return s.userProvider.GetAllUsersInFamily(ctx, familyID)
}
