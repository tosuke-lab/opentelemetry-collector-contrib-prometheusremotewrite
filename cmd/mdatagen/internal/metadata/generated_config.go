// Code generated by mdatagen. DO NOT EDIT.

package metadata

import "go.opentelemetry.io/collector/confmap"

// MetricConfig provides common config for a particular metric.
type MetricConfig struct {
	Enabled bool `mapstructure:"enabled"`

	enabledSetByUser bool
}

func (ms *MetricConfig) Unmarshal(parser *confmap.Conf) error {
	if parser == nil {
		return nil
	}
	err := parser.Unmarshal(ms, confmap.WithErrorUnused())
	if err != nil {
		return err
	}
	ms.enabledSetByUser = parser.IsSet("enabled")
	return nil
}

// MetricsConfig provides config for file metrics.
type MetricsConfig struct {
	DefaultMetric            MetricConfig `mapstructure:"default.metric"`
	DefaultMetricToBeRemoved MetricConfig `mapstructure:"default.metric.to_be_removed"`
	OptionalMetric           MetricConfig `mapstructure:"optional.metric"`
}

func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		DefaultMetric: MetricConfig{
			Enabled: true,
		},
		DefaultMetricToBeRemoved: MetricConfig{
			Enabled: true,
		},
		OptionalMetric: MetricConfig{
			Enabled: false,
		},
	}
}

// ResourceAttributeConfig provides common config for a particular resource attribute.
type ResourceAttributeConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// ResourceAttributesConfig provides config for file resource attributes.
type ResourceAttributesConfig struct {
	MapResourceAttr        ResourceAttributeConfig `mapstructure:"map.resource.attr"`
	OptionalResourceAttr   ResourceAttributeConfig `mapstructure:"optional.resource.attr"`
	SliceResourceAttr      ResourceAttributeConfig `mapstructure:"slice.resource.attr"`
	StringEnumResourceAttr ResourceAttributeConfig `mapstructure:"string.enum.resource.attr"`
	StringResourceAttr     ResourceAttributeConfig `mapstructure:"string.resource.attr"`
}

func DefaultResourceAttributesConfig() ResourceAttributesConfig {
	return ResourceAttributesConfig{
		MapResourceAttr: ResourceAttributeConfig{
			Enabled: true,
		},
		OptionalResourceAttr: ResourceAttributeConfig{
			Enabled: false,
		},
		SliceResourceAttr: ResourceAttributeConfig{
			Enabled: true,
		},
		StringEnumResourceAttr: ResourceAttributeConfig{
			Enabled: true,
		},
		StringResourceAttr: ResourceAttributeConfig{
			Enabled: true,
		},
	}
}

// MetricsBuilderConfig is a configuration for file metrics builder.
type MetricsBuilderConfig struct {
	Metrics            MetricsConfig            `mapstructure:"metrics"`
	ResourceAttributes ResourceAttributesConfig `mapstructure:"resource_attributes"`
}

func DefaultMetricsBuilderConfig() MetricsBuilderConfig {
	return MetricsBuilderConfig{
		Metrics:            DefaultMetricsConfig(),
		ResourceAttributes: DefaultResourceAttributesConfig(),
	}
}
