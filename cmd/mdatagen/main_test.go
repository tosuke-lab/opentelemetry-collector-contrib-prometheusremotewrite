// Copyright The OpenTelemetry Authors
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

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	validMetadata = `
name: metricreceiver
attributes:
  cpu_type:
    value: type
    description: The type of CPU consumption
    type: string
    enum:
    - user
    - io_wait
    - system
  host:
    description: The type of CPU consumption
    type: string
metrics:
  system.cpu.time:
    enabled: true
    description: Total CPU seconds broken down by different states.
    extended_documentation: Additional information on CPU Time can be found [here](https://en.wikipedia.org/wiki/CPU_time).
    unit: s
    sum:
      aggregation: cumulative
      value_type: double
    attributes: [host, cpu_type]
`
)

func Test_runContents(t *testing.T) {
	type args struct {
		yml string
	}
	tests := []struct {
		name                  string
		args                  args
		expectedDocumentation string
		want                  string
		wantErr               bool
	}{
		{
			name:                  "valid metadata",
			args:                  args{validMetadata},
			expectedDocumentation: "testdata/documentation.md",
			want:                  "",
		},
		{
			name:    "invalid yaml",
			args:    args{"invalid"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir := t.TempDir()

			metadataFile := filepath.Join(tmpdir, "metadata.yaml")
			require.NoError(t, os.WriteFile(metadataFile, []byte(tt.args.yml), 0600))

			err := run(metadataFile)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.FileExists(t, filepath.Join(tmpdir, "internal/metadata/generated_metrics.go"))

			actualDocumentation := filepath.Join(tmpdir, "documentation.md")
			require.FileExists(t, actualDocumentation)
			if tt.expectedDocumentation != "" {
				expectedFileBytes, err := os.ReadFile(tt.expectedDocumentation)
				require.NoError(t, err)

				actualFileBytes, err := os.ReadFile(actualDocumentation)
				require.NoError(t, err)

				require.Equal(t, expectedFileBytes, actualFileBytes)
			}
		})
	}
}

func Test_run(t *testing.T) {
	type args struct {
		ymlPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "no argument",
			args:    args{""},
			wantErr: true,
		},
		{
			name:    "no such file",
			args:    args{"/no/such/file"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args.ymlPath); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
