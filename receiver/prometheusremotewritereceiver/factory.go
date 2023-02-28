// Copyright  The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheusremotewritereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusremotewritereceiver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/receiver"
)

const (
	typeStr              = "prometheusremotewrite"
	defaultBindEndpoint  = "0.0.0.0:19291"
	stability            = component.StabilityLevelDevelopment
	defaultTimeThreshold = 24
)

// NewFactory - remote write
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		newDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, stability))
}

func newDefaultConfig() component.Config {
	return &Config{
		ReceiverSettings: obsreport.ReceiverSettings{
			ReceiverID:   component.NewID(typeStr),
			LongLivedCtx: false,
		},
		HTTPServerSettings: confighttp.HTTPServerSettings{
			Endpoint: defaultBindEndpoint,
		},
		TimeThreshold: defaultTimeThreshold,
	}
}

// createMetricsReceiver creates a metrics receiver based on provided config.
func createMetricsReceiver(
	_ context.Context,
	params receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	rCfg := cfg.(*Config)
	return NewPrometheusRemoteWriteReceiver(params, rCfg, consumer)
}
