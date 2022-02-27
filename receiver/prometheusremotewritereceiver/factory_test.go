// Copyright Sabre

package prometheusremotewritereceiver

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configtest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"testing"
)

func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	assert.Equal(t, "prometheusremotewrite", string(factory.Type()))

	config := factory.CreateDefaultConfig()
	assert.NotNil(t, config, "failed to create default config")
	assert.NoError(t, configtest.CheckConfigStruct(config))
}

func TestCreateReceiver(t *testing.T) {
	factory := NewFactory()
	config := factory.CreateDefaultConfig()

	params := componenttest.NewNopReceiverCreateSettings()
	traceReceiver, err := factory.CreateTracesReceiver(context.Background(), params, config, consumertest.NewNop())
	assert.ErrorIs(t, err, componenterror.ErrDataTypeIsNotSupported)
	assert.Nil(t, traceReceiver)

	metricReceiver, err := factory.CreateMetricsReceiver(context.Background(), params, config, consumertest.NewNop())
	assert.NoError(t, err, "Metric receiver creation failed")
	assert.NotNil(t, metricReceiver, "Receiver creation failed")
}
