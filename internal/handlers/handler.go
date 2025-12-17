package handlers

import (
	"cashly/internal/entity"
	userservice "cashly/internal/service/user"
	"cashly/pkg/slogx"
	"context"
	"errors"
	"time"

	tb "gopkg.in/telebot.v3"
)

const UfsKey = "user_family_state"

var (
	ErrInternalServerForUser = errors.New("Сталася помилка на боці серверу, спробуйте пізніше.")
	ErrUnableToGetUserState  = errors.New("Не вдалося отримати стан поточної сім'ї. Спробуйте пізніше.")
)

type UseCase interface {
	CreateFamily(context.Context, string, int64) (*entity.Family, string, time.Time, error)
	SelectFamily(context.Context, int, int64) (bool, bool, *entity.Family, error)
	RegisterUser(context.Context, *entity.User) (*entity.User, error)
	LeaveFamily(context.Context, *entity.Family, int64) error
	JoinFamily(context.Context, string, int64) (*entity.Family, error)

	GetBalance(context.Context, int, int64, string, string) (float64, error)
	GetFamilyMembers(context.Context, *entity.Family, int64) ([]userservice.Member, error)
	GetFamiliesByUserID(context.Context, int64) ([]entity.Family, error)
	GetUserByID(context.Context, int64) (*entity.User, error)

	SaveBankToken(context.Context, int, int64, string) (*entity.UserBankToken, error)
	DeleteUserBankToken(context.Context, int, int64) error

	// admin usecases
	RemoveMember(context.Context, int, int64, int64) error
	DeleteFamily(context.Context, *entity.Family, int64) error
	CreateNewInviteCode(context.Context, *entity.Family, int64) (string, time.Time, error)
}

type Handler struct {
	bot     *tb.Bot
	logger  slogx.Logger
	usecase UseCase

	eventCh chan entity.EventNotification

	AuthPassword string
}

func New(
	uc UseCase,
	bot *tb.Bot,
	eventCh chan entity.EventNotification,
	l slogx.Logger,
) *Handler {
	return &Handler{
		bot:     bot,
		logger:  l,
		usecase: uc,
		eventCh: eventCh,
	}
}
