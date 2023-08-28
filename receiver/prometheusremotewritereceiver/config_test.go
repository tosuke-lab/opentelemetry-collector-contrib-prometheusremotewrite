// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheusremotewritereceiver

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configauth"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/confmap/confmaptest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusremotewritereceiver/internal/metadata"
)

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	tests := []struct {
		id       component.ID
		expected component.Config
	}{
		{
			expected: createDefaultConfig(),
		},
		{
			id: component.NewIDWithName(metadata.Type, "customname"),
			expected: &Config{
				HTTPServerSettings: confighttp.HTTPServerSettings{
					Endpoint:           "0.0.0.0:19291",
					TLSSetting:         (*configtls.TLSServerSetting)(nil),
					CORS:               (*confighttp.CORSSettings)(nil),
					Auth:               (*configauth.Authentication)(nil),
					MaxRequestBodySize: 0,
					IncludeMetadata:    false},
				TimeThreshold: 24,
			},
		},
	}

	for _, tt := range tests {
		t.Run("prometheusremotewrite", func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub("prometheusremotewrite")
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}
