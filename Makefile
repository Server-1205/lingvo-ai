# --- Makefile for Lingvo AI ---
# Usage: make [target]

SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

# --- Project ---
PROJECT    ?= lingvo
GO         ?= go
GOFLAGS    ?=
LDFLAGS    ?= -s -w
NPM        ?= pnpm

# --- Git ---
VERSION   ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT    ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# --- Build ---
BIN_DIR    := bin
MAIN_PKG   := ./cmd/server
BINARY     := $(BIN_DIR)/server

# --- Frontend ---
WEB_DIR    := web
WEB_DIST   := web/dist

# --- Docker ---
DOCKER_REGISTRY ?= ghcr.io
DOCKER_IMAGE    ?= $(DOCKER_REGISTRY)/lingvo-ai
DOCKER_TAG      ?= $(VERSION)

# --- Tools ---
GOLANGCI_LINT ?= golangci-lint
GOSEC         ?= gosec
GOTEST        ?= $(GO) test
GOTESTFLAGS   ?= -race -count=1
AIR           ?= air

# ============================================================================
.DEFAULT_GOAL := help

##@ Development

.PHONY: all
all: build build-frontend ## Build everything

.PHONY: install
install: ## Install dev tools (air, golangci-lint, gosec)
	@which $(AIR) 2>/dev/null || go install github.com/air-verse/air@latest
	@which $(GOLANGCI_LINT) 2>/dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@which $(GOSEC) 2>/dev/null || go install github.com/securego/gosec/v2/cmd/gosec@latest

.PHONY: setup
setup: tidy install frontend-install ## Full project setup (deps + tools)

.PHONY: dev
dev: ## Run backend with hot reload (requires air)
	$(AIR)

.PHONY: dev-frontend
dev-frontend: ## Run frontend dev server
	cd $(WEB_DIR) && $(NPM) dev

.PHONY: dev-all
dev-all: ## Run backend + frontend concurrently
	@echo "Starting backend (air) and frontend (vite) concurrently..."
	@trap 'kill 0' EXIT; ($(AIR)) & (cd $(WEB_DIR) && $(NPM) dev) & wait

.PHONY: build
build: ## Build backend binary
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BINARY) $(MAIN_PKG)

.PHONY: build-frontend
build-frontend: frontend-install ## Build frontend assets
	cd $(WEB_DIR) && $(NPM) run build

.PHONY: frontend-install
frontend-install: ## Install frontend dependencies
	cd $(WEB_DIR) && $(NPM) install

.PHONY: run
run: build ## Build and run backend
	$(BINARY)

.PHONY: generate
generate: ## Run go generate
	$(GO) generate ./...

##@ Testing

.PHONY: test
test: ## Run Go tests
	$(GOTEST) $(GOTESTFLAGS) ./...

.PHONY: test-frontend
test-frontend: ## Run frontend tests (vitest)
	cd $(WEB_DIR) && npx vitest run

.PHONY: test-cover
test-cover: ## Run tests with coverage report
	$(GOTEST) $(GOTESTFLAGS) -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: test-integration
test-integration: ## Run integration tests
	$(GOTEST) $(GOTESTFLAGS) -tags=integration ./...

.PHONY: bench
bench: ## Run benchmarks
	$(GO) test -bench=. -benchmem ./...

##@ Code Quality

.PHONY: lint
lint: ## Run linters
	$(GOLANGCI_LINT) run ./...

.PHONY: lint-frontend
lint-frontend: ## Run frontend linter
	cd $(WEB_DIR) && $(NPM) lint

.PHONY: fmt
fmt: ## Format Go code
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: tidy
tidy: ## Tidy and verify go.mod
	$(GO) mod tidy
	$(GO) mod verify

.PHONY: typecheck
typecheck: ## Run frontend type checking
	cd $(WEB_DIR) && npx tsc --noEmit

.PHONY: audit
audit: ## Run security audit (Go + frontend)
	$(GO) mod verify
	cd $(WEB_DIR) && $(NPM) audit --audit-level=high || true

.PHONY: security
security: ## Run gosec on Go code
	$(GOSEC) -quiet -fmt=text ./...

.PHONY: check
check: lint vet test test-frontend build build-frontend ## Full pre-commit gate

##@ Docker

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):latest \
		.

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run --rm -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: docker-push
docker-push: ## Push Docker image
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest

##@ CI

.PHONY: ci
ci: check ## Run full CI pipeline (alias for check)

##@ Database

.PHONY: reset-db
reset-db: ## Reset SQLite database (remove data file)
	rm -f data/*.db data/*.db-shm data/*.db-wal
	@echo "Database reset. Will be recreated on next server start."

##@ Cleanup

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BIN_DIR) coverage.out coverage.html
	rm -rf $(WEB_DIST)/assets

##@ Help

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2} \
		/^##@/ {printf "\n\033[1m%s\033[0m\n", substr($$0, 5)}' $(MAKEFILE_LIST)
