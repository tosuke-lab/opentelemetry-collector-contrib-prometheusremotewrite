// Copyright The OpenTelemetry Authors
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

package filterprocessor

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/goldendataset"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter/filtermetric"
)

const filteredMetric = "p0_metric_1"
const filteredAttrKey = "pt-label-key-1"

var filteredAttrVal = pcommon.NewValueStr("pt-label-val-1")

func TestExprError(t *testing.T) {
	testMatchError(t, pmetric.MetricTypeGauge, pmetric.NumberDataPointValueTypeInt)
	testMatchError(t, pmetric.MetricTypeGauge, pmetric.NumberDataPointValueTypeDouble)
	testMatchError(t, pmetric.MetricTypeSum, pmetric.NumberDataPointValueTypeInt)
	testMatchError(t, pmetric.MetricTypeSum, pmetric.NumberDataPointValueTypeDouble)
	testMatchError(t, pmetric.MetricTypeHistogram, pmetric.NumberDataPointValueTypeEmpty)
}

func testMatchError(t *testing.T, mdType pmetric.MetricType, mvType pmetric.NumberDataPointValueType) {
	// the "foo" expr expression will cause expr Run() to return an error
	proc, next, logs := testProcessor(t, nil, []string{"foo"})
	pdm := testData("", 1, mdType, mvType)
	err := proc.ConsumeMetrics(context.Background(), pdm)
	assert.NoError(t, err)
	// assert that metrics not be filtered as a result
	assert.Equal(t, []pmetric.Metrics{pdm}, next.AllMetrics())
	assert.Equal(t, 1, logs.Len())
	assert.Equal(t, "shouldKeepMetric failed", logs.All()[0].Message)
}

func TestExprProcessor(t *testing.T) {
	testFilter(t, pmetric.MetricTypeGauge, pmetric.NumberDataPointValueTypeInt)
	testFilter(t, pmetric.MetricTypeGauge, pmetric.NumberDataPointValueTypeDouble)
	testFilter(t, pmetric.MetricTypeSum, pmetric.NumberDataPointValueTypeInt)
	testFilter(t, pmetric.MetricTypeSum, pmetric.NumberDataPointValueTypeDouble)
	testFilter(t, pmetric.MetricTypeHistogram, pmetric.NumberDataPointValueTypeEmpty)
}

func testFilter(t *testing.T, mdType pmetric.MetricType, mvType pmetric.NumberDataPointValueType) {
	format := "MetricName == '%s' && Label('%s') == '%s'"
	q := fmt.Sprintf(format, filteredMetric, filteredAttrKey, filteredAttrVal.Str())

	mds := testDataSlice(2, mdType, mvType)
	totMetricCount := 0
	for _, md := range mds {
		totMetricCount += md.MetricCount()
	}
	expectedMetricCount := totMetricCount - 1
	filtered := filterMetrics(t, nil, []string{q}, mds)
	filteredMetricCount := 0
	for _, metrics := range filtered {
		filteredMetricCount += metrics.MetricCount()
		rmsSlice := metrics.ResourceMetrics()
		for i := 0; i < rmsSlice.Len(); i++ {
			rms := rmsSlice.At(i)
			ilms := rms.ScopeMetrics()
			for j := 0; j < ilms.Len(); j++ {
				ilm := ilms.At(j)
				metricSlice := ilm.Metrics()
				for k := 0; k < metricSlice.Len(); k++ {
					metric := metricSlice.At(k)
					if metric.Name() == filteredMetric {
						dt := metric.Type()
						switch dt {
						case pmetric.MetricTypeGauge:
							pts := metric.Gauge().DataPoints()
							for l := 0; l < pts.Len(); l++ {
								assertFiltered(t, pts.At(l).Attributes())
							}
						case pmetric.MetricTypeSum:
							pts := metric.Sum().DataPoints()
							for l := 0; l < pts.Len(); l++ {
								assertFiltered(t, pts.At(l).Attributes())
							}
						case pmetric.MetricTypeHistogram:
							pts := metric.Histogram().DataPoints()
							for l := 0; l < pts.Len(); l++ {
								assertFiltered(t, pts.At(l).Attributes())
							}
						}
					}
				}
			}
		}
	}
	assert.Equal(t, expectedMetricCount, filteredMetricCount)
}

func assertFiltered(t *testing.T, lm pcommon.Map) {
	lm.Range(func(k string, v pcommon.Value) bool {
		if k == filteredAttrKey && v.Equal(filteredAttrVal) {
			assert.Fail(t, "found metric that should have been filtered out")
			return false
		}
		return true
	})
}

func filterMetrics(t *testing.T, include []string, exclude []string, mds []pmetric.Metrics) []pmetric.Metrics {
	proc, next, _ := testProcessor(t, include, exclude)
	for _, md := range mds {
		err := proc.ConsumeMetrics(context.Background(), md)
		require.NoError(t, err)
	}
	return next.AllMetrics()
}

func testProcessor(t *testing.T, include []string, exclude []string) (component.MetricsProcessor, *consumertest.MetricsSink, *observer.ObservedLogs) {
	factory := NewFactory()
	cfg := exprConfig(factory, include, exclude)
	ctx := context.Background()
	next := &consumertest.MetricsSink{}
	core, logs := observer.New(zapcore.WarnLevel)
	proc, err := factory.CreateMetricsProcessor(
		ctx,
		component.ProcessorCreateSettings{
			TelemetrySettings: component.TelemetrySettings{
				Logger: zap.New(core),
			},
		},
		cfg,
		next,
	)
	require.NoError(t, err)
	require.NotNil(t, proc)
	return proc, next, logs
}

func exprConfig(factory component.ProcessorFactory, include []string, exclude []string) component.ProcessorConfig {
	cfg := factory.CreateDefaultConfig()
	pCfg := cfg.(*Config)
	pCfg.Metrics = MetricFilters{}
	if include != nil {
		pCfg.Metrics.Include = &filtermetric.MatchProperties{
			MatchType:   "expr",
			Expressions: include,
		}
	}
	if exclude != nil {
		pCfg.Metrics.Exclude = &filtermetric.MatchProperties{
			MatchType:   "expr",
			Expressions: exclude,
		}
	}
	return cfg
}

func testDataSlice(size int, mdType pmetric.MetricType, mvType pmetric.NumberDataPointValueType) []pmetric.Metrics {
	var out []pmetric.Metrics
	for i := 0; i < 16; i++ {
		out = append(out, testData(fmt.Sprintf("p%d_", i), size, mdType, mvType))
	}
	return out
}

func testData(prefix string, size int, mdType pmetric.MetricType, mvType pmetric.NumberDataPointValueType) pmetric.Metrics {
	c := goldendataset.MetricsCfg{
		MetricDescriptorType: mdType,
		MetricValueType:      mvType,
		MetricNamePrefix:     prefix,
		NumILMPerResource:    size,
		NumMetricsPerILM:     size,
		NumPtLabels:          size,
		NumPtsPerMetric:      size,
		NumResourceAttrs:     size,
		NumResourceMetrics:   size,
	}
	return goldendataset.MetricsFromCfg(c)
}
