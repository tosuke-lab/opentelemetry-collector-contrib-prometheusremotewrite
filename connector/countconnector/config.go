// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package countconnector // import "github.com/open-telemetry/opentelemetry-collector-contrib/connector/countconnector"

const (
	defaultMetricNameSpans = "trace.span.count"
	defaultMetricDescSpans = "The number of spans observed."

	defaultMetricNameDataPoints = "metric.data_point.count"
	defaultMetricDescDataPoints = "The number of data points observed."

	defaultMetricNameLogRecords = "log.record.count"
	defaultMetricDescLogRecords = "The number of log records observed."
)

// MetricInfo for a data type
type MetricInfo struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
}

// Config for the connector
type Config struct {
	Traces  MetricInfo `mapstructure:"traces"`
	Metrics MetricInfo `mapstructure:"metrics"`
	Logs    MetricInfo `mapstructure:"logs"`
}
