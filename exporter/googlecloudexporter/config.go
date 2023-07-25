// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package googlecloudexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter"

import (
	"fmt"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry"
)

// Config defines configuration for Google Cloud exporter.
type Config struct {
	collector.Config `mapstructure:",squash"`

	// Timeout for all API calls. If not set, defaults to 12 seconds.
	exporterhelper.TimeoutSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`

	// ResourceToTelemetrySettings is the option for converting resource attributes to telemetry attributes.
	// "Enabled" - A boolean field to enable/disable this option. Default is `false`.
	// If enabled, all the resource attributes will be converted to metric labels by default.
	ResourceToTelemetrySettings resourcetotelemetry.Settings `mapstructure:"resource_to_telemetry_conversion"`

	// Max number of metric's labels. Metrics exceeding this would be dropped. Default 0 = no limit
	LabelsLimit int `mapstructure:"labels_limit"`

	LabelsToResources []LabelsToResource `mapstructure:"labels_to_resources"`
}

type LabelsToResource struct {
	RequiredLabel string `mapstructure:"required_label"`
	TargetType    string `mapstructure:"target_type"`

	LabelToResources []LabelToResource `mapstructure:"label_to_resource"`
}

type LabelToResource struct {
	SourceLabel         string `mapstructure:"source_label"`
	TargetResourceLabel string `mapstructure:"target_resource_label"`
}

func (cfg *Config) Validate() error {
	if err := collector.ValidateConfig(cfg.Config); err != nil {
		return fmt.Errorf("googlecloud exporter settings are invalid :%w", err)
	}
	return nil
}
