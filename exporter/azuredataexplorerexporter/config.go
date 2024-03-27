// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package azuredataexplorerexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/azuredataexplorerexporter"

import (
	"errors"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/config/configopaque"
)

// Config defines configuration for Azure Data Explorer Exporter
type Config struct {
	ClusterURI         string              `mapstructure:"cluster_uri"`
	ApplicationID      string              `mapstructure:"application_id"`
	ApplicationKey     configopaque.String `mapstructure:"application_key"`
	TenantID           string              `mapstructure:"tenant_id"`
	Database           string              `mapstructure:"db_name"`
	MetricTable        string              `mapstructure:"metrics_table_name"`
	LogTable           string              `mapstructure:"logs_table_name"`
	TraceTable         string              `mapstructure:"traces_table_name"`
	MetricTableMapping string              `mapstructure:"metrics_table_json_mapping"`
	LogTableMapping    string              `mapstructure:"logs_table_json_mapping"`
	TraceTableMapping  string              `mapstructure:"traces_table_json_mapping"`
	IngestionType      string              `mapstructure:"ingestion_type"`
}

// Validate checks if the exporter configuration is valid
func (adxCfg *Config) Validate() error {
	if adxCfg == nil {
		return errors.New("ADX config is nil / not provided")
	}

	if isEmpty(adxCfg.ClusterURI) || isEmpty(adxCfg.ApplicationID) || isEmpty(string(adxCfg.ApplicationKey)) || isEmpty(adxCfg.TenantID) {
		return errors.New(`mandatory configurations "cluster_uri" ,"application_id" , "application_key" and "tenant_id" are missing or empty `)
	}

	if !(adxCfg.IngestionType == managedIngestType || adxCfg.IngestionType == queuedIngestTest || isEmpty(adxCfg.IngestionType)) {
		return fmt.Errorf("unsupported configuration for ingestion_type. Accepted types [%s, %s] Provided [%s]", managedIngestType, queuedIngestTest, adxCfg.IngestionType)
	}

	return nil
}

func isEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}
