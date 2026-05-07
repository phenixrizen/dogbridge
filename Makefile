SHELL := /bin/bash

COMPOSE_FILE := examples/docker-compose/docker-compose.yaml
SIGNOZ_COMPOSE_FILE := examples/docker-compose/signoz/docker-compose.yaml
OPENOBSERVE_COMPOSE_FILE := examples/docker-compose/openobserve/docker-compose.yaml

.PHONY: help build test fmt demo-up demo-down demo-smoke-traces demo-smoke-metrics demo-signoz-up demo-signoz-down demo-signoz-smoke demo-openobserve-up demo-openobserve-down demo-openobserve-smoke

help:
	@echo "dogbridge targets"
	@echo "  build               Build dogbridge binary"
	@echo "  test                Run go tests"
	@echo "  fmt                 Run go fmt"
	@echo "  demo-up             Start Tempo demo stack"
	@echo "  demo-down           Stop Tempo demo stack"
	@echo "  demo-smoke-traces   Run dd-trace smoke test against Tempo demo"
	@echo "  demo-smoke-metrics  Run statsd smoke test against Tempo demo"
	@echo "  demo-signoz-up      Start SigNoz demo stack"
	@echo "  demo-signoz-down    Stop SigNoz demo stack"
	@echo "  demo-signoz-smoke   Run dd-trace smoke test against SigNoz demo"
	@echo "  demo-openobserve-up    Start OpenObserve demo stack"
	@echo "  demo-openobserve-down  Stop OpenObserve demo stack"
	@echo "  demo-openobserve-smoke Run dd-trace smoke test against OpenObserve demo"

build:
	go build ./cmd/dogbridge

test:
	go test ./...

fmt:
	go fmt ./...

demo-up:
	docker compose -f $(COMPOSE_FILE) up -d
	@echo "Grafana: http://localhost:3000 (admin/admin)"


demo-down:
	docker compose -f $(COMPOSE_FILE) down -v

demo-smoke-traces:
	bash examples/docker-compose/smoke-traces.sh

demo-smoke-metrics:
	bash examples/docker-compose/smoke-metrics.sh

demo-signoz-up:
	docker compose -f $(SIGNOZ_COMPOSE_FILE) up -d --remove-orphans
	@echo "SigNoz: http://localhost:3301"

demo-signoz-down:
	docker compose -f $(SIGNOZ_COMPOSE_FILE) down -v

demo-signoz-smoke:
	DOGBRIDGE_ENDPOINT=http://localhost:8126 bash examples/docker-compose/signoz/smoke-traces-signoz.sh

demo-openobserve-up:
	docker compose -f $(OPENOBSERVE_COMPOSE_FILE) up -d --remove-orphans
	@echo "OpenObserve: http://localhost:5080 (root@example.com / Complexpass#123)"

demo-openobserve-down:
	docker compose -f $(OPENOBSERVE_COMPOSE_FILE) down -v

demo-openobserve-smoke:
	bash examples/docker-compose/openobserve/smoke-openobserve.sh
