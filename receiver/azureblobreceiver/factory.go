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

package azureblobreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureblobreceiver"

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent"
)

const (
	// The value of "type" key in configuration.
	typeStr             = "azureblob"
	logsContainerName   = "logs"
	tracesContainerName = "traces"
)

var (
	errUnexpectedConfigurationType = errors.New("failed to cast configuration to Azure Blob Config")
)

type blobReceiverFactory struct {
	receivers *sharedcomponent.SharedComponents
}

// NewFactory returns a factory for Azure Blob receiver.
func NewFactory() component.ReceiverFactory {
	f := &blobReceiverFactory{
		receivers: sharedcomponent.NewSharedComponents(),
	}

	return component.NewReceiverFactory(
		typeStr,
		f.createDefaultConfig,
		component.WithTracesReceiver(f.createTracesReceiver, component.StabilityLevelBeta),
		component.WithLogsReceiver(f.createLogsReceiver, component.StabilityLevelBeta))
}

func (f *blobReceiverFactory) createDefaultConfig() component.ReceiverConfig {
	return &Config{
		ReceiverSettings: config.NewReceiverSettings(component.NewID(typeStr)),
		Logs:             LogsConfig{ContainerName: logsContainerName},
		Traces:           TracesConfig{ContainerName: tracesContainerName},
	}
}

func (f *blobReceiverFactory) createLogsReceiver(
	ctx context.Context,
	set component.ReceiverCreateSettings,
	cfg component.ReceiverConfig,
	nextConsumer consumer.Logs,
) (component.LogsReceiver, error) {

	receiver, err := f.getReceiver(set, cfg)

	if err != nil {
		set.Logger.Error(err.Error())
		return nil, err
	}

	receiver.(logsDataConsumer).setNextLogsConsumer(nextConsumer)

	return receiver, nil
}

func (f *blobReceiverFactory) createTracesReceiver(
	ctx context.Context,
	set component.ReceiverCreateSettings,
	cfg component.ReceiverConfig,
	nextConsumer consumer.Traces,
) (component.TracesReceiver, error) {

	receiver, err := f.getReceiver(set, cfg)

	if err != nil {
		set.Logger.Error(err.Error())
		return nil, err
	}

	receiver.(tracesDataConsumer).setNextTracesConsumer(nextConsumer)
	return receiver, nil
}

func (f *blobReceiverFactory) getReceiver(
	set component.ReceiverCreateSettings,
	cfg component.ReceiverConfig) (component.Component, error) {

	var err error
	r := f.receivers.GetOrAdd(cfg, func() component.Component {
		receiverConfig, ok := cfg.(*Config)

		if !ok {
			err = errUnexpectedConfigurationType
			return nil
		}

		var beh blobEventHandler
		beh, err = f.getBlobEventHandler(receiverConfig, set.Logger)
		if err != nil {
			return nil
		}

		var receiver component.Receiver
		receiver, err = newReceiver(*receiverConfig, set, beh)
		return receiver
	})

	if err != nil {
		return nil, err
	}

	return r.Unwrap(), err
}

func (f *blobReceiverFactory) getBlobEventHandler(cfg *Config, logger *zap.Logger) (blobEventHandler, error) {
	bc, err := newBlobClient(cfg.ConnectionString, logger)
	if err != nil {
		return nil, err
	}

	return newBlobEventHandler(cfg.EventHub.EndPoint, cfg.Logs.ContainerName, cfg.Traces.ContainerName, bc, logger),
		nil
}
