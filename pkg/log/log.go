package log

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func DefaultConsoleLogger() *slog.Logger {
	w := os.Stderr
	// create a new logger
	logger := slog.New(tint.NewHandler(w, nil))

	// set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.DateTime,
		}),
	))
	return logger
}
