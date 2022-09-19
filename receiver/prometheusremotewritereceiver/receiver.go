// Copyright Sabre

package prometheusremotewritereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusremotewritereceiver"

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"net"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage/remote"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
)

const (
	receiverFormat = "protobuf"
)

// PrometheusRemoteWriteReceiver - remote write
type PrometheusRemoteWriteReceiver struct {
	params       component.ReceiverCreateSettings
	host         component.Host
	nextConsumer consumer.Metrics

	mu         sync.Mutex
	startOnce  sync.Once
	stopOnce   sync.Once
	shutdownWG sync.WaitGroup

	server  *http.Server
	config  *Config
	logger  *zap.Logger
	obsrecv *obsreport.Receiver
}

// NewReceiver - remote write
func NewReceiver(params component.ReceiverCreateSettings, config *Config, consumer consumer.Metrics) (*PrometheusRemoteWriteReceiver, error) {
	zr := &PrometheusRemoteWriteReceiver{
		params:       params,
		nextConsumer: consumer,
		config:       config,
		logger:       params.Logger,
		obsrecv: obsreport.NewReceiver(obsreport.ReceiverSettings{
			ReceiverID:             config.ID(),
			ReceiverCreateSettings: params,
		}),
	}
	return zr, nil
}

// Start - remote write
func (rec *PrometheusRemoteWriteReceiver) Start(ctx context.Context, host component.Host) error {
	if host == nil {
		return errors.New("nil host")
	}
	rec.mu.Lock()
	defer rec.mu.Unlock()
	var err = component.ErrNilNextConsumer
	rec.startOnce.Do(func() {
		err = nil
		rec.host = host
		rec.server, err = rec.config.HTTPServerSettings.ToServer(host, rec.params.TelemetrySettings, rec)
		var listener net.Listener
		listener, err = rec.config.HTTPServerSettings.ToListener()
		if err != nil {
			return
		}
		rec.shutdownWG.Add(1)
		go func() {
			defer rec.shutdownWG.Done()
			if errHTTP := rec.server.Serve(listener); errHTTP != http.ErrServerClosed {
				host.ReportFatalError(errHTTP)
			}
		}()
	})
	return err
}

func (rec *PrometheusRemoteWriteReceiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	req, err := remote.DecodeWriteRequest(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := rec.obsrecv.StartMetricsOp(context.Background())

	reg := regexp.MustCompile(`(\w+)_(\w+)_(\w+)\z`)
	timeThreshold := time.Now().Add(-time.Hour * 24)
	pms := pmetric.NewMetrics()

	for _, ts := range req.Timeseries {

		prm := pmetric.NewResourceMetrics()

		metricName := finaName(ts.Labels)
		pm := pmetric.NewMetric()
		pm.SetName(metricName)
		rec.logger.Debug("Metric name", zap.String("metric_name", pm.Name()))

		match := reg.FindStringSubmatch(metricName)
		metricsType := ""
		if len(match) > 1 {
			lastSuffixInMetricName := match[len(match)-1]
			if IsValidSuffix(lastSuffixInMetricName) {
				metricsType = lastSuffixInMetricName
				if len(match) > 2 {
					secondSuffixInMetricName := match[len(match)-2]
					if IsValidUnit(secondSuffixInMetricName) {
						pm.SetUnit(secondSuffixInMetricName)
					}
				}
			} else if IsValidUnit(lastSuffixInMetricName) {
				pm.SetUnit(lastSuffixInMetricName)
			}
		}
		rec.logger.Debug("Metric unit", zap.String("metric name", pm.Name()), zap.String("metric_unit", pm.Unit()))

		for _, s := range ts.Samples {
			ppoint := pmetric.NewNumberDataPoint()
			ppoint.SetDoubleVal(s.Value)
			ppoint.SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, s.Timestamp*int64(time.Millisecond))))
			if ppoint.Timestamp().AsTime().Before(timeThreshold) {
				rec.logger.Debug("Metric older than the threshold", zap.String("metric name", pm.Name()), zap.Time("metric_timestamp", ppoint.Timestamp().AsTime()))
				continue
			}
			for _, l := range ts.Labels {
				labelName := l.Name
				if l.Name == "__name__" {
					labelName = "key_name"
				}
				ppoint.Attributes().Insert(labelName, pcommon.NewValueString(l.Value))
			}
			if IsValidCumulativeSuffix(metricsType) {
				pm.SetDataType(pmetric.MetricDataTypeSum)
				pm.Sum().SetIsMonotonic(true)
				pm.Sum().SetAggregationTemporality(pmetric.MetricAggregationTemporalityCumulative)
				ppoint.CopyTo(pm.Sum().DataPoints().AppendEmpty())
			} else {
				pm.SetDataType(pmetric.MetricDataTypeGauge)
				ppoint.CopyTo(pm.Gauge().DataPoints().AppendEmpty())
			}
			rec.logger.Debug("Metric sample",
				zap.String("metric_name", pm.Name()),
				zap.String("metric_unit", pm.Unit()),
				zap.Float64("metric_value", ppoint.DoubleVal()),
				zap.Time("metric_timestamp", ppoint.Timestamp().AsTime()),
				zap.String("metric_labels", fmt.Sprintf("%#v", ppoint.Attributes())),
			)
		}
		pilm := pmetric.NewScopeMetrics()
		pm.CopyTo(pilm.Metrics().AppendEmpty())
		pilm.CopyTo(prm.ScopeMetrics().AppendEmpty())
		prm.CopyTo(pms.ResourceMetrics().AppendEmpty())
	}

	metricCount := pms.ResourceMetrics().Len()
	dataPointCount := pms.DataPointCount()
	if metricCount != 0 {
		err = rec.nextConsumer.ConsumeMetrics(ctx, pms)
	}
	rec.obsrecv.EndMetricsOp(ctx, receiverFormat, dataPointCount, err)
	w.WriteHeader(http.StatusAccepted)
}

// Shutdown - remote write
func (rec *PrometheusRemoteWriteReceiver) Shutdown(context.Context) error {
	var err = component.ErrNilNextConsumer
	rec.stopOnce.Do(func() {
		err = rec.server.Close()
		rec.shutdownWG.Wait()
	})
	return err
}

func finaName(labels []prompb.Label) (ret string) {
	for _, label := range labels {
		if label.Name == "__name__" {
			return label.Value
		}
	}
	return "error"
}

// IsValidSuffix - remote write
func IsValidSuffix(suffix string) bool {
	switch suffix {
	case
		"max",
		"sum",
		"count",
		"total":
		return true
	}
	return false
}

// IsValidCumulativeSuffix - remote write
func IsValidCumulativeSuffix(suffix string) bool {
	switch suffix {
	case
		"sum",
		"count",
		"total":
		return true
	}
	return false
}

// IsValidUnit - remote write
func IsValidUnit(unit string) bool {
	switch unit {
	case
		"seconds",
		"bytes":
		return true
	}
	return false
}
