package server

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func ginLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		if !logger.Enabled(context, slog.LevelDebug) {
			context.Next()

			return
		}

		start := time.Now()

		context.Next()

		elapsed := time.Since(start)

		req := context.Request
		logger.LogAttrs(context, slog.LevelDebug, "request handled", slog.String("route", req.RequestURI), slog.Duration("elapsed", elapsed), slog.Int("status_code", context.Writer.Status()))
	}
}
