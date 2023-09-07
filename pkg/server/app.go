package server

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/mutation"
	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/relabel"
	"gopkg.in/yaml.v3"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/config"
	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/handlers"
	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/logger"
	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/util"
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

	l.Info("server initializing", slog.Any("config", k.config))

	mutator, err := k.buildMutator(l)
	util.ErrCheck(err)

	prom := ginprom.New(
		ginprom.Ignore("/metrics"),
	)
	metricsHandlers := &handlers.MetricsHandlers{Prometheus: prom}
	webhookHandlers := handlers.NewWebhookHandlers(l, mutator)

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

func (k *kubeRelabelApp) buildMutator(logger *slog.Logger) (*mutation.Mutator, error) {
	buf, err := os.ReadFile(k.config.RelabelConfigFile)
	if err != nil {
		return nil, fmt.Errorf("err reading relabel config %w", err)
	}

	var cfg config.RelabelConfig
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, fmt.Errorf("err deserializing relabel config %w", err)
	}

	rules, err := cfg.ToRules()
	if err != nil {
		return nil, fmt.Errorf("err building rules %w", err)
	}

	return mutation.NewMutator(relabel.NewRelabeler(rules), logger), nil
}

// NewRelabelApp creates a relabeling app from a given config.
func NewRelabelApp(config config.AppConfig) KubeRelabelApp {
	return &kubeRelabelApp{config: config}
}
