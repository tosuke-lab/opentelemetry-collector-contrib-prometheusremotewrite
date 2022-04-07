// Copyright Sabre

package prometheusremotewritereceiver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage/remote"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/pdata"
	"go.opentelemetry.io/collector/obsreport"
)

const (
	receiverFormat = "protobuf"
)

// Receiver - remote write
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

// NewPrometheusRemoteWriteReceiver - remote write
func NewPrometheusRemoteWriteReceiver(config *Config, consumer consumer.Metrics, params component.ReceiverCreateSettings) (*PrometheusRemoteWriteReceiver, error) {
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
	var err = componenterror.ErrNilNextConsumer
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
	pms := pdata.NewMetrics()

	for _, ts := range req.Timeseries {

		prm := pdata.NewResourceMetrics()

		metricName := finaName(ts.Labels)
		pm := pdata.NewMetric()
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
			ppoint := pdata.NewNumberDataPoint()
			ppoint.SetDoubleVal(s.Value)
			ppoint.SetTimestamp(pdata.NewTimestampFromTime(time.Unix(0, s.Timestamp*int64(time.Millisecond))))
			for _, l := range ts.Labels {
				ppoint.Attributes().Insert(l.Name, pdata.NewAttributeValueString(l.Value))
			}
			if IsValidCumulativeSuffix(metricsType) {
				pm.SetDataType(pdata.MetricDataTypeSum)
				pm.Sum().SetIsMonotonic(true)
				pm.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
				ppoint.CopyTo(pm.Sum().DataPoints().AppendEmpty())
			} else {
				pm.SetDataType(pdata.MetricDataTypeGauge)
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
		pilm := pdata.NewInstrumentationLibraryMetrics()
		pm.CopyTo(pilm.Metrics().AppendEmpty())
		pilm.CopyTo(prm.InstrumentationLibraryMetrics().AppendEmpty())
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
	var err = componenterror.ErrNilNextConsumer
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
