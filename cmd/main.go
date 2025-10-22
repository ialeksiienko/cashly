package main

import (
	"cashly/internal/app"
	"cashly/internal/config"
	"cashly/internal/pkg/sl"
	"io"
	"log/slog"

	"os"
	"strings"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	var tokenEncrKey [32]byte
	copy(tokenEncrKey[:], []byte(cfg.Mono.EncryptKey))

	logger := setupLogger(cfg.Env)

	pgsqlxpool, _, err := app.NewDBPool(app.DatabaseConfig{
		Username: cfg.DB.User,
		Password: cfg.DB.Pass,
		Hostname: cfg.DB.Host,
		Port:     cfg.DB.Port,
		DBName:   cfg.DB.Name,

		Logger: logger,
	})
	if err != nil {
		logger.Fatal("unexpected error while tried to connect to database", slog.String("err", err.Error()))
	}

	defer pgsqlxpool.Close()

	tgBot, err := app.NewBot(app.TBConfig{
		BotToken:     cfg.Bot.Token,
		LongPoller:   cfg.Bot.LongPoller,
		Pgsqlxpool:   pgsqlxpool,
		EncrKey:      tokenEncrKey,
		MonoApiUrl:   cfg.Mono.ApiURL,
		AuthPassword: cfg.Bot.Password,
		Logger:       logger,
	})
	if err != nil {
		logger.Fatal("failed to build tg bot:", slog.String("error", err.Error()))
	}

	logger.Info("bot is running")

	tgBot.RunBot()
}

func setupLogger(env string) *sl.MyLogger {
	var level zerolog.Level
	var writer io.Writer

	var slogLevel slog.Level

	switch strings.ToLower(env) {
	case envLocal:
		level = zerolog.DebugLevel
		slogLevel = slog.LevelDebug

		writer = zerolog.ConsoleWriter{Out: os.Stdout}
	case envDev:
		level = zerolog.DebugLevel
		slogLevel = slog.LevelDebug

		writer = os.Stdout
	case envProd:
		level = zerolog.InfoLevel
		slogLevel = slog.LevelInfo

		writer = os.Stdout
	default:
		level = zerolog.InfoLevel
		slogLevel = slog.LevelInfo

		writer = os.Stdout
	}

	zlogger := zerolog.New(writer).Level(level).With().Timestamp().Logger()
	zl := &zlogger
	handler := slogzerolog.Option{
		Level:  slogLevel,
		Logger: zl,
	}.NewZerologHandler()

	return sl.New(slog.New(handler), func(msg string, attrs ...any) {
		fieldMap := make(map[string]any)

		for _, attr := range attrs {
			if a, ok := attr.(slog.Attr); ok {
				fieldMap[a.Key] = a.Value.Any()
			}
		}

		zlogger.Fatal().Fields(fieldMap).Msg(msg)
	})
}
