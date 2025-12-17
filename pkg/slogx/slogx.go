package slogx

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog"
)

type Logger interface {
	Info(msg string, attrs ...slog.Attr)
	Error(msg string, attrs ...slog.Attr)
	Debug(msg string, attrs ...slog.Attr)
	Warn(msg string, attrs ...slog.Attr)
	Fatal(msg string, attrs ...slog.Attr)
	With(attrs ...slog.Attr) Logger
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type CustomLogger struct {
	logger    *slog.Logger
	fatalFunc func(msg string, attrs ...any)
}

func New(env string) *CustomLogger {
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

	return &CustomLogger{logger: slog.New(handler), fatalFunc: func(msg string, attrs ...any) {
		fieldMap := make(map[string]any)

		for _, attr := range attrs {
			if a, ok := attr.(slog.Attr); ok {
				fieldMap[a.Key] = a.Value.Any()
			}
		}

		zlogger.Fatal().Fields(fieldMap).Msg(msg)
	}}
}

func (l *CustomLogger) Info(msg string, attrs ...slog.Attr) {
	l.logger.Info(msg, attrsToAny(attrs)...)
}

func (l *CustomLogger) Error(msg string, attrs ...slog.Attr) {
	l.logger.Error(msg, attrsToAny(attrs)...)
}

func (l *CustomLogger) Debug(msg string, attrs ...slog.Attr) {
	l.logger.Debug(msg, attrsToAny(attrs)...)
}

func (l *CustomLogger) Warn(msg string, attrs ...slog.Attr) {
	l.logger.Warn(msg, attrsToAny(attrs)...)
}

func (l *CustomLogger) Fatal(msg string, attrs ...slog.Attr) {
	l.fatalFunc(msg, attrsToAny(attrs)...)
}

func (l *CustomLogger) With(attrs ...slog.Attr) Logger {
	return &CustomLogger{logger: l.logger.With(attrsToAny(attrs)...)}
}

func attrsToAny(attrs []slog.Attr) []any {
	args := make([]any, len(attrs))
	for i, a := range attrs {
		args[i] = a
	}
	return args
}
