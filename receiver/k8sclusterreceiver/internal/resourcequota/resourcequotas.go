// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcequota // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver/internal/resourcequota"

import (
	"strings"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	corev1 "k8s.io/api/core/v1"

	imetadata "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver/internal/resourcequota/internal/metadata"
)

func GetMetrics(set receiver.CreateSettings, rq *corev1.ResourceQuota) pmetric.Metrics {
	mb := imetadata.NewMetricsBuilder(imetadata.DefaultMetricsBuilderConfig(), set)
	ts := pcommon.NewTimestampFromTime(time.Now())

	for k, v := range rq.Status.Hard {
		val := v.Value()
		if strings.HasSuffix(string(k), ".cpu") {
			val = v.MilliValue()
		}
		mb.RecordK8sResourceQuotaHardLimitDataPoint(ts, val, string(k))
	}

	for k, v := range rq.Status.Used {
		val := v.Value()
		if strings.HasSuffix(string(k), ".cpu") {
			val = v.MilliValue()
		}
		mb.RecordK8sResourceQuotaUsedDataPoint(ts, val, string(k))
	}

	return mb.Emit(imetadata.WithK8sResourcequotaUID(string(rq.UID)), imetadata.WithK8sResourcequotaName(rq.Name), imetadata.WithK8sNamespaceName(rq.Namespace), imetadata.WithOpencensusResourcetype("k8s"))
}
