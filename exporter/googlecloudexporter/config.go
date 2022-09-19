// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package googlecloudexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter"

import (
	"fmt"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for Google Cloud exporter.
type Config struct {
	config.ExporterSettings `mapstructure:",squash"`
	collector.Config        `mapstructure:",squash"`

	// Timeout for all API calls. If not set, defaults to 12 seconds.
	exporterhelper.TimeoutSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`

	// ClientOption for authentication via API key
	CredentialFileName string `mapstructure:"credential_file_name"`

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
	if err := cfg.ExporterSettings.Validate(); err != nil {
		return fmt.Errorf("exporter settings are invalid :%w", err)
	}
	if err := collector.ValidateConfig(cfg.Config); err != nil {
		return fmt.Errorf("googlecloud exporter settings are invalid :%w", err)
	}
	return nil
}
