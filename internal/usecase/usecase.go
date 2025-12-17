package usecase

import (
	"cashly/internal/entity"
	userservice "cashly/internal/service/user"
	"context"
	"time"
)

type userService interface {
	Register(context.Context, *entity.User) (*entity.User, error)
	SaveToFamily(context.Context, int, int64) error
	GetBalance(context.Context, int, int64, string, string) (float64, error)
	GetByID(context.Context, int64) (*entity.User, error)
	GetUsersByFamilyID(context.Context, int) ([]entity.User, error)
	GetFamilyMembers(context.Context, *entity.Family, int64) ([]userservice.Member, error)
	DeleteFromFamily(context.Context, int, int64) error
}

type adminService interface {
	DeleteFromFamily(context.Context, int, int64) error
}
type familyService interface {
	Create(context.Context, string, int64) (*entity.Family, error)
	GetFamiliesByUserID(context.Context, int64) ([]entity.Family, error)
	GetByCode(context.Context, string) (*entity.Family, time.Time, error)
	GetByID(context.Context, int) (*entity.Family, error)
	CreateNewInviteCode(context.Context, *entity.Family, int64) (string, time.Time, error)
	Delete(context.Context, int) error
}

type tokenService interface {
	Save(context.Context, int, int64, string) (*entity.UserBankToken, error)
	Get(context.Context, int, int64) (bool, *entity.UserBankToken, error)
	Delete(context.Context, int, int64) error
}

type UseCase struct {
	userService   userService
	adminService  adminService
	familyService familyService
	tokenService  tokenService
}

func New(
	userService userService,
	adminService adminService,
	familyService familyService,
	tokenService tokenService,
) *UseCase {
	return &UseCase{
		userService:   userService,
		adminService:  adminService,
		familyService: familyService,
		tokenService:  tokenService,
	}
}
