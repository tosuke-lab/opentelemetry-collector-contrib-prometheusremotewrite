// Copyright Sabre

package prometheusremotewritereceiver

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pmetric"
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
					Value: "http_server_requests_seconds_sum",
				},
				{
					Name:  "method",
					Value: "GET",
				},
			},
			Datapoint: Datapoint{
				Timestamp: now,
				Value:     1.0,
			},
		},
		{
			Labels: []Label{
				{
					Name:  "__name__",
					Value: "http_server_requests_seconds_count",
				},
				{
					Name:  "method",
					Value: "GET",
				},
			},
			Datapoint: Datapoint{
				Timestamp: now,
				Value:     2.0,
			},
		},
		{
			Labels: []Label{
				{
					Name:  "__name__",
					Value: "pdc_jvm_memory_used_bytes",
				},
				{
					Name:  "xxx",
					Value: "yyy",
				},
			},
			Datapoint: Datapoint{
				Timestamp: now,
				Value:     0.0,
			},
		},
	}

	r, writeErr := c.WriteTimeSeries(context.Background(), tsList, WriteOptions{})
	require.Nil(t, writeErr)
	require.Equal(t, http.StatusAccepted, r.StatusCode)

	require.Equal(t, pmetric.MetricDataTypeSum, cms.AllMetrics()[0].ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).DataType())
	require.Equal(t, true, cms.AllMetrics()[0].ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().IsMonotonic())
	require.Equal(t, pmetric.MetricAggregationTemporalityCumulative, cms.AllMetrics()[0].ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().AggregationTemporality())

	require.Equal(t, pmetric.MetricDataTypeSum, cms.AllMetrics()[0].ResourceMetrics().At(1).ScopeMetrics().At(0).Metrics().At(0).DataType())
	require.Equal(t, true, cms.AllMetrics()[0].ResourceMetrics().At(1).ScopeMetrics().At(0).Metrics().At(0).Sum().IsMonotonic())
	require.Equal(t, pmetric.MetricAggregationTemporalityCumulative, cms.AllMetrics()[0].ResourceMetrics().At(1).ScopeMetrics().At(0).Metrics().At(0).Sum().AggregationTemporality())

	require.Equal(t, pmetric.MetricDataTypeGauge, cms.AllMetrics()[0].ResourceMetrics().At(2).ScopeMetrics().At(0).Metrics().At(0).DataType())
}
