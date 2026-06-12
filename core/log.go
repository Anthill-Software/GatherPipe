package core

import (
	"log/slog"
	"strings"
)

// Intercepteur magique qui convertit hclog en slog
type hclogToSlogRedirector struct{}

func (r hclogToSlogRedirector) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))

	// On détecte le préfixe hclog, on le nettoie, et on envoie au bon niveau slog
	switch {
	case strings.HasPrefix(msg, "[DEBUG]"):
		cleanMsg := strings.TrimSpace(strings.TrimPrefix(msg, "[DEBUG]"))
		slog.Debug(cleanMsg)
	case strings.HasPrefix(msg, "[INFO]"):
		cleanMsg := strings.TrimSpace(strings.TrimPrefix(msg, "[INFO]"))
		slog.Info(cleanMsg)
	case strings.HasPrefix(msg, "[WARN]"):
		cleanMsg := strings.TrimSpace(strings.TrimPrefix(msg, "[WARN]"))
		slog.Warn(cleanMsg)
	case strings.HasPrefix(msg, "[ERROR]"):
		cleanMsg := strings.TrimSpace(strings.TrimPrefix(msg, "[ERROR]"))
		slog.Error(cleanMsg)
	default:
		slog.Info(msg)
	}

	return len(p), nil
}
