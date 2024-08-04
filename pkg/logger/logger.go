package logger

import (
	"log/slog"
	"os"
)

const (
	debug = "debug"
	dev   = "dev"
	prod  = "prod"
)

func New(logLevel string) *slog.Logger {
	var log *slog.Logger

	switch logLevel {
	case debug:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case dev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case prod:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
