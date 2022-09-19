// Copyright Sabre

package prometheusremotewritereceiver

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/service/servicetest"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factory := NewFactory()
	factories.Receivers[typeStr] = factory
	cfg, err := servicetest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)

	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, 1, len(cfg.Receivers))

	defaultConfig := cfg.Receivers[config.NewComponentID(typeStr)]
	assert.Equal(t, factory.CreateDefaultConfig(), defaultConfig)

	dcfg := defaultConfig.(*Config)
	assert.Equal(t, "prometheusremotewrite", dcfg.ID().String())
	assert.Equal(t, "0.0.0.0:19291", dcfg.HTTPServerSettings.Endpoint)

}
