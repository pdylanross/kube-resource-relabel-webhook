package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Depado/ginprom"

	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/config"
)

type internalServerConfig struct {
	ListenAddress string
	TLS           *config.ServerTLSConfig

	FriendlyName string
	SetupFunc    func(router *gin.Engine) error
	Prometheus   *ginprom.Prometheus
}

func newInternalServerConfigFrom(cfg *config.ServerConfig, friendlyName string, setupFunc func(router *gin.Engine) error, prometheus *ginprom.Prometheus) internalServerConfig {
	return internalServerConfig{
		ListenAddress: cfg.ListenAddress,
		TLS:           cfg.TLS,
		FriendlyName:  friendlyName,
		SetupFunc:     setupFunc,
		Prometheus:    prometheus,
	}
}

type server struct {
	srv    http.Server
	logger *slog.Logger

	config internalServerConfig
}

func (s *server) start(group *run.Group) {
	group.Add(func() error {
		return s.beginListener()
	}, func(err error) {
		s.logger.Info("draining connections")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.srv.Shutdown(ctx); err != nil {
			s.logger.Error("error in shutdown", slog.Any("error", err))
		} else {
			s.logger.Info("shutdown gracefully")
		}
	})
}

func (s *server) beginListener() error {
	if s.config.TLS != nil {
		s.logger.Info("begin listen tls",
			slog.String("addr", s.config.ListenAddress),
			slog.String("tls_key_file", s.config.TLS.TLSKeyFilePath),
			slog.String("tls_cert_file", s.config.TLS.TLSCertFilePath),
		)
		return s.srv.ListenAndServeTLS(s.config.TLS.TLSCertFilePath, s.config.TLS.TLSKeyFilePath)
	}

	s.logger.Info("begin listen", slog.String("addr", s.config.ListenAddress))
	return s.srv.ListenAndServe()
}

func newServer(config internalServerConfig, logger *slog.Logger) (*server, error) {
	serverLogger := logger.WithGroup("server").With(slog.String("server_name", config.FriendlyName))

	srv := gin.New()
	srv.Use(config.Prometheus.Instrument())
	srv.Use(ginLogger(serverLogger))
	srv.Use(gin.Recovery())

	if err := config.SetupFunc(srv); err != nil {
		return nil, fmt.Errorf("error setting up server %s: %w", config.FriendlyName, err)
	}

	return &server{config: config, srv: http.Server{Addr: config.ListenAddress, Handler: srv}, logger: serverLogger}, nil
}
