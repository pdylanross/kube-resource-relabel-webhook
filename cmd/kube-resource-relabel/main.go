package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/config"
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/server"
)

type GlobalOpts struct {
	LogLevel  string
	LogFormat string
}

func (o *GlobalOpts) ToLogConfig() *config.LoggerConfig {
	return &config.LoggerConfig{
		Level:  o.LogLevel,
		Format: o.LogFormat,
	}
}

type ServeCmd struct {
	RelabelConfigFile string `help:"filepath for the relabel config" required:"true" type:"path"`

	WebhookHost        string `help:"hostname for the webhook" default:"0.0.0.0"`
	WebhookPort        uint16 `help:"port for the webhook" default:"8443"`
	WebhookTLSCertFile string `help:"path to the TLS cert file for the webhook" type:"path"`
	WebhookTLSKeyFile  string `help:"path to the TLS key file for the webhook" type:"path"`

	MetricsEnabled     bool   `help:"if the metrics endpoint is enabled"`
	MetricsHost        string `help:"hostname for the webhook" default:"0.0.0.0"`
	MetricsPort        uint16 `help:"port for the webhook" default:"8001"`
	MetricsTLSCertFile string `help:"path to the TLS cert file for the webhook" type:"path"`
	MetricsTLSKeyFile  string `help:"path to the TLS key file for the webhook" type:"path"`
}

func (s *ServeCmd) ToWebhookConfig() config.ServerConfig {
	var tlsConfig *config.ServerTLSConfig
	if s.WebhookTLSCertFile != "" && s.WebhookTLSKeyFile != "" {
		tlsConfig = &config.ServerTLSConfig{
			TLSCertFilePath: s.WebhookTLSCertFile,
			TLSKeyFilePath:  s.WebhookTLSKeyFile,
		}
	}

	return config.ServerConfig{
		ListenAddress: fmt.Sprintf("%s:%d", s.WebhookHost, s.WebhookPort),
		TLS:           tlsConfig,
	}
}

func (s *ServeCmd) ToMetricsConfig() *config.ServerConfig {
	if !s.MetricsEnabled {
		return nil
	}

	var tlsConfig *config.ServerTLSConfig
	if s.MetricsTLSCertFile != "" && s.MetricsTLSKeyFile != "" {
		tlsConfig = &config.ServerTLSConfig{
			TLSCertFilePath: s.MetricsTLSCertFile,
			TLSKeyFilePath:  s.MetricsTLSKeyFile,
		}
	}

	return &config.ServerConfig{
		ListenAddress: fmt.Sprintf("%s:%d", s.MetricsHost, s.MetricsPort),
		TLS:           tlsConfig,
	}
}

func (s *ServeCmd) Run(opts *GlobalOpts) error {
	appConf := config.AppConfig{
		WebhookServer:     s.ToWebhookConfig(),
		MetricsServer:     s.ToMetricsConfig(),
		Logger:            opts.ToLogConfig(),
		RelabelConfigFile: s.RelabelConfigFile,
	}

	app := server.NewRelabelApp(appConf)
	return app.Run()
}

var CLI struct {
	LogLevel  string `help:"the log level" enum:"debug,info,warn,error" default:"info"`
	LogFormat string `help:"the log format" enum:"json,text" default:"text"`

	Serve ServeCmd `cmd:"" help:"Run the server"`
}

func main() {
	ctx := kong.Parse(
		&CLI,
		kong.Name("kube-resource-relabel-webhook"),
		kong.ConfigureHelp(kong.HelpOptions{
			NoAppSummary:        false,
			Summary:             false,
			Compact:             false,
			Tree:                false,
			FlagsLast:           false,
			NoExpandSubcommands: false,
		}),
	)
	globals := GlobalOpts{
		LogLevel:  CLI.LogLevel,
		LogFormat: CLI.LogFormat,
	}

	if err := ctx.Run(&globals); err != nil {
		ctx.Errorf("err: %s", err)
	}
}
