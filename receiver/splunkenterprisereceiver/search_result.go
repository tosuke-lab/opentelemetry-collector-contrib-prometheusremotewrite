// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package splunkenterprisereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/splunkenterprisereceiver"

// metric name and its associated search as a key value pair
var searchDict = map[string]string{
	`SplunkLicenseIndexUsageSearch`: `search=search index=_internal source=*license_usage.log type="Usage"| fields idx, b| eval indexname = if(len(idx)=0 OR isnull(idx),"(UNKNOWN)",idx)| stats sum(b) as b by indexname| eval By=round(b, 9)| fields indexname, By`,
}

var apiDict = map[string]string{
	`SplunkIndexerThroughput`: `/services/server/introspection/indexer?output_mode=json`,
}

type searchResponse struct {
	search string
	Jobid  *string `xml:"sid"`
	Return int
	Fields []*field `xml:"result>field"`
}

type field struct {
	FieldName string `xml:"k,attr"`
	Value     string `xml:"value>text"`
}

// '/services/server/introspection/indexer'
type indexThroughput struct {
	Entries []idxTEntry `json:"entry"`
}

type idxTEntry struct {
	Content idxTContent `json:"content"`
}

type idxTContent struct {
	Status string  `json:"status"`
	AvgKb  float64 `json:"average_KBps"`
}
