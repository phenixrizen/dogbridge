# dogbridge Implementation Plan

This plan prioritizes delivering a real, usable migration bridge quickly, then expanding signal coverage and production hardening.

## Phase 0 (Immediate cleanup)
- Replace placeholder docs with factual current-state docs.
- Make config examples internally consistent and backend-specific.
- Ensure repository metadata (`README`, `LICENSE`, module path) reflects intended ownership and release strategy.

## Phase 1 (MVP 1: Datadog traces -> Tempo)
**Goal:** Run a local end-to-end flow from `dd-trace-go` to Tempo via dogbridge.

### Tasks
1. Build a runnable collector distribution:
   - Wire `otelcol` main with component factories.
   - Include `datadogreceiver`, `otlpreceiver`, `batch`, `memory_limiter`, `resource`, `otlp` exporter.
2. Add default config for traces pipeline on port `8126`.
3. Complete `examples/go-ddtrace` with real tracer instrumentation and Dockerfile.
4. Complete `examples/docker-compose` with dogbridge, tempo, and grafana wiring.
5. Add a smoke test script that emits one trace and verifies it reaches Tempo APIs.

### Exit criteria
- `docker compose up` works locally for traces path.
- A sample trace appears in Tempo/Grafana.
- Docs include exact run commands and expected outputs.

## Phase 2 (MVP 2: DogStatsD metrics -> VictoriaMetrics)
**Goal:** Ingest DogStatsD metrics and export through Prometheus remote write.

### Tasks
1. Enable `statsdreceiver` and `prometheusremotewrite` exporter in distro.
2. Add transform processor rules for default high-cardinality tag drops.
3. Add configurable allow-list/deny-list label controls.
4. Extend compose stack with VictoriaMetrics and Grafana dashboards.
5. Add e2e test that sends counters/gauges/histogram-like traffic and validates writes.

### Exit criteria
- Metrics pipeline receives DogStatsD and writes to VictoriaMetrics.
- Cardinality defaults documented and test-covered.

## Phase 3 (MVP 3: Kubernetes logs -> Loki/OpenSearch)
**Goal:** Provide practical Kubernetes log collection with metadata enrichment.

### Tasks
1. Enable `filelogreceiver`, `k8sattributesprocessor`, and log exporters (Loki + Elasticsearch/OpenSearch).
2. Provide DaemonSet profile for node log collection.
3. Add JSON parsing and trace/span ID extraction operators.
4. Add Kubernetes deployment docs and minimal RBAC manifests.
5. Add integration tests (kind/minikube-friendly) for log ingestion path.

### Exit criteria
- Logs from pod stdout/stderr reach Loki (and optionally OpenSearch).
- Kubernetes metadata fields are present and documented.

## Phase 4 (Helm productionization)
**Goal:** Turn Helm scaffold into production-usable chart.

### Tasks
1. Implement real templates for deployment/daemonset/service/configmap/rbac.
2. Support `mode: deployment|daemonset|gateway`.
3. Add configurable ports, resources, autoscaling, node selectors, tolerations.
4. Add health probes and ServiceMonitor support.
5. Add chart tests and lint checks in CI.

### Exit criteria
- `helm lint` passes.
- Chart can deploy in both deployment and daemonset modes.

## Phase 5 (Hardening and operations)
**Goal:** Make dogbridge robust under production load.

### Tasks
1. Tune memory limiter, batch, queues, retry/backoff, and timeouts.
2. Add self-observability dashboards and dropped-telemetry counters.
3. Add load test scenarios for traces, metrics, and logs independently.
4. Document backpressure behavior and scaling guidance by signal type.
5. Add security guidance (mTLS, secrets, NetworkPolicy examples).

### Exit criteria
- Load tests show stable behavior under sustained traffic.
- Operational runbook exists for SRE handoff.

## Phase 6 (Migration enablement)
**Goal:** Help teams move from Datadog clients to native OTel with low risk.

### Tasks
1. Publish Datadog compatibility matrix and known limitations.
2. Add Go migration guide from `dd-trace-go` to OTel SDK/OTLP.
3. Provide dashboard and alert migration examples.
4. Add phased rollout playbook (single service -> dual write -> cutover).

### Exit criteria
- Migration docs are actionable and reviewed against real pilot services.

## Suggested execution order for next PRs
1. Runnable collector binary + traces MVP compose demo.
2. Real docs for architecture/datadog-compatibility/limitations.
3. Helm templates (minimum viable deployment mode).
4. Metrics and logs e2e expansions.
