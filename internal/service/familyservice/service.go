package familyservice

import (
	"cashly/internal/entity"
	"cashly/internal/pkg/sl"
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

type FamilyServiceIface interface {
	Create(ctx context.Context, inp *entity.Family) (*entity.Family, error)
	GetFamiliesByUserID(ctx context.Context, userID int64) ([]entity.Family, error)
	GetFamilyByCode(ctx context.Context, code string) (*entity.Family, time.Time, error)
	GetFamilyByID(ctx context.Context, id int) (*entity.Family, error)
	GetInviteCode(ctx context.Context, familyID int) (string, time.Time, error)
	Delete(ctx context.Context, fn pgx.Tx, familyID int) error
	SaveFamilyInviteCode(ctx context.Context, userId int64, familyId int, code string) (time.Time, error)
	ClearInviteCodes(ctx context.Context) error

	WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error
}

type FamilyService struct {
	familyCreator            FamilyCreator
	familyProvider           FamilyProvider
	familyDeletor            FamilyDeletor
	familyInviteCodeSaver    FamilyInviteCodeSaver
	familyInviteCodeCleaner  FamilyInviteCodeCleaner
	familyInviteCodeProvider FamilyInviteCodeProvider
	withTransaction          WithTransaction
	sl                       sl.Logger
}

func New(
	familyIface FamilyServiceIface,
	sl sl.Logger,
) *FamilyService {
	return &FamilyService{
		familyCreator:            familyIface,
		familyProvider:           familyIface,
		familyDeletor:            familyIface,
		familyInviteCodeSaver:    familyIface,
		familyInviteCodeCleaner:  familyIface,
		familyInviteCodeProvider: familyIface,
		withTransaction:          familyIface,
		sl:                       sl,
	}
}
