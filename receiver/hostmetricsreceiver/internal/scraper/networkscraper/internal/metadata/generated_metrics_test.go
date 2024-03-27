// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestDefaultMetrics(t *testing.T) {
	start := pcommon.Timestamp(1_000_000_000)
	ts := pcommon.Timestamp(1_000_001_000)
	mb := NewMetricsBuilder(DefaultMetricsSettings(), component.BuildInfo{}, WithStartTime(start))
	enabledMetrics := make(map[string]bool)

	enabledMetrics["system.network.connections"] = true
	mb.RecordSystemNetworkConnectionsDataPoint(ts, 1, AttributeProtocol(1), "attr-val")

	mb.RecordSystemNetworkConntrackCountDataPoint(ts, 1)

	mb.RecordSystemNetworkConntrackMaxDataPoint(ts, 1)

	enabledMetrics["system.network.dropped"] = true
	mb.RecordSystemNetworkDroppedDataPoint(ts, 1, "attr-val", AttributeDirection(1))

	enabledMetrics["system.network.errors"] = true
	mb.RecordSystemNetworkErrorsDataPoint(ts, 1, "attr-val", AttributeDirection(1))

	enabledMetrics["system.network.io"] = true
	mb.RecordSystemNetworkIoDataPoint(ts, 1, "attr-val", AttributeDirection(1))

	enabledMetrics["system.network.packets"] = true
	mb.RecordSystemNetworkPacketsDataPoint(ts, 1, "attr-val", AttributeDirection(1))

	metrics := mb.Emit()

	assert.Equal(t, 1, metrics.ResourceMetrics().Len())
	sm := metrics.ResourceMetrics().At(0).ScopeMetrics()
	assert.Equal(t, 1, sm.Len())
	ms := sm.At(0).Metrics()
	assert.Equal(t, len(enabledMetrics), ms.Len())
	seenMetrics := make(map[string]bool)
	for i := 0; i < ms.Len(); i++ {
		assert.True(t, enabledMetrics[ms.At(i).Name()])
		seenMetrics[ms.At(i).Name()] = true
	}
	assert.Equal(t, len(enabledMetrics), len(seenMetrics))
}

