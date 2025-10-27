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
	"os"
	"os/signal"
	"syscall"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func(ctx context.Context) {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Info("invite codes cleaner stopped")
				return
			case <-ticker.C:
				if err := familyservice.ClearInviteCodes(context.Background()); err != nil {
					logger.Error("failed to clear invite codes: " + err.Error())
				} else {
					logger.Debug("invite codes cleared successfully")
				}
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				logger.Info("event listener stopped")
				return
			case eventNt, ok := <-eventCh:
				if !ok {
					logger.Warn("event channel closed")
					return
				}

				var text string
				switch eventNt.Event {
				case entity.EventBalanceChecked:
					text = fmt.Sprintf("👤 Користувач (ID: %d) перевірив твій баланс у сім'ї [ %s ].",
						eventNt.Data["checked_by_user_id"].(int64), eventNt.FamilyName)

				case entity.EventJoinedFamily:
					text = fmt.Sprintf("🎉 Користувач (ID: %d) приєднався до твоєї сім’ї [ %s ].",
						eventNt.Data["joined_user_id"].(int64), eventNt.FamilyName)

				case entity.EventDeletedFromFamily:
					text = fmt.Sprintf("🥲 На жаль, вас видалили з сім'ї [ %s ].", eventNt.FamilyName)

				case entity.EventLeavedFromFamily:
					text = fmt.Sprintf("😔 Користувач (ID : %d) вийшов з твоєї сім'ї [ %s ].",
						eventNt.Data["leaved_user_id"].(int64), eventNt.FamilyName)
				}

				if text != "" {
					tgbot.bot.Send(&tb.User{ID: eventNt.RecipientID}, text)
				}
			}
		}
	}(ctx)

	telegram.SetupRoutes(tgbot.bot, tgbot.authPassword, handler)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(fmt.Sprintf("bot panic recovered: %v", r))
			}
		}()
		logger.Info("bot is running")
		tgbot.bot.Start()
	}()

	<-sigCh
	logger.Info("shutdown signal received")

	cancel()
	tgbot.bot.Stop()
	close(eventCh)

	logger.Info("bot stopped gracefully")
}
