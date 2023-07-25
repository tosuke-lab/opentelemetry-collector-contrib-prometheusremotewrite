// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheusremotewritereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusremotewritereceiver"

import (
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/obsreport"
)

// Config - remote write
type Config struct {
	obsreport.ReceiverSettings    `mapstructure:",squash"`
	confighttp.HTTPServerSettings `mapstructure:",squash"`
	TimeThreshold                 int64 `mapstructure:"time_threshold"`
}
