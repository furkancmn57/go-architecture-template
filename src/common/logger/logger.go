// Package logger configures the process-wide slog default handler.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Init sets the global slog default. Call once after config load.
// local/development → text + debug (visible in make run / make run-dev);
// other envs → JSON + info.
func Init(appEnv string) {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}

	var handler slog.Handler
	switch strings.ToLower(strings.TrimSpace(appEnv)) {
	case "local", "development", "dev":
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

// IsDev reports whether the process should use developer-friendly logging (HTTP access logs, etc.).
func IsDev(appEnv string) bool {
	switch strings.ToLower(strings.TrimSpace(appEnv)) {
	case "local", "development", "dev":
		return true
	default:
		return false
	}
}
