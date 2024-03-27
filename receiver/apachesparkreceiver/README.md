# Apache Spark Receiver

| Status                   |                  |
| ------------------------ | ---------------- |
| Stability                | [development] |
| Supported pipeline types | metrics          |
| Distributions            | [contrib]        |

This receiver fetches metrics for an Apache Spark cluster through the Apache Spark REST API - specifically, the /metrics/json, /api/v1/applications/[app-id]/stages, /api/v1/applications/[app-id]/executors, and /api/v1/applications/[app-id]/jobs endpoints.

## Purpose

The purpose of this component is to allow monitoring of Apache Spark clusters and the applications running on them through the collection of performance metrics like memory utilization, CPU utilization, shuffle operations, garbage collection time, I/O operations, and more.

## Prerequisites

This receiver supports Apache Spark versions:

- 3.3.2+

## Configuration

These configuration options are for connecting to an Apache Spark application.

The following settings are optional:

- `collection_interval`: (default = `60s`): This receiver collects metrics on an interval. This value must be a string readable by Golang's [time.ParseDuration](https://pkg.go.dev/time#ParseDuration). Valid time units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`.
- `endpoint`: (default = `http://localhost:4040`): Apache Spark endpoint to connect to in the form of `[http][://]{host}[:{port}]`
- `whitelisted_application_names`: An array of Spark application names for which metrics should be collected. If no application names are specified, metrics will be collected for all Spark applications running on the cluster at the specified endpoint.

### Example Configuration

```yaml
receivers:
  apachespark:
    collection_interval: 60s
    endpoint: http://localhost:4040
    whitelisted_application_names:
    - PythonStatusAPIDemo
    - PythonLR
```

The full list of settings exposed for this receiver are documented [here](./config.go) with detailed sample configurations [here](./testdata/config.yaml).

## Metrics

Details about the metrics produced by this receiver can be found in [metadata.yaml](./metadata.yaml)

[development]: https://github.com/open-telemetry/opentelemetry-collector#development
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
