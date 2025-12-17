package userservice

import (
	"cashly/internal/adapter/database"
	"cashly/internal/entity"
	"cashly/pkg/slogx"
	"context"

	"github.com/jackc/pgx/v4"
)

type UserIface interface {
	Save(ctx context.Context, user *entity.User) (*entity.User, error)
	SaveToFamily(ctx context.Context, familyID int, userID int64) error
	GetAllUsersInFamily(ctx context.Context, familyID int) ([]entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	DeleteFromFamily(ctx context.Context, tx pgx.Tx, familyID int, userID int64) error
}

type saver interface {
	Save(ctx context.Context, user *entity.User) (*entity.User, error)
	SaveToFamily(ctx context.Context, familyID int, userID int64) error
}

type Service struct {
	saver    saver
	provider provider
	deletor  deletor
	logger   slogx.Logger

	monoApiUrl string

	tokenProvider   TokenProvider
	withTransaction database.WithTransaction
}

func New(
	userIface UserIface,
	withTransaction database.WithTransaction,
	monoApiUrl string,
	tokenProvider TokenProvider,
	l slogx.Logger,
) Service {
	return Service{
		saver:    userIface,
		provider: userIface,
		deletor:  userIface,
		logger:   l,

		monoApiUrl: monoApiUrl,

		tokenProvider:   tokenProvider,
		withTransaction: withTransaction,
	}
}
