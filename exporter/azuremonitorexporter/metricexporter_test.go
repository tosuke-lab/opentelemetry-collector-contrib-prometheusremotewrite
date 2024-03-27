// Copyright OpenTelemetry Authors
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

package azuremonitorexporter

/*
Contains tests for metricexporter.go and metric_to_envelopes.go
*/

import (
	"context"
	"testing"

	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

// Test onMetricData callback for the test metrics data
func TestExporterMetricDataCallback(t *testing.T) {
	mockTransportChannel := getMockTransportChannel()
	exporter := getMetricExporter(defaultConfig, mockTransportChannel)

	metrics := getTestMetrics()

	assert.NoError(t, exporter.onMetricData(context.Background(), metrics))

	mockTransportChannel.AssertNumberOfCalls(t, "Send", 5)
}

func TestGaugeEnvelopes(t *testing.T) {
	gaugeMetric := getTestGaugeMetric()
	dataPoint := getDataPoint(t, gaugeMetric)

	assert.Equal(t, dataPoint.Name, "Gauge")
	assert.Equal(t, dataPoint.Value, float64(1))
	assert.Equal(t, dataPoint.Count, 1)
	assert.Equal(t, dataPoint.Kind, contracts.Measurement)
}

func TestSumEnvelopes(t *testing.T) {
	sumMetric := getTestSumMetric()
	dataPoint := getDataPoint(t, sumMetric)

	assert.Equal(t, dataPoint.Name, "Sum")
	assert.Equal(t, dataPoint.Value, float64(2))
	assert.Equal(t, dataPoint.Count, 1)
	assert.Equal(t, dataPoint.Kind, contracts.Measurement)
}

func TestHistogramEnvelopes(t *testing.T) {
	histogramMetric := getTestHistogramMetric()
	dataPoint := getDataPoint(t, histogramMetric)

	assert.Equal(t, dataPoint.Name, "Histogram")
	assert.Equal(t, dataPoint.Value, float64(3))
	assert.Equal(t, dataPoint.Count, 3)
	assert.Equal(t, dataPoint.Min, float64(0))
	assert.Equal(t, dataPoint.Max, float64(2))
	assert.Equal(t, dataPoint.Kind, contracts.Aggregation)
}

func TestExponentialHistogramEnvelopes(t *testing.T) {
	exponentialHistogramMetric := getTestExponentialHistogramMetric()
	dataPoint := getDataPoint(t, exponentialHistogramMetric)

	assert.Equal(t, dataPoint.Name, "ExponentialHistogram")
	assert.Equal(t, dataPoint.Value, float64(4))
	assert.Equal(t, dataPoint.Count, 4)
	assert.Equal(t, dataPoint.Min, float64(1))
	assert.Equal(t, dataPoint.Max, float64(3))
	assert.Equal(t, dataPoint.Kind, contracts.Aggregation)
}

func TestSummaryEnvelopes(t *testing.T) {
	summaryMetric := getTestSummaryMetric()
	dataPoint := getDataPoint(t, summaryMetric)

	assert.Equal(t, dataPoint.Name, "Summary")
	assert.Equal(t, dataPoint.Value, float64(5))
	assert.Equal(t, dataPoint.Count, 5)
	assert.Equal(t, dataPoint.Kind, contracts.Aggregation)
}

func getDataPoint(t testing.TB, metric pmetric.Metric) *contracts.DataPoint {
	var envelopes []*contracts.Envelope = getMetricPacker().MetricToEnvelopes(metric, getResource(), getScope())
	require.Equal(t, len(envelopes), 1)
	envelope := envelopes[0]
	require.NotNil(t, envelope)

	assert.NotNil(t, envelope.Tags)
	assert.NotNil(t, envelope.Time)

	require.NotNil(t, envelope.Data)
	envelopeData := envelope.Data.(*contracts.Data)
	assert.Equal(t, envelopeData.BaseType, "MetricData")

	require.NotNil(t, envelopeData.BaseData)

	metricData := envelopeData.BaseData.(*contracts.MetricData)

	require.Equal(t, len(metricData.Metrics), 1)

	dataPoint := metricData.Metrics[0]
	require.NotNil(t, dataPoint)

	return dataPoint
}

func getMetricExporter(config *Config, transportChannel transportChannel) *metricExporter {
	return &metricExporter{
		config,
		transportChannel,
		zap.NewNop(),
		newMetricPacker(zap.NewNop()),
	}
}

func getMetricPacker() *metricPacker {
	return newMetricPacker(zap.NewNop())
}

func getTestMetrics() pmetric.Metrics {
	metrics := pmetric.NewMetrics()
	resourceMetricsSlice := metrics.ResourceMetrics()
	resourceMetric := resourceMetricsSlice.AppendEmpty()
	scopeMetricsSlice := resourceMetric.ScopeMetrics()
	scopeMetrics := scopeMetricsSlice.AppendEmpty()
	metricSlice := scopeMetrics.Metrics()

	metric := metricSlice.AppendEmpty()
	gaugeMetric := getTestGaugeMetric()
	gaugeMetric.CopyTo(metric)

	metric = metricSlice.AppendEmpty()
	sumMetric := getTestSumMetric()
	sumMetric.CopyTo(metric)

	metric = metricSlice.AppendEmpty()
	histogramMetric := getTestHistogramMetric()
	histogramMetric.CopyTo(metric)

	metric = metricSlice.AppendEmpty()
	exponentialHistogramMetric := getTestExponentialHistogramMetric()
	exponentialHistogramMetric.CopyTo(metric)

	metric = metricSlice.AppendEmpty()
	summaryMetric := getTestSummaryMetric()
	summaryMetric.CopyTo(metric)

	return metrics
}

func getTestGaugeMetric() pmetric.Metric {
	metric := pmetric.NewMetric()
	metric.SetName("Gauge")
	metric.SetEmptyGauge()
	datapoints := metric.Gauge().DataPoints()
	datapoint := datapoints.AppendEmpty()
	datapoint.SetDoubleValue(1)
	return metric
}

func getTestSumMetric() pmetric.Metric {
	metric := pmetric.NewMetric()
	metric.SetName("Sum")
	metric.SetEmptySum()
	datapoints := metric.Sum().DataPoints()
	datapoint := datapoints.AppendEmpty()
	datapoint.SetDoubleValue(2)
	return metric
}

func getTestHistogramMetric() pmetric.Metric {
	metric := pmetric.NewMetric()
	metric.SetName("Histogram")
	metric.SetEmptyHistogram()
	datapoints := metric.Histogram().DataPoints()
	datapoint := datapoints.AppendEmpty()
	datapoint.SetSum(3)
	datapoint.SetCount(3)
	datapoint.SetMin(0)
	datapoint.SetMax(2)
	return metric
}

func getTestExponentialHistogramMetric() pmetric.Metric {
	metric := pmetric.NewMetric()
	metric.SetName("ExponentialHistogram")
	metric.SetEmptyExponentialHistogram()
	datapoints := metric.ExponentialHistogram().DataPoints()
	datapoint := datapoints.AppendEmpty()
	datapoint.SetSum(4)
	datapoint.SetCount(4)
	datapoint.SetMin(1)
	datapoint.SetMax(3)
	return metric
}

func getTestSummaryMetric() pmetric.Metric {
	metric := pmetric.NewMetric()
	metric.SetName("Summary")
	metric.SetEmptySummary()
	datapoints := metric.Summary().DataPoints()
	datapoint := datapoints.AppendEmpty()
	datapoint.SetSum(5)
	datapoint.SetCount(5)
	return metric
}
