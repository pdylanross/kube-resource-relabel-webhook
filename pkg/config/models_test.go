package config

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/relabel"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

func TestRelabelConfig_ToRules(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	exampleDir := path.Join(path.Dir(filename), "../../example/relabel-config")
	exampleConfigs, err := filepath.Glob(path.Join(exampleDir, "*-good-config.yaml"))

	require.Nil(t, err, "err listing config dir")

	for _, cfg := range exampleConfigs {
		t.Run(cfg, func(t *testing.T) {
			LoadRelabelRule(cfg, t)
		})
	}
}

func TestRelabelConfig_ToRules_AirflowDagConfig(t *testing.T) {
	configPath := GetSampleConfigFilePath("airflow-dag-good-config.yaml", t)
	rules := LoadRelabelRule(configPath, t)

	assert.NotEmpty(t, rules)
	assert.Equal(t, 1, len(rules))
}

func GetSampleConfigFilePath(configName string, t *testing.T) string {
	_, filename, _, _ := runtime.Caller(0)
	exampleDir := path.Join(path.Dir(filename), "../../example/relabel-config")
	configPath := path.Join(exampleDir, configName)

	_, err := os.Stat(configPath)
	require.Nil(t, err, "GetSampleConfigFilePath file error")

	return configPath
}

func LoadRelabelRule(file string, t *testing.T) []relabel.Rule {
	buf, err := os.ReadFile(file)
	require.Nil(t, err, "err reading sample config file")

	var cfg RelabelConfig
	err = yaml.Unmarshal(buf, &cfg)

	assert.Nil(t, err, "err reading ")

	rules, err := cfg.ToRules()

	assert.Nil(t, err)
	assert.NotEmpty(t, rules)

	return rules
}
