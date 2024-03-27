// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

func TestMetricsBuilderConfig(t *testing.T) {
	tests := []struct {
		name string
		want MetricsBuilderConfig
	}{
		{
			name: "default",
			want: DefaultMetricsBuilderConfig(),
		},
		{
			name: "all_set",
			want: MetricsBuilderConfig{
				Metrics: MetricsConfig{
					K8sHpaCurrentReplicas: MetricConfig{Enabled: true},
					K8sHpaDesiredReplicas: MetricConfig{Enabled: true},
					K8sHpaMaxReplicas:     MetricConfig{Enabled: true},
					K8sHpaMinReplicas:     MetricConfig{Enabled: true},
				},
				ResourceAttributes: ResourceAttributesConfig{
					K8sHpaName:       ResourceAttributeConfig{Enabled: true},
					K8sHpaUID:        ResourceAttributeConfig{Enabled: true},
					K8sNamespaceName: ResourceAttributeConfig{Enabled: true},
				},
			},
		},
		{
			name: "none_set",
			want: MetricsBuilderConfig{
				Metrics: MetricsConfig{
					K8sHpaCurrentReplicas: MetricConfig{Enabled: false},
					K8sHpaDesiredReplicas: MetricConfig{Enabled: false},
					K8sHpaMaxReplicas:     MetricConfig{Enabled: false},
					K8sHpaMinReplicas:     MetricConfig{Enabled: false},
				},
				ResourceAttributes: ResourceAttributesConfig{
					K8sHpaName:       ResourceAttributeConfig{Enabled: false},
					K8sHpaUID:        ResourceAttributeConfig{Enabled: false},
					K8sNamespaceName: ResourceAttributeConfig{Enabled: false},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := loadMetricsBuilderConfig(t, tt.name)
			if diff := cmp.Diff(tt.want, cfg, cmpopts.IgnoreUnexported(MetricConfig{}, ResourceAttributeConfig{})); diff != "" {
				t.Errorf("Config mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}

func loadMetricsBuilderConfig(t *testing.T, name string) MetricsBuilderConfig {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	sub, err := cm.Sub(name)
	require.NoError(t, err)
	cfg := DefaultMetricsBuilderConfig()
	require.NoError(t, component.UnmarshalConfig(sub, &cfg))
	return cfg
}
