# Datadog Processor

| Status                   |               |
|--------------------------|---------------|
| Stability                | [beta]        |
| Supported pipeline types | traces        |
| Distributions            | [contrib]     |

## Description

The Datadog Processor can be used to compute Datadog APM Stats pre-sampling. For example, when using the [tailsamplingprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/tailsamplingprocessor#tail-sampling-processor) or [probabilisticsamplerprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/probabilisticsamplerprocessor) components, the `datadogprocessor` can be prepended into the pipeline to ensure that Datadog APM Stats are accurate and include the dropped traces.

## Usage

To use the Datadog Processor, simply prepend it into a pipeline before any sampling processor. The Datadog Processor will compute APM Stats on all spans that it sees. Here is an example on how to add it to a pipeline using the [probabilisticsampler](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/probabilisticsamplerprocessor):

<table>
<tr>
<td> Before </td> <td> After </td>
</tr>
<tr>
<td valign="top">

```yaml
# ...
processors:
  # ...
  probabilistic_sampler:
    sampling_percentage: 20

exporters:
  datadog:
    api:
      key: ${env:DD_API_KEY}

service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [datadog]
    traces:
      receivers: [otlp]
      processors: [batch, probabilistic_sampler]
      exporters: [datadog]
```

</td><td valign="top">

```yaml
# ...
processors:
  # ...
  probabilistic_sampler:
    sampling_percentage: 20
  # add the "datadog" processor definition
  datadog:

exporters:
  datadog:
    api:
      key: ${env:DD_API_KEY}

service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [datadog]
    traces:
      receivers: [otlp]
      # prepend it to the sampler in your pipeline:
      processors: [batch, datadog, probabilistic_sampler]
      exporters: [datadog]
```

</tr></table>

Simply add the Datadog Processor into your list of processors and prepend it to the sampler in the traces pipeline to ensure it sees all spans.

## Configuration

By default, when used in conjunction with the Datadog Exporter, the processor should detect its presence (as long as it is configured within a pipeline), and use it to export the Datadog APM Stats. No configuration needed!

If using within a gateway deployment or running alongside the Datadog Agent where the Datadog Exporter is not present, then you must specify an alternative exporter to use, such as for example an OTLP exporter:

```yaml
processors:
  datadog:
    metrics_exporter: otlp
```

The default value for `metrics_exporter` is `datadog`. If your Datadog Exporter has a different name, you must specify it via config. Any configured metrics exporter must exist as part of a metrics pipeline.

When using in conjunction with the Datadog Agent's OTLP Ingest, the minimum required Datadog Agent version that supports this processor is 7.42.0.

If not using the Datadog backend, the processor will still create valid RED metrics, but in that situation you may prefer to use the [spanmetricsprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/spanmetricsprocessor) instead.

[beta]: https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]:https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
