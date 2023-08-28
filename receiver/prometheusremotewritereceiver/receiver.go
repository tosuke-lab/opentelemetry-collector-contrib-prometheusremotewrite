// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheusremotewritereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusremotewritereceiver"

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/prometheus/prometheus/storage/remote"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite"
)

const (
	receiverFormat    = "protobuf"
	receiverTransport = "http"
)

// PrometheusRemoteWriteReceiver - remote write
type PrometheusRemoteWriteReceiver struct {
	params       receiver.CreateSettings
	host         component.Host
	nextConsumer consumer.Metrics

	mu         sync.Mutex
	startOnce  sync.Once
	cancelFunc context.CancelFunc
	shutdownWG sync.WaitGroup

	server        *http.Server
	config        *Config
	timeThreshold *int64
	logger        *zap.Logger
	obsrecv       *obsreport.Receiver
}

// NewPrometheusRemoteWriteReceiver - remote write
func NewPrometheusRemoteWriteReceiver(params receiver.CreateSettings, config *Config, consumer consumer.Metrics) (*PrometheusRemoteWriteReceiver, error) {
	obsrecv, err := obsreport.NewReceiver(obsreport.ReceiverSettings{
		ReceiverID:             params.ID,
		Transport:              receiverTransport,
		LongLivedCtx:           true,
		ReceiverCreateSettings: params,
	})
	zr := &PrometheusRemoteWriteReceiver{
		params:        params,
		nextConsumer:  consumer,
		config:        config,
		logger:        params.Logger,
		obsrecv:       obsrecv,
		timeThreshold: &config.TimeThreshold,
	}
	return zr, err
}

// Start - remote write
func (rec *PrometheusRemoteWriteReceiver) Start(_ context.Context, host component.Host) error {
	_, cancel := context.WithCancel(context.Background())
	rec.cancelFunc = cancel
	if host == nil {
		return errors.New("nil host")
	}
	rec.mu.Lock()
	defer rec.mu.Unlock()
	var err = component.ErrNilNextConsumer
	rec.startOnce.Do(func() {
		err = nil
		rec.host = host
		rec.server, err = rec.config.HTTPServerSettings.ToServer(host, rec.params.TelemetrySettings, rec, confighttp.WithDecoder("snappy", func(body io.ReadCloser) (io.ReadCloser, error) { return body, nil }))
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

	ctx := rec.obsrecv.StartMetricsOp(r.Context())
	req, err := remote.DecodeWriteRequest(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pms, err := prometheusremotewrite.FromTimeSeries(req.Timeseries, prometheusremotewrite.Settings{
		TimeThreshold: *rec.timeThreshold,
		Logger:        *rec.logger,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricCount := pms.ResourceMetrics().Len()
	dataPointCount := pms.DataPointCount()
	if metricCount != 0 {
		err = rec.nextConsumer.ConsumeMetrics(ctx, pms)
	}
	rec.obsrecv.EndMetricsOp(ctx, receiverFormat, dataPointCount, err)
	w.WriteHeader(http.StatusAccepted)
}

// Shutdown tells the receiver that should stop reception,
// giving it a chance to perform any necessary clean-up and
// shutting down its HTTP server.
func (rec *PrometheusRemoteWriteReceiver) Shutdown(context.Context) error {
	if rec.server == nil {
		return nil
	}
	err := rec.server.Close()
	rec.shutdownWG.Wait()
	return err
}
