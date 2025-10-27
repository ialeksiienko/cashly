package handler

import (
	"cashly/internal/entity"
	"cashly/internal/pkg/sl"
	"cashly/internal/service/userservice"
	"context"
	"errors"
	"time"

	tb "gopkg.in/telebot.v3"
)

var (
	ErrInternalServerForUser = errors.New("Сталася помилка на боці серверу, спробуйте пізніше.")
	ErrUnableToGetUserState  = errors.New("Не вдалося отримати стан поточної сім'ї. Спробуйте пізніше.")
)

type UseCase interface {
	CreateFamily(ctx context.Context, familyName string, userID int64) (*entity.Family, string, time.Time, error)
	SelectFamily(ctx context.Context, familyID int, userID int64) (bool, bool, *entity.Family, error)
	RegisterUser(ctx context.Context, user *entity.User) (*entity.User, error)
	LeaveFamily(ctx context.Context, family *entity.Family, userID int64) error
	JoinFamily(ctx context.Context, code string, userID int64) (*entity.Family, error)

	GetBalance(ctx context.Context, familyID int, checkedUserID int64, cardType string, currency string) (float64, error)
	GetFamilyMembersInfo(ctx context.Context, family *entity.Family, userID int64) ([]userservice.MemberInfo, error)
	GetFamiliesByUserID(ctx context.Context, userID int64) ([]entity.Family, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)

	SaveBankToken(ctx context.Context, familyID int, userID int64, token string) (*entity.UserBankToken, error)
	DeleteUserBankToken(ctx context.Context, familyID int, userID int64) error

	// admin usecases
	RemoveMember(ctx context.Context, familyID int, userID int64, memberID int64) error
	DeleteFamily(ctx context.Context, family *entity.Family, userID int64) error
	CreateNewInviteCode(ctx context.Context, family *entity.Family, userID int64) (string, time.Time, error)
}

type Handler struct {
	bot     *tb.Bot
	sl      sl.Logger
	usecase UseCase

	eventCh chan *entity.EventNotification
}

func New(
	uc UseCase,
	bot *tb.Bot,
	sl sl.Logger,
	eventCh chan *entity.EventNotification,
) *Handler {
	return &Handler{
		bot:     bot,
		sl:      sl,
		usecase: uc,
		eventCh: eventCh,
	}
}
