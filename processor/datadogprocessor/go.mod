module github.com/open-telemetry/opentelemetry-collector-contrib/processor/datadogprocessor

go 1.19

require (
	github.com/DataDog/datadog-agent/pkg/trace v0.43.0-rc.3.0.20230227093656-ca3af9ce4c3c
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/metrics v0.1.1
	github.com/DataDog/sketches-go v1.4.1
	github.com/stretchr/testify v1.8.2
	go.opentelemetry.io/collector v0.73.0
	go.opentelemetry.io/collector/component v0.73.0
	go.opentelemetry.io/collector/consumer v0.73.0
	go.opentelemetry.io/collector/exporter v0.73.0
	go.opentelemetry.io/collector/pdata v1.0.0-rc7
	go.uber.org/zap v1.24.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/DataDog/datadog-agent/pkg/obfuscate v0.43.0-rc.3.0.20230206114529-17c7dfde736c // indirect
	github.com/DataDog/datadog-agent/pkg/remoteconfig/state v0.43.0-rc.3.0.20230206114529-17c7dfde736c // indirect
	github.com/DataDog/datadog-agent/pkg/util/cgroups v0.43.0-rc.3.0.20230206114529-17c7dfde736c // indirect
	github.com/DataDog/datadog-agent/pkg/util/log v0.43.0-rc.3.0.20230206114529-17c7dfde736c // indirect
	github.com/DataDog/datadog-agent/pkg/util/pointer v0.43.0-rc.3 // indirect
	github.com/DataDog/datadog-agent/pkg/util/scrubber v0.43.0-rc.3.0.20230206114529-17c7dfde736c // indirect
	github.com/DataDog/datadog-go/v5 v5.1.1 // indirect
	github.com/DataDog/go-tuf v0.3.0--fix-localmeta-fork // indirect
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes v0.1.2 // indirect
	github.com/DataDog/opentelemetry-mapping-go/pkg/quantile v0.1.2 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/containerd/cgroups v1.0.4 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/godbus/dbus/v5 v5.0.6 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/karrick/godirwalk v1.17.0 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20220517141722-cf486979b281 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417 // indirect
	github.com/outcaste-io/ristretto v0.2.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.4.0 // indirect
	github.com/shirou/gopsutil/v3 v3.23.2 // indirect
	github.com/tinylib/msgp v1.1.8 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.opentelemetry.io/collector/confmap v0.73.0 // indirect
	go.opentelemetry.io/collector/featuregate v0.73.0 // indirect
	go.opentelemetry.io/collector/semconv v0.73.0 // indirect
	go.opentelemetry.io/otel v1.14.0 // indirect
	go.opentelemetry.io/otel/metric v0.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.14.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/grpc v1.53.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// It appears that the v0.2.0 tag was modified.  Replacing with v0.2.1
replace github.com/outcaste-io/ristretto v0.2.0 => github.com/outcaste-io/ristretto v0.2.1
