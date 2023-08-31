package server

import (
	"log/slog"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/config"
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/handlers"
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/logger"
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/util"
)

// KubeRelabelApp is the main app entrypoint interface.
type KubeRelabelApp interface {
	// Run initializes and runs the app, waiting until completion
	Run() error
}

type kubeRelabelApp struct {
	config config.AppConfig
}

func (k *kubeRelabelApp) Run() error {
	l, err := logger.BuildLogger(k.config.Logger)
	util.ErrCheck(err)

	l.Info("server startup", slog.Any("config", k.config))

	prom := ginprom.New(
		ginprom.Ignore("/metrics"),
	)
	metricsHandlers := &handlers.MetricsHandlers{Prometheus: prom}
	webhookHandlers := &handlers.WebhookHandlers{}

	gin.SetMode(gin.ReleaseMode)
	rg := util.MakeRunGroup()

	{
		webhookServerConfig := newInternalServerConfigFrom(&k.config.WebhookServer, "webhook", webhookHandlers.SetupWebhookHandlers, prom)
		webhookSrv, err := newServer(webhookServerConfig, l)
		util.ErrCheck(err)

		webhookSrv.start(rg)
	}

	if k.config.MetricsServer != nil {
		metricsServerConfig := newInternalServerConfigFrom(k.config.MetricsServer, "metrics", metricsHandlers.SetupMetricsHandlers, prom)
		metricsSrv, err := newServer(metricsServerConfig, l)
		util.ErrCheck(err)

		metricsSrv.start(rg)
	}

	return rg.Run()
}

// NewRelabelApp creates a relabeling app from a given config.
func NewRelabelApp(config config.AppConfig) KubeRelabelApp {
	return &kubeRelabelApp{config: config}
}
