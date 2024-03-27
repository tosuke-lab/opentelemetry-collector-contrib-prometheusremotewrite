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

package logicmonitorexporter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	lmsdktraces "github.com/logicmonitor/lm-data-sdk-go/api/traces"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func Test_NewTracesExporter(t *testing.T) {
	t.Run("should create Traces exporter", func(t *testing.T) {
		config := &Config{
			HTTPClientSettings: confighttp.HTTPClientSettings{
				Endpoint: "http://example.logicmonitor.com/rest",
			},
			APIToken: APIToken{AccessID: "testid", AccessKey: "testkey"},
		}
		set := exportertest.NewNopCreateSettings()
		exp := newTracesExporter(context.Background(), config, set)
		assert.NotNil(t, exp)
	})
}

func TestPushTraceData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := lmsdktraces.LMTraceIngestResponse{
			Success: true,
			Message: "",
		}
		assert.NoError(t, json.NewEncoder(w).Encode(&response))
	}))
	defer ts.Close()

	params := exportertest.NewNopCreateSettings()
	f := NewFactory()
	config := &Config{
		HTTPClientSettings: confighttp.HTTPClientSettings{
			Endpoint: ts.URL,
		},
		APIToken: APIToken{AccessID: "testid", AccessKey: "testkey"},
	}
	ctx := context.Background()
	exp, err := f.CreateTracesExporter(ctx, params, config)
	assert.NoError(t, err)
	assert.NoError(t, exp.Start(ctx, componenttest.NewNopHost()))

	testTraces := ptrace.NewTraces()
	generateTraces().CopyTo(testTraces)
	err = exp.ConsumeTraces(context.Background(), testTraces)
	assert.NoError(t, err)
}

func generateTraces() ptrace.Traces {
	traces := ptrace.NewTraces()
	resourceSpans := traces.ResourceSpans()
	rs := resourceSpans.AppendEmpty()
	rs.ScopeSpans().AppendEmpty().Spans().AppendEmpty()
	return traces
}