func TestAllMetrics(t *testing.T) {
	start := pcommon.Timestamp(1_000_000_000)
	ts := pcommon.Timestamp(1_000_001_000)
	settings := MetricsSettings{
		SystemNetworkConnections:    MetricSettings{Enabled: true},
		SystemNetworkConntrackCount: MetricSettings{Enabled: true},
		SystemNetworkConntrackMax:   MetricSettings{Enabled: true},
		SystemNetworkDropped:        MetricSettings{Enabled: true},
		SystemNetworkErrors:         MetricSettings{Enabled: true},
		SystemNetworkIo:             MetricSettings{Enabled: true},
		SystemNetworkPackets:        MetricSettings{Enabled: true},
	}
	mb := NewMetricsBuilder(settings, component.BuildInfo{}, WithStartTime(start))

	mb.RecordSystemNetworkConnectionsDataPoint(ts, 1, AttributeProtocol(1), "attr-val")
	mb.RecordSystemNetworkConntrackCountDataPoint(ts, 1)
	mb.RecordSystemNetworkConntrackMaxDataPoint(ts, 1)
	mb.RecordSystemNetworkDroppedDataPoint(ts, 1, "attr-val", AttributeDirection(1))
	mb.RecordSystemNetworkErrorsDataPoint(ts, 1, "attr-val", AttributeDirection(1))
	mb.RecordSystemNetworkIoDataPoint(ts, 1, "attr-val", AttributeDirection(1))
	mb.RecordSystemNetworkPacketsDataPoint(ts, 1, "attr-val", AttributeDirection(1))

	metrics := mb.Emit()

	assert.Equal(t, 1, metrics.ResourceMetrics().Len())
	rm := metrics.ResourceMetrics().At(0)
	attrCount := 0
	assert.Equal(t, attrCount, rm.Resource().Attributes().Len())

	assert.Equal(t, 1, rm.ScopeMetrics().Len())
	ms := rm.ScopeMetrics().At(0).Metrics()
	allMetricsCount := reflect.TypeOf(MetricsSettings{}).NumField()
	assert.Equal(t, allMetricsCount, ms.Len())
	validatedMetrics := make(map[string]struct{})
	for i := 0; i < ms.Len(); i++ {
		switch ms.At(i).Name() {
		case "system.network.connections":
			assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
			assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
			assert.Equal(t, "The number of connections.", ms.At(i).Description())
			assert.Equal(t, "{connections}", ms.At(i).Unit())
			assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
			assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
			dp := ms.At(i).Sum().DataPoints().At(0)
			assert.Equal(t, start, dp.StartTimestamp())
			assert.Equal(t, ts, dp.Timestamp())
			assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
			assert.Equal(t, int64(1), dp.IntValue())
			attrVal, ok := dp.Attributes().Get("protocol")
			assert.True(t, ok)
			assert.Equal(t, AttributeProtocol(1).String(), attrVal.Str())
			attrVal, ok = dp.Attributes().Get("state")
			assert.True(t, ok)
			assert.EqualValues(t, "attr-val", attrVal.Str())
			validatedMetrics["system.network.connections"] = struct{}{}
		case "system.network.conntrack.count":
			assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
			assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
			assert.Equal(t, "The count of entries in conntrack table.", ms.At(i).Description())
			assert.Equal(t, "{entries}", ms.At(i).Unit())
			assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
			assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
			dp := ms.At(i).Sum().DataPoints().At(0)
			assert.Equal(t, start, dp.StartTimestamp())
			assert.Equal(t, ts, dp.Timestamp())
			assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
			assert.Equal(t, int64(1), dp.IntValue())
			validatedMetrics["system.network.conntrack.count"] = struct{}{}
		case "system.network.conntrack.max":
			assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
			assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
			assert.Equal(t, "The limit for entries in the conntrack table.", ms.At(i).Description())
			assert.Equal(t, "{entries}", ms.At(i).Unit())
			assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
			assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
			dp := ms.At(i).Sum().DataPoints().At(0)
			assert.Equal(t, start, dp.StartTimestamp())
			assert.Equal(t, ts, dp.Timestamp())
			assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
			assert.Equal(t, int64(1), dp.IntValue())
			validatedMetrics["system.network.conntrack.max"] = struct{}{}
		case "system.network.dropped":
			assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
			assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
			assert.Equal(t, "The number of packets dropped.", ms.At(i).Description())
			assert.Equal(t, "{packets}", ms.At(i).Unit())
			assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
			assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
			dp := ms.At(i).Sum().DataPoints().At(0)
			assert.Equal(t, start, dp.StartTimestamp())
			assert.Equal(t, ts, dp.Timestamp())
			assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
			assert.Equal(t, int64(1), dp.IntValue())
			attrVal, ok := dp.Attributes().Get("device")
			assert.True(t, ok)
			assert.EqualValues(t, "attr-val", attrVal.Str())
			attrVal, ok = dp.Attributes().Get("direction")
			assert.True(t, ok)
			assert.Equal(t, AttributeDirection(1).String(), attrVal.Str())
			validatedMetrics["system.network.dropped"] = struct{}{}
		case "system.network.errors":
			assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
			assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
			assert.Equal(t, "The number of errors encountered.", ms.At(i).Description())
			assert.Equal(t, "{errors}", ms.At(i).Unit())
			assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
			assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
			dp := ms.At(i).Sum().DataPoints().At(0)
			assert.Equal(t, start, dp.StartTimestamp())
			assert.Equal(t, ts, dp.Timestamp())
			assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
			assert.Equal(t, int64(1), dp.IntValue())
			attrVal, ok := dp.Attributes().Get("device")
			assert.True(t, ok)
			assert.EqualValues(t, "attr-val", attrVal.Str())
			attrVal, ok = dp.Attributes().Get("direction")
			assert.True(t, ok)
			assert.Equal(t, AttributeDirection(1).String(), attrVal.Str())
			validatedMetrics["system.network.errors"] = struct{}{}
		case "system.network.io":
			assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
			assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
			assert.Equal(t, "The number of bytes transmitted and received.", ms.At(i).Description())
			assert.Equal(t, "By", ms.At(i).Unit())
			assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
			assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
			dp := ms.At(i).Sum().DataPoints().At(0)
			assert.Equal(t, start, dp.StartTimestamp())
			assert.Equal(t, ts, dp.Timestamp())
			assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
			assert.Equal(t, int64(1), dp.IntValue())
			attrVal, ok := dp.Attributes().Get("device")
			assert.True(t, ok)
			assert.EqualValues(t, "attr-val", attrVal.Str())
			attrVal, ok = dp.Attributes().Get("direction")
			assert.True(t, ok)
			assert.Equal(t, AttributeDirection(1).String(), attrVal.Str())
			validatedMetrics["system.network.io"] = struct{}{}
		case "system.network.packets":
			assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
			assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
			assert.Equal(t, "The number of packets transferred.", ms.At(i).Description())
			assert.Equal(t, "{packets}", ms.At(i).Unit())
			assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
			assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
			dp := ms.At(i).Sum().DataPoints().At(0)
			assert.Equal(t, start, dp.StartTimestamp())
			assert.Equal(t, ts, dp.Timestamp())
			assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
			assert.Equal(t, int64(1), dp.IntValue())
			attrVal, ok := dp.Attributes().Get("device")
			assert.True(t, ok)
			assert.EqualValues(t, "attr-val", attrVal.Str())
			attrVal, ok = dp.Attributes().Get("direction")
			assert.True(t, ok)
			assert.Equal(t, AttributeDirection(1).String(), attrVal.Str())
			validatedMetrics["system.network.packets"] = struct{}{}
		}
	}
	assert.Equal(t, allMetricsCount, len(validatedMetrics))
}

