package userservice

import (
	"context"
	"monofamily/internal/entity"
	"monofamily/internal/pkg/sl"
	"monofamily/internal/service/familyservice"

	"github.com/jackc/pgx/v4"
)

type UserServiceIface interface {
	SaveUser(ctx context.Context, user *entity.User) (*entity.User, error)
	SaveUserToFamily(ctx context.Context, familyID int, userID int64) error
	GetAllUsersInFamily(ctx context.Context, familyID int) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	DeleteUserFromFamily(ctx context.Context, tx pgx.Tx, familyID int, userID int64) error

	WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error
}

type userSaver interface {
	SaveUser(ctx context.Context, user *entity.User) (*entity.User, error)
	SaveUserToFamily(ctx context.Context, familyID int, userID int64) error
}

type UserService struct {
	userSaver    userSaver
	userProvider userProvider
	userDeletor  userDeletor
	sl           sl.Logger

	monoApiUrl string

	tokenProvider   tokenProvider
	withTransaction familyservice.WithTransaction
}

func New(
	userIface UserServiceIface,
	sl sl.Logger,
	monoApiUrl string,
	tokenProvider tokenProvider,
) *UserService {
	return &UserService{
		userSaver:    userIface,
		userProvider: userIface,
		userDeletor:  userIface,
		sl:           sl,

		monoApiUrl: monoApiUrl,

		tokenProvider:   tokenProvider,
		withTransaction: userIface,
	}
}
