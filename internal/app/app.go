package app

import (
	"cashly/internal/adapter/database"
	familyrepo "cashly/internal/adapter/repository/family"
	tokenrepo "cashly/internal/adapter/repository/token"
	userrepo "cashly/internal/adapter/repository/user"
	"cashly/internal/config"
	"cashly/internal/entity"
	"cashly/internal/handler"
	"cashly/internal/router"
	familyservice "cashly/internal/service/family"
	tokenservice "cashly/internal/service/token"
	userservice "cashly/internal/service/user"
	"cashly/internal/usecase"
	"cashly/pkg/slogx"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	tb "gopkg.in/telebot.v3"
)

type App struct {
	bot          *tb.Bot
	cfg          config.Config
	monoApiUrl   string
	authPassword string
	logger       slogx.Logger
}

func New(cfg config.Config, l slogx.Logger) (App, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.Bot.Token,
		Poller: &tb.LongPoller{Timeout: time.Duration(cfg.Bot.LongPoller) * time.Second},
	})
	if err != nil {
		return App{}, err
	}

	return App{
		bot:          b,
		cfg:          cfg,
		monoApiUrl:   cfg.Mono.ApiURL,
		authPassword: cfg.Bot.Password,
		logger:       l,
	}, nil
}

func (a App) Run() error {
	l := a.logger
	tbot := a.bot

	key, err := config.ConvertTokenKeyToBytes(a.cfg.Mono.EncryptKey)
	if err != nil {
		return err
	}

	pool, closedb, err := database.NewDBPool(database.Config{
		User: a.cfg.DB.User,
		Pass: a.cfg.DB.Pass,
		Host: a.cfg.DB.Host,
		Port: a.cfg.DB.Port,
		Name: a.cfg.DB.Name,

		Logger: l,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("unexpected error while trying to connect to database: %s", err.Error()))
	}
	defer closedb()

	db := database.New(pool, l)

	eventCh := make(chan entity.EventNotification, 100)

	familyRepo := familyrepo.New(db, l)
	tokenRepo := tokenrepo.New(db, l)
	userRepo := userrepo.New(db, l)

	familyService := familyservice.New(familyRepo, db, l)
	tokenService := tokenservice.New(tokenservice.NewEncrypt(key), tokenRepo, l)
	userService := userservice.New(userRepo, db, a.monoApiUrl, tokenService, l)

	adminService := userService

	uc := usecase.New(userService, adminService, familyService, tokenService)

	h := handler.New(uc, tbot, eventCh, l)

	h.AuthPassword = a.authPassword

	router.SetupRoutes(tbot, h)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go a.inviteCodesCleaner(ctx, familyService)

	go a.eventListener(ctx, eventCh)

	go a.start()

	<-sigCh
	l.Info("shutdown signal received")

	cancel()
	tbot.Stop()

	l.Info("bot stopped gracefully")

	return nil
}

const lkey = "func"

func (a App) inviteCodesCleaner(ctx context.Context, fs familyservice.Service) {
	l := a.logger.With(slog.String(lkey, "invite codes cleaner"))

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			l.Info("invite codes cleaner stopped")
			return
		case <-ticker.C:
			if err := fs.ClearInviteCodes(context.Background()); err != nil {
				l.Error("failed to clear invite codes: " + err.Error())
			} else {
				l.Debug("invite codes cleared successfully")
			}
		}
	}
}

func (a App) eventListener(ctx context.Context, ch chan entity.EventNotification) {
	l := a.logger.With(slog.String(lkey, "event listener"))

	for {
		select {
		case <-ctx.Done():
			l.Info("event listener stopped")
			return
		case e, ok := <-ch:
			if !ok {
				l.Warn("event channel closed")
				return
			}

			var text string
			switch e.Type {
			case entity.EventBalanceChecked:
				text = fmt.Sprintf("👤 Користувач (ID: %d) перевірив твій баланс у сім'ї [ %s ].",
					e.Data["checked_by_user_id"].(int64), e.FamilyName)

			case entity.EventJoinedFamily:
				text = fmt.Sprintf("🎉 Користувач (ID: %d) приєднався до твоєї сім’ї [ %s ].",
					e.Data["joined_user_id"].(int64), e.FamilyName)

			case entity.EventDeletedFromFamily:
				text = fmt.Sprintf("🥲 На жаль, вас видалили з сім'ї [ %s ].", e.FamilyName)

			case entity.EventLeavedFromFamily:
				text = fmt.Sprintf("😔 Користувач (ID : %d) вийшов з твоєї сім'ї [ %s ].",
					e.Data["leaved_user_id"].(int64), e.FamilyName)
			}

			if text != "" {
				a.bot.Send(&tb.User{ID: e.RecipientID}, text)
			}
		}
	}
}

func (a App) start() {
	l := a.logger

	defer func() {
		if r := recover(); r != nil {
			l.Error(fmt.Sprintf("bot panic recovered: %v", r))
		}
	}()

	l.Info("bot is running")

	a.bot.Start()
}
