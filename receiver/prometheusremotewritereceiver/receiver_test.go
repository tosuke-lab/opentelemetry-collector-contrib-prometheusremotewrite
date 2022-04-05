// Copyright Sabre

package prometheusremotewritereceiver

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"net/http"
	"testing"
)

func TestPrometheusRemoteWriteReceiver(t *testing.T) {
	ctx := context.Background()
	cms := new(consumertest.MetricsSink)

	receiver, _ := NewPrometheusRemoteWriteReceiver(
		&Config{
			ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(typeStr)),
			HTTPServerSettings: confighttp.HTTPServerSettings{
				Endpoint: defaultBindEndpoint,
			},
		},
		cms,
		componenttest.NewNopReceiverCreateSettings())

	require.NoError(t, receiver.Start(ctx, componenttest.NewNopHost()))

	url := "http://localhost:19291"

	cfg := NewConfig(
		WriteURLOption(url),
	)
	c, err := NewClient(cfg)
	require.NoError(t, err)

	tsList := TSList{
		{
			Labels: []Label{
				{
					Name:  "__name__",
					Value: "foo_bar",
				},
				{
					Name:  "biz",
					Value: "baz",
				},
			},
			Datapoint: Datapoint{
				Timestamp: now,
				Value:     1415.92,
			},
		},
	}

	r, writeErr := c.WriteTimeSeries(context.Background(), tsList, WriteOptions{})
	require.Nil(t, writeErr)
	require.Equal(t, http.StatusAccepted, r.StatusCode)
}
