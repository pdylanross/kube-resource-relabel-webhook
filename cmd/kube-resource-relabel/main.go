package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/config"
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/handlers"
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/logger"
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/server"
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/util"
)

func main() {
	metricsHandlers := &handlers.MetricsHandlers{MetricsPath: "metrics"}
	relabelHandlers := &handlers.WebhookHandlers{}

	webhookServerConfig := config.ServerConfig{
		ListenAddress: ":8080",
		TLS:           nil,
		FriendlyName:  "webhook",
		SetupFunc:     relabelHandlers.SetupWebhookHandlers,
	}

	metricsServerConfig := config.ServerConfig{
		ListenAddress: ":8081",
		TLS:           nil,
		FriendlyName:  "metrics",
		SetupFunc:     metricsHandlers.SetupMetricsHandlers,
	}

	logCfg := config.LoggerConfig{
		Level:  "debug",
		Format: "text",
	}
	l, err := logger.BuildLogger(&logCfg)
	util.ErrCheck(err)

	gin.SetMode(gin.ReleaseMode)

	webhookSrv, err := server.NewServer(&webhookServerConfig, l)
	util.ErrCheck(err)

	metricsSrv, err := server.NewServer(&metricsServerConfig, l)
	util.ErrCheck(err)

	rg := util.MakeRunGroup()

	webhookSrv.Start(rg)
	metricsSrv.Start(rg)

	util.ErrCheck(rg.Run())
}