func TestNoMetrics(t *testing.T) {
	start := pcommon.Timestamp(1_000_000_000)
	ts := pcommon.Timestamp(1_000_001_000)
	settings := MetricsSettings{
		SystemNetworkConnections:    MetricSettings{Enabled: false},
		SystemNetworkConntrackCount: MetricSettings{Enabled: false},
		SystemNetworkConntrackMax:   MetricSettings{Enabled: false},
		SystemNetworkDropped:        MetricSettings{Enabled: false},
		SystemNetworkErrors:         MetricSettings{Enabled: false},
		SystemNetworkIo:             MetricSettings{Enabled: false},
		SystemNetworkPackets:        MetricSettings{Enabled: false},
	}
	mb := NewMetricsBuilder(settings, component.BuildInfo{}, WithStartTime(start))
	mb.RecordSystemNetworkConnectionsDataPoint(ts, 1, AttributeProtocol(1), "attr-val")
	mb.RecordSystemNetworkConntrackCountDataPoint(ts, 1)
	mb.RecordSystemNetworkConntrackMaxDataPoint(ts, 1)
	mb.RecordSystemNetworkDroppedDataPoint(ts, 1, "attr-val", AttributeDirection(1))
	mb.RecordSystemNetworkErrorsDataPoint(ts, 1, "attr-val", AttributeDirection(1))
	mb.RecordSystemNetworkIoDataPoint(ts, 1, "attr-val", AttributeDirection(1))
	mb.RecordSystemNetworkPacketsDataPoint(ts, 1, "attr-val", AttributeDirection(1))

	metrics := mb.Emit()

	assert.Equal(t, 0, metrics.ResourceMetrics().Len())
}
