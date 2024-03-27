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

package splunkhecreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/splunkhecreceiver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk"
)

// This file implements factory for Splunk HEC receiver.

const (
	// The value of "type" key in configuration.
	typeStr = "splunk_hec"
	// The stability level of the receiver.
	stability = component.StabilityLevelBeta

	// Default endpoints to bind to.
	defaultEndpoint = ":8088"
)

// NewFactory creates a factory for Splunk HEC receiver.
func NewFactory() component.ReceiverFactory {
	return component.NewReceiverFactory(
		typeStr,
		createDefaultConfig,
		component.WithMetricsReceiver(createMetricsReceiver, stability),
		component.WithLogsReceiver(createLogsReceiver, stability))
}

// CreateDefaultConfig creates the default configuration for Splunk HEC receiver.
func createDefaultConfig() component.ReceiverConfig {
	return &Config{
		ReceiverSettings: config.NewReceiverSettings(component.NewID(typeStr)),
		HTTPServerSettings: confighttp.HTTPServerSettings{
			Endpoint: defaultEndpoint,
		},
		AccessTokenPassthroughConfig: splunk.AccessTokenPassthroughConfig{},
		HecToOtelAttrs: splunk.HecToOtelAttrs{
			Source:     splunk.DefaultSourceLabel,
			SourceType: splunk.DefaultSourceTypeLabel,
			Index:      splunk.DefaultIndexLabel,
			Host:       conventions.AttributeHostName,
		},
		RawPath:    splunk.DefaultRawPath,
		HealthPath: splunk.DefaultHealthPath,
	}
}

// CreateMetricsReceiver creates a metrics receiver based on provided config.
func createMetricsReceiver(
	_ context.Context,
	params component.ReceiverCreateSettings,
	cfg component.ReceiverConfig,
	consumer consumer.Metrics,
) (component.MetricsReceiver, error) {

	rCfg := cfg.(*Config)

	if rCfg.Path != "" {
		params.Logger.Warn("splunk_hec receiver path is deprecated", zap.String("path", rCfg.Path))
	}

	return newMetricsReceiver(params, *rCfg, consumer)
}

// createLogsReceiver creates a logs receiver based on provided config.
func createLogsReceiver(
	_ context.Context,
	params component.ReceiverCreateSettings,
	cfg component.ReceiverConfig,
	consumer consumer.Logs,
) (component.LogsReceiver, error) {

	rCfg := cfg.(*Config)

	if rCfg.Path != "" {
		params.Logger.Warn("splunk_hec receiver path is deprecated", zap.String("path", rCfg.Path))
	}

	return newLogsReceiver(params, *rCfg, consumer)
}
