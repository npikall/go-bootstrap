.DEFAULT_GOAL := help

# The variables in .env are exported to every recipe.
-include .env
export $(shell sed -n 's/^\([A-Za-z_][A-Za-z0-9_]*\)=.*/\1/p' .env 2>/dev/null)

# Missing development tools are installed system-wide with `go install`.
GOBIN_DIR := $(or $(shell go env GOBIN),$(shell go env GOPATH)/bin)
# A tool on PATH wins, otherwise it is used from where `go install` puts it.
tool = $(or $(shell command -v $(1) 2>/dev/null),$(GOBIN_DIR)/$(1))
AIR := $(call tool,air)
GOLANGCI_LINT := $(call tool,golangci-lint)
TEMPL := $(call tool,templ)
TAILWINDCSS := $(call tool,tailwindcss)
SQLC := $(call tool,sqlc)
GOOSE := $(call tool,goose)
GOTESTSUM := $(call tool,gotestsum)

AIR_VERSION := v1.67.4
GOLANGCI_LINT_VERSION := v2.15.0
SQLC_VERSION := v1.31.1
GOOSE_VERSION := v3.30.1
GOTESTSUM_VERSION := v1.13.0

.PHONY: all
all: build test ## Build and test the application

.PHONY: build
build: ## Build the application
	go build -o main ./cmd/api

.PHONY: run
run: ## Run the application
	go run ./cmd/api

.PHONY: test
test: install-gotestsum ## Test the application
	$(GOTESTSUM) -- ./...

.PHONY: test-watch
test-watch: install-gotestsum ## Rerun the tests of changed packages
	$(GOTESTSUM) --watch

.PHONY: sqlc
sqlc: install-sqlc ## Generate Go code from the SQL queries
	$(SQLC) generate

.PHONY: migrate-up
migrate-up: install-goose ## Apply all pending migrations
	$(GOOSE) up

.PHONY: migrate-down
migrate-down: install-goose ## Roll back the latest migration
	$(GOOSE) down

.PHONY: migrate-status
migrate-status: install-goose ## Show the migration status
	$(GOOSE) status

.PHONY: migrate-new
migrate-new: install-goose ## Create a migration, e.g. make migrate-new name=add_email
	$(GOOSE) -s create $(name) sql

.PHONY: lint
lint: install-golangci-lint ## Lint the application
	$(GOLANGCI_LINT) run

.PHONY: clean
clean: ## Clean the binary
	rm -f main
	rm -f demo.db

.PHONY: watch
watch: install-air ## Live Reload
	$(AIR)

.PHONY: help
help: ## Show this help
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@awk 'BEGIN { FS = ":.*## " } /^[a-zA-Z0-9_-]+:.*## / { printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: install-air
install-air:
	@command -v $(AIR) >/dev/null || { \
		echo "warning: air not found, installing github.com/air-verse/air@$(AIR_VERSION) system-wide with go install" >&2 && \
		go install github.com/air-verse/air@$(AIR_VERSION); \
	}

.PHONY: install-golangci-lint
install-golangci-lint:
	@command -v $(GOLANGCI_LINT) >/dev/null || { \
		echo "warning: golangci-lint not found, installing github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) system-wide with go install" >&2 && \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
	}

.PHONY: install-sqlc
install-sqlc:
	@command -v $(SQLC) >/dev/null || { \
		echo "warning: sqlc not found, installing github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION) system-wide with go install" >&2 && \
		go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION); \
	}

.PHONY: install-goose
install-goose:
	@command -v $(GOOSE) >/dev/null || { \
		echo "warning: goose not found, installing github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION) system-wide with go install" >&2 && \
		go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION); \
	}

.PHONY: install-gotestsum
install-gotestsum:
	@command -v $(GOTESTSUM) >/dev/null || { \
		echo "warning: gotestsum not found, installing gotest.tools/gotestsum@$(GOTESTSUM_VERSION) system-wide with go install" >&2 && \
		go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION); \
	}
