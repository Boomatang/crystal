package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

// Init initializes the global logger with specified level and format
func Init(level slog.Level, format string) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	Log = slog.New(handler)
}

// InitFromEnv initializes logger from environment variables
// LOG_LEVEL: DEBUG, INFO, WARN, ERROR (default: INFO)
// LOG_FORMAT: text, json (default: text)
func InitFromEnv() {
	level := slog.LevelInfo
	format := "text"

	// Read from environment
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		switch envLevel {
		case "DEBUG":
			level = slog.LevelDebug
		case "INFO":
			level = slog.LevelInfo
		case "WARN":
			level = slog.LevelWarn
		case "ERROR":
			level = slog.LevelError
		}
	}

	if envFormat := os.Getenv("LOG_FORMAT"); envFormat != "" {
		format = envFormat
	}

	Init(level, format)
}

func init() {
	// Default initialization
	Init(slog.LevelInfo, "text")
}
