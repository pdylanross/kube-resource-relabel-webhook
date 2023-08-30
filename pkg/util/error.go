package util

import (
	"log/slog"
	"os"
)

func ErrCheck(err error) {
	if err != nil {
		slog.Error("fatal error", slog.Any("error", err))
		os.Exit(-1)
	}
}
