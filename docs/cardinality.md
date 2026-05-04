# Metrics cardinality controls

DogStatsD tags can create unbounded label cardinality in downstream Prometheus-compatible backends. `dogbridge` applies two layers of controls in the metrics pipeline used for VictoriaMetrics.

## Default high-cardinality drops

The `transform/drop_default_high_cardinality` processor removes common per-request identifiers:

- `claim_id`
- `raw_claim_id`
- `trace_id`
- `span_id`
- `request_id`
- `partition_key`
- `offset`
- `kafka_sequence`
- `runtime-id`
- `key_hash`

## Configurable allow-list / deny-list

The `transform/label_controls` processor supports environment-variable controls:

- `DOGBRIDGE_METRICS_LABEL_ALLOW_REGEX` (default: `.*`): keep only matching label keys.
- `DOGBRIDGE_METRICS_LABEL_DENY_REGEX` (default: `^$`): drop matching label keys after allow-list filtering.

These controls are applied in order: **allow-list first, deny-list second**.

## Example

Keep only `env`, `service`, `version`, and `team`, then drop `team`:

```bash
export DOGBRIDGE_METRICS_LABEL_ALLOW_REGEX='^(env|service|version|team)$'
export DOGBRIDGE_METRICS_LABEL_DENY_REGEX='^team$'
```

With this configuration, only `env`, `service`, and `version` labels survive.
