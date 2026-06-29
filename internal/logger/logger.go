package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	Log = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(Log)
}
