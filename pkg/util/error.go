package util

import (
	"log/slog"
	"os"
)

// ErrCheck checks an error and exits the program if it's not nil.
func ErrCheck(err error) {
	if err != nil {
		slog.Error("fatal error", slog.Any("error", err))
		os.Exit(-1)
	}
}
