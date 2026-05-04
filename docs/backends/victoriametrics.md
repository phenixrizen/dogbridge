# VictoriaMetrics backend

Use `config/examples/metrics-to-victoriametrics.yaml` to receive DogStatsD traffic on UDP `8125` and export metrics via Prometheus Remote Write.

## Pipeline summary

- Receiver: `statsd`
- Processors: `memory_limiter`, `transform/drop_default_high_cardinality`, `transform/label_controls`, `batch`
- Exporter: `prometheusremotewrite`

## Local compose target

The `examples/docker-compose/docker-compose.yaml` stack includes a single-node VictoriaMetrics service and wires dogbridge remote write traffic to:

- `http://victoriametrics:8428/api/v1/write`

For label control options, see `docs/cardinality.md`.
