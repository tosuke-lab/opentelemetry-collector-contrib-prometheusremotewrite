// Copyright Sabre

package prometheusremotewritereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusremotewritereceiver"

import (
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
)

// Config - remote write
type Config struct {
	config.ReceiverSettings       `mapstructure:",squash"`
	confighttp.HTTPServerSettings `mapstructure:",squash"`
}

var _ config.Receiver = (*Config)(nil)
