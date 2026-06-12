package core

import (
	"context"
	"log/slog"
	"path/filepath"
	"strings"
)

type hclogToSlogRedirector struct {
	pluginName string
	pluginType string
}

func NewLinkRedirector(s string, t string) hclogToSlogRedirector {
	return hclogToSlogRedirector{
		pluginName: s,
		pluginType: t,
	}
}

func (r hclogToSlogRedirector) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	ctx := context.Background()

	tokens := map[string]slog.Level{
		"[TRACE]": slog.LevelDebug,
		"[DEBUG]": slog.LevelDebug,
		"[INFO]":  slog.LevelInfo,
		"[WARN]":  slog.LevelWarn,
		"[ERROR]": slog.LevelError,
	}

	level := slog.LevelInfo

	for token, lvl := range tokens {
		if strings.HasPrefix(msg, token) {
			level = lvl
			msg = strings.TrimSpace(strings.TrimPrefix(msg, token))
			break
		}
	}

	if idx := strings.Index(msg, ":"); idx != -1 {
		msg = strings.TrimSpace(msg[idx+1:])
	}

	for token := range tokens {
		if strings.HasPrefix(msg, token) {
			msg = strings.TrimSpace(strings.TrimPrefix(msg, token))
			break
		}
	}

	slog.Log(ctx, level, msg,
		"type", r.pluginType,
		"sender", filepath.Base(r.pluginName),
	)

	return len(p), nil
}
