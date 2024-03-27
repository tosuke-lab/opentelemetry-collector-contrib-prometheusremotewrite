# Websocket Processor

The WebSocket processor, which can be positioned anywhere in a pipeline, allows
data to pass through to the next component. Simultaneously, it makes a portion
of the data accessible to WebSocket clients connecting on a configurable port.
This functionality resembles that of the Unix `tee` command, which enables data
to flow through while duplicating and redirecting it for inspection.

To avoid overloading clients, the amount of telemetry duplicated over 
any open WebSockets is rate limited by an adjustable amount.

## Config

The WebSocket processor has two configurable fields: `port` and `limit`:

- `port`: The port on which the WebSocket processor listens. Optional. Defaults
  to `12001`.
- `limit`: The rate limit over the WebSocket in messages per second. Can be a
  float or an integer. Optional. Defaults to `1`.

Example configuration:

```yaml
websocket:
  port: 12001
  limit: 1 # rate limit 1 msg/sec
```
