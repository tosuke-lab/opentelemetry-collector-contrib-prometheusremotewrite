// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheusremotewrite // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite"

import "go.uber.org/zap"

type Settings struct {
	Namespace           string
	ExternalLabels      map[string]string
	DisableTargetInfo   bool
	TimeThreshold       int64
	ExportCreatedMetric bool
	Logger              zap.Logger
	AddMetricSuffixes   bool
}
