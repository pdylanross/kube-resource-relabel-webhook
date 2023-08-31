package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/config"
)

// BuildLogger creates a log/slog logger from the given log config.
func BuildLogger(cfg *config.LoggerConfig) (*slog.Logger, error) {
	opts, err := buildHandlerOptions(cfg)
	if err != nil {
		return nil, fmt.Errorf("error building logger: %s", err)
	}

	handler, err := buildHandler(cfg, opts)
	if err != nil {
		return nil, fmt.Errorf("error building logger: %s", err)
	}

	l := slog.New(handler)

	slog.SetDefault(l)

	return l, nil
}

func buildHandlerOptions(cfg *config.LoggerConfig) (*slog.HandlerOptions, error) {
	var programLevel = new(slog.LevelVar)

	switch strings.ToLower(cfg.Level) {
	case "debug":
		programLevel.Set(slog.LevelDebug)
	case "info":
		programLevel.Set(slog.LevelInfo)
	case "warn":
		programLevel.Set(slog.LevelWarn)
	case "error":
		programLevel.Set(slog.LevelError)
	default:
		return nil, fmt.Errorf("unknown log level %s", cfg.Level)
	}

	return &slog.HandlerOptions{Level: programLevel}, nil
}

func buildHandler(cfg *config.LoggerConfig, opts *slog.HandlerOptions) (slog.Handler, error) {
	switch strings.ToLower(cfg.Format) {
	case "text":
		return slog.NewTextHandler(os.Stdout, opts), nil
	case "json":
		return slog.NewJSONHandler(os.Stdout, opts), nil
	default:
		return nil, fmt.Errorf("unknown log format %s", cfg.Format)
	}
}
