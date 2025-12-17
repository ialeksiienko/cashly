package main

import (
	"cashly/internal/app"
	"cashly/internal/config"
	"cashly/pkg/slogx"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	logger := slogx.New(cfg.Env)

	bot, err := app.New(cfg, logger)
	if err != nil {
		logger.Fatal("unable to create bot app", slog.String("err", err.Error()))
	}

	if runErr := bot.Run(); runErr != nil {
		logger.Fatal("unable to start bot", slog.String("err", runErr.Error()))
	}
}
