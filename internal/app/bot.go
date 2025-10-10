package app

import (
	"context"
	"monofamily/internal/adapter/database"
	"monofamily/internal/adapter/database/familyrepo"
	"monofamily/internal/adapter/database/tokenrepo"
	"monofamily/internal/adapter/database/userrepo"
	"monofamily/internal/delivery/telegram"
	"monofamily/internal/delivery/telegram/handler"
	"monofamily/internal/pkg/sl"
	"monofamily/internal/service/familyservice"
	"monofamily/internal/service/tokenservice"
	"monofamily/internal/service/userservice"
	"monofamily/internal/usecase"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	tb "gopkg.in/telebot.v3"
)

type TelegramBot struct {
	bot *tb.Bot

	pgsqlxpool *pgxpool.Pool

	encrKey [32]byte

	monoApiUrl string

	sl sl.Logger
}

type TBConfig struct {
	BotToken   string
	LongPoller int

	Pgsqlxpool *pgxpool.Pool

	EncrKey [32]byte

	MonoApiUrl string

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
		bot:        b,
		sl:         cfg.Logger,
		pgsqlxpool: cfg.Pgsqlxpool,
		encrKey:    cfg.EncrKey,
		monoApiUrl: cfg.MonoApiUrl,
	}

	return tgBot, nil
}

func (tgbot *TelegramBot) RunBot() {
	logger := tgbot.sl

	db := database.New(tgbot.pgsqlxpool)

	familyrepo := familyrepo.New(db.DB, logger)
	tokenrepo := tokenrepo.New(db.DB, logger)
	userrepo := userrepo.New(db.DB, logger)

	familyservice := familyservice.New(familyrepo, logger)
	tokenservice := tokenservice.New(tgbot.encrKey, tokenrepo, logger)
	userservice := userservice.New(userrepo, logger, tgbot.monoApiUrl, tokenservice)

	usecase := usecase.New(userservice, userservice, familyservice, tokenservice)

	handler := handler.New(usecase, tgbot.bot, logger)

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

	telegram.SetupRoutes(tgbot.bot, handler)

	tgbot.bot.Start()
}
