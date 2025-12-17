package familyservice

import (
	"cashly/internal/adapter/database"
	"cashly/internal/entity"
	"cashly/pkg/slogx"
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

type FamilyIface interface {
	Create(ctx context.Context, inp *entity.Family) (*entity.Family, error)
	GetFamiliesByUserID(ctx context.Context, uid int64) ([]entity.Family, error)
	GetByCode(ctx context.Context, code string) (*entity.Family, time.Time, error)
	GetByID(ctx context.Context, id int) (*entity.Family, error)
	GetInviteCode(ctx context.Context, fid int) (string, time.Time, error)
	Delete(ctx context.Context, fn pgx.Tx, fid int) error
	SaveInviteCode(ctx context.Context, uid int64, fid int, code string) (time.Time, error)
	ClearInviteCodes(ctx context.Context) error
}

type Service struct {
	creator            creator
	provider           provider
	deletor            deletor
	inviteCodeSaver    inviteCodeSaver
	inviteCodeCleaner  inviteCodeCleaner
	inviteCodeProvider inviteCodeProvider
	withTransaction    database.WithTransaction
	logger             slogx.Logger
}

func New(
	fiface FamilyIface,
	tx database.WithTransaction,
	l slogx.Logger,
) Service {
	return Service{
		creator:            fiface,
		provider:           fiface,
		deletor:            fiface,
		inviteCodeSaver:    fiface,
		inviteCodeCleaner:  fiface,
		inviteCodeProvider: fiface,
		withTransaction:    tx,
		logger:             l,
	}
}
