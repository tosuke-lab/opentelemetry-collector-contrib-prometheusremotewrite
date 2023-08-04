// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheusremotewritereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusremotewritereceiver"

import (
	"go.opentelemetry.io/collector/config/confighttp"
)

// Config - remote write
type Config struct {
	confighttp.HTTPServerSettings `mapstructure:",squash"`
	TimeThreshold                 int64 `mapstructure:"time_threshold"`
}

// Validate checks the receiver configuration is valid
func (cfg *Config) Validate() error {
	return nil
}
