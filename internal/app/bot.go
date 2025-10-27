package app

import (
	"cashly/internal/adapter/database"
	"cashly/internal/adapter/database/familyrepo"
	"cashly/internal/adapter/database/tokenrepo"
	"cashly/internal/adapter/database/userrepo"
	"cashly/internal/delivery/telegram"
	"cashly/internal/delivery/telegram/handler"
	"cashly/internal/entity"
	"cashly/internal/pkg/sl"
	"cashly/internal/service/familyservice"
	"cashly/internal/service/tokenservice"
	"cashly/internal/service/userservice"
	"cashly/internal/usecase"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	tb "gopkg.in/telebot.v3"
)

type TelegramBot struct {
	bot *tb.Bot

	pgsqlxpool *pgxpool.Pool

	encrKey [32]byte

	monoApiUrl string

	authPassword string

	sl sl.Logger
}

type TBConfig struct {
	BotToken   string
	LongPoller int

	Pgsqlxpool *pgxpool.Pool

	EncrKey [32]byte

	MonoApiUrl string

	AuthPassword string

	Logger sl.Logger
}

func NewBot(cfg TBConfig) (*TelegramBot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.BotToken,
		Poller: &tb.LongPoller{Timeout: time.Duration(cfg.LongPoller) * time.Second},
	})
	if err != nil {
		return nil, err
	}

	tgBot := &TelegramBot{
		bot:          b,
		sl:           cfg.Logger,
		pgsqlxpool:   cfg.Pgsqlxpool,
		encrKey:      cfg.EncrKey,
		monoApiUrl:   cfg.MonoApiUrl,
		authPassword: cfg.AuthPassword,
	}

	return tgBot, nil
}

func (tgbot *TelegramBot) RunBot() {
	logger := tgbot.sl
	eventCh := make(chan *entity.EventNotification, 100)

	db := database.New(tgbot.pgsqlxpool)

	familyrepo := familyrepo.New(db.DB, logger)
	tokenrepo := tokenrepo.New(db.DB, logger)
	userrepo := userrepo.New(db.DB, logger)

	familyservice := familyservice.New(familyrepo, logger)
	tokenservice := tokenservice.New(tgbot.encrKey, tokenrepo, logger)
	userservice := userservice.New(userrepo, logger, tgbot.monoApiUrl, tokenservice)

	usecase := usecase.New(userservice, userservice, familyservice, tokenservice)

	handler := handler.New(usecase, tgbot.bot, logger, eventCh)

	go func() {
		for {
			err := familyservice.ClearInviteCodes(context.Background())
			if err != nil {
				logger.Error(err.Error())
			} else {
				logger.Debug("invite codes cleared successfully")
			}
			time.Sleep(24 * time.Hour)
		}
	}()

	go func() {
		for eventNt := range eventCh {
			var text string

			switch eventNt.Event {
			case entity.EventBalanceChecked:
				text = fmt.Sprintf("👤 Користувач (ID: %d) перевірив твій баланс у сім'ї [ %s ].", eventNt.Data["checked_by_user_id"].(int64), eventNt.FamilyName)

			case entity.EventJoinedFamily:
				text = fmt.Sprintf("🎉 Користувач (ID: %d) приєднався до твоєї сім’ї [ %s ].", eventNt.Data["joined_user_id"].(int64), eventNt.FamilyName)

			case entity.EventDeletedFromFamily:
				text = fmt.Sprintf("🥲 На жаль, вас видалили з сім'ї [ %s ].", eventNt.FamilyName)

			case entity.EventLeavedFromFamily:
				text = fmt.Sprintf("😔 Користувач (ID : %d) вийшов з твоєї сім'ї [ %s ].", eventNt.Data["leaved_user_id"].(int64), eventNt.FamilyName)
			}

			if text != "" {
				tgbot.bot.Send(&tb.User{ID: eventNt.RecipientID}, text)
			}
		}
	}()

	telegram.SetupRoutes(tgbot.bot, tgbot.authPassword, handler)

	tgbot.bot.Start()
}
