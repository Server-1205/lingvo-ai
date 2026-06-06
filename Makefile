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
GOTEST        ?= $(GO) test
GOTESTFLAGS   ?= -race -count=1

# ============================================================================
.DEFAULT_GOAL := help

##@ Development

.PHONY: dev
dev: ## Run backend with hot reload (requires air)
	air

.PHONY: dev-frontend
dev-frontend: ## Run frontend dev server
	cd $(WEB_DIR) && $(NPM) dev

.PHONY: build
build: ## Build backend binary
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BINARY) $(MAIN_PKG)

.PHONY: build-frontend
build-frontend: ## Build frontend assets
	cd $(WEB_DIR) && $(NPM) install && $(NPM) run build

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
test-frontend: ## Run frontend tests
	cd $(WEB_DIR) && $(NPM) test run

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
ci: lint vet test build build-frontend ## Run full CI pipeline

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
