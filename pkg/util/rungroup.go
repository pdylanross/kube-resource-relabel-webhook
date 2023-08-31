package util

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
)

// MakeRunGroup creates a new rungroup and registers a SIGTERM and SIGINT handler to it.
func MakeRunGroup() *run.Group {
	var rg run.Group

	addExitHandlersToRunGroup(&rg)
	return &rg
}

func addExitHandlersToRunGroup(rg *run.Group) {
	sigC := make(chan os.Signal, 1)
	sigExit := make(chan struct{})
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)

	rg.Add(func() error {
		select {
		case s := <-sigC:
			slog.Info("exit signal received", slog.String("signal", s.String()))
			return nil
		case <-sigExit:
			return nil
		}
	}, func(_ error) {
		close(sigExit)
	})
}
