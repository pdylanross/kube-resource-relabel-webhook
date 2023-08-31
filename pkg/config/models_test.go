package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

func TestRelabelConfig_ToRules(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	exampleDir := path.Join(path.Dir(filename), "../../example")
	exampleConfigs, err := filepath.Glob(path.Join(exampleDir, "*-good-config.yaml"))

	require.Nil(t, err, "err listing config dir")

	for _, cfg := range exampleConfigs {
		t.Run(fmt.Sprintf("TestRelabelConfig_ToRules-Example-%s", cfg), func(t *testing.T) {
			buf, err := os.ReadFile(cfg) //nolint:govet
			require.Nil(t, err, "err reading sample config file")

			var cfg RelabelConfig
			err = yaml.Unmarshal(buf, &cfg)

			assert.Nil(t, err, "err reading ")

			rules, err := cfg.ToRules()

			assert.Nil(t, err)
			assert.NotEmpty(t, rules)
		})
	}
}
