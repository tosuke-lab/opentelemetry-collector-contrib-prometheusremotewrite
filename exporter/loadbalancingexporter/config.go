// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package loadbalancingexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/loadbalancingexporter"

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/servicediscovery/types"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
)

type routingKey int

const (
	traceIDRouting routingKey = iota
	svcRouting
	metricNameRouting
	resourceRouting
)

// Config defines configuration for the exporter.
type Config struct {
	Protocol   Protocol         `mapstructure:"protocol"`
	Resolver   ResolverSettings `mapstructure:"resolver"`
	RoutingKey string           `mapstructure:"routing_key"`
}

// Protocol holds the individual protocol-specific settings. Only OTLP is supported at the moment.
type Protocol struct {
	OTLP otlpexporter.Config `mapstructure:"otlp"`
}

// ResolverSettings defines the configurations for the backend resolver
type ResolverSettings struct {
	Static      *StaticResolver      `mapstructure:"static"`
	DNS         *DNSResolver         `mapstructure:"dns"`
	K8sSvc      *K8sSvcResolver      `mapstructure:"k8s"`
	AWSCloudMap *AWSCloudMapResolver `mapstructure:"awsCloudMap"`
}

// StaticResolver defines the configuration for the resolver providing a fixed list of backends
type StaticResolver struct {
	Hostnames []string `mapstructure:"hostnames"`
}

// DNSResolver defines the configuration for the DNS resolver
type DNSResolver struct {
	Hostname string        `mapstructure:"hostname"`
	Port     string        `mapstructure:"port"`
	Interval time.Duration `mapstructure:"interval"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// K8sSvcResolver defines the configuration for the DNS resolver
type K8sSvcResolver struct {
	Service string  `mapstructure:"service"`
	Ports   []int32 `mapstructure:"ports"`
}

type AWSCloudMapResolver struct {
	NamespaceName string                   `mapstructure:"namespace"`
	ServiceName   string                   `mapstructure:"serviceName"`
	HealthStatus  types.HealthStatusFilter `mapstructure:"healthStatus"`
	Interval      time.Duration            `mapstructure:"interval"`
	Timeout       time.Duration            `mapstructure:"timeout"`
	Port          *uint16                  `mapstructure:"port"`
}
