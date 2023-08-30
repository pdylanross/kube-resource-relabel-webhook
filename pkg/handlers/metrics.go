package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandlers struct {
	MetricsPath string
}

func (mh *MetricsHandlers) SetupMetricsHandlers(router *gin.Engine) error {
	setupHealthHandlers(router)

	router.Any(mh.MetricsPath, gin.WrapH(promhttp.Handler()))
	return nil
}

func setupHealthHandlers(router *gin.Engine) {
	router.GET("livez", func(context *gin.Context) {
		context.Status(200)
	})
}
