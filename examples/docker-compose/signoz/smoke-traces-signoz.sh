#!/usr/bin/env bash
set -euo pipefail

echo "Running dd-trace smoke emitter against dogbridge..."
go run ./examples/go-ddtrace

echo "Trace emitted. Verify in SigNoz UI at http://localhost:3301"
