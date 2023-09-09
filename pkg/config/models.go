package config

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/relabel"

	"gopkg.in/yaml.v3"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/util"
)

// AppConfig describes how the app should behave.
type AppConfig struct {
	// WebhookServer describes how the webhook http server should behave
	WebhookServer ServerConfig
	// MetricsServer describes how the metrics server should behave, or disable if nil
	MetricsServer *ServerConfig
	// Logger describes how to setup the logger
	Logger *LoggerConfig
	// RelabelConfigFile is the file path for the ruleset configuring how we should relabel
	RelabelConfigFile string
	// Version is the version of the currently running app
	Version AppVersionConfig
}

func (a *AppConfig) LogValue() slog.Value {
	marshalled, err := json.Marshal(a)

	// failure to serialize the appconfig should never happen & is panic worthy
	// if this ever hits you in the wild, i'm sorry
	util.ErrCheck(err)

	return slog.StringValue(string(marshalled))
}

// ServerConfig describes how a server should behave.
type ServerConfig struct {
	// ListenAddress defines the address the server should listen on
	ListenAddress string
	// TLS describes TLS options our server should use, or use plain http if nil
	TLS *ServerTLSConfig
}

// ServerTLSConfig describes how we should listen on TLS.
type ServerTLSConfig struct {
	// TLSCertFilePath describes the file path for our cert file
	TLSCertFilePath string
	// TLSKeyFilePath describes the private key file
	TLSKeyFilePath string
}

// LoggerConfig describes how we should setup the logger.
type LoggerConfig struct {
	// Level describes the lowest log level we should emit
	// one of debug, info, warn, or error
	Level string
	// Format describes the log format we should write
	// one of text, json
	Format string
}

// AppVersionConfig describes the current app version.
type AppVersionConfig struct {
	// Version is the semver version of the app
	Version string
	// CommitHash is the git version of the app
	CommitHash string
	// BuildTimestamp is the ts of the current build
	BuildTimestamp string
}

func (a *AppVersionConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("version", a.Version),
		slog.String("commit-hash", a.CommitHash),
		slog.String("build-timestamp", a.BuildTimestamp),
	)
}

// RelabelConfig is the config defining how we should relabel k8s objects.
type RelabelConfig struct {
	// Relabel is the list of rules
	Relabel []RelabelConfigRule `yaml:"relabel"`
}

func (r *RelabelConfig) ToRules() ([]relabel.Rule, error) {
	rules := make([]relabel.Rule, len(r.Relabel))

	for idx, r := range r.Relabel {
		newRule, err := r.ToRelabelRule()
		if err != nil {
			return nil, fmt.Errorf("error building rules: %s", err)
		}

		rules[idx] = *newRule
	}

	return rules, nil
}

// RelabelConfigRule is a single relabel rule.
type RelabelConfigRule struct {
	// Name is the friendly name for the rule
	Name string `yaml:"name"`
	// Conditions are the conditions required to perform the actions
	// these are all required for the actions to be applied
	Conditions []RelabelConfigRuleCondition `yaml:"conditions"`
	// Actions are what are performed when all of the conditions are met
	Actions []RelabelConfigRuleAction `yaml:"actions"`
}

func (r *RelabelConfigRule) ToRelabelRule() (*relabel.Rule, error) {
	conditions := make([]relabel.Condition, len(r.Conditions))
	actions := make([]relabel.Action, len(r.Actions))

	for idx, c := range r.Conditions {
		concrete, err := c.GetConcrete()
		if err != nil {
			return nil, fmt.Errorf("error constructing condition for rule %s: %s", r.Name, err)
		}

		conditions[idx] = concrete.ToCondition()
	}

	for idx, a := range r.Actions {
		concrete, err := a.GetConcrete()
		if err != nil {
			return nil, fmt.Errorf("error constructing action for rule %s: %s", r.Name, err)
		}

		actions[idx] = concrete.ToAction()
	}

	return relabel.NewRelabelRule(r.Name, conditions, actions), nil
}

// RelabelConfigRuleCondition is a singular condition to be met.
type RelabelConfigRuleCondition struct {
	// Type is the type of condition
	Type string `yaml:"type"`
	// Value is the specific configuration for this condition
	Value yaml.Node `yaml:"value"`
}

// RelabelConfigRuleAction is a singular action to be performed.
type RelabelConfigRuleAction struct {
	// Type is the type of action to be performed
	Type string `yaml:"type"`
	// Value is the specific configuration for this action
	Value yaml.Node `yaml:"value"`
}
