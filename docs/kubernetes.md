# Kubernetes logs (MVP)

This guide configures dogbridge as a DaemonSet to collect pod stdout/stderr logs from each node and forward them to Loki or OpenSearch-compatible endpoints.

## What is implemented

- `filelog` receiver for `/var/log/pods/*/*/*.log` container log files.
- `k8sattributes` processor for namespace, pod, container, and node metadata enrichment.
- JSON parsing plus trace/span extraction from structured log fields (`trace_id`, `span_id`).
- Export path examples for Loki and OpenSearch (`elasticsearch` exporter).

## DaemonSet profile notes

- Run dogbridge in DaemonSet mode so every node can read local pod log files.
- Mount host paths:
  - `/var/log/pods`
  - `/var/log/containers` (symlink resolution in some distros)
- Use `serviceAccount` auth for Kubernetes metadata lookup.

## Minimal RBAC manifest

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: dogbridge
  namespace: observability
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: dogbridge-log-reader
rules:
  - apiGroups: [""]
    resources: ["pods", "namespaces", "nodes"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: dogbridge-log-reader
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: dogbridge-log-reader
subjects:
  - kind: ServiceAccount
    name: dogbridge
    namespace: observability
```

## Collector config examples

- Loki: `config/examples/logs-to-loki.yaml`
- OpenSearch: `config/examples/logs-to-opensearch.yaml`

## Kind/minikube integration testing

The repository includes config-level integration tests that validate required pipeline elements for both Loki and OpenSearch examples:

```bash
go test ./tests/e2e -run Logs
```

For a cluster-level test:
1. Deploy Loki in-cluster or expose a reachable endpoint.
2. Deploy dogbridge with the Loki config example.
3. Emit a pod log line with JSON containing `trace_id` and `span_id`.
4. Verify the log in Loki and check for `k8s.*` resource attributes.

## Remaining placeholders

- No bundled kind/minikube orchestration script yet.
- Helm chart now supports deployment and daemonset modes with RBAC, service, and collector config map wiring.
