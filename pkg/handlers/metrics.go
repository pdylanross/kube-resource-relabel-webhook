package handlers

import (
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
)

// MetricsHandlers encapsulates properties required to setup metrics endpoints.
type MetricsHandlers struct {
	Prometheus *ginprom.Prometheus
}

// SetupMetricsHandlers registers http handlers for metrics endpoints.
func (mh *MetricsHandlers) SetupMetricsHandlers(router *gin.Engine) error {
	setupHealthHandlers(router)

	mh.Prometheus.Use(router)
	return nil
}
