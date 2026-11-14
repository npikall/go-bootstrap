.DEFAULT_GOAL := help

# Missing development tools are installed system-wide with `go install`.
GOBIN_DIR := $(or $(shell go env GOBIN),$(shell go env GOPATH)/bin)
# A tool on PATH wins, otherwise it is used from where `go install` puts it.
tool = $(or $(shell command -v $(1) 2>/dev/null),$(GOBIN_DIR)/$(1))
AIR := $(call tool,air)
GOLANGCI_LINT := $(call tool,golangci-lint)
TEMPL := $(call tool,templ)
TAILWINDCSS := $(call tool,tailwindcss)

AIR_VERSION := v1.67.4
GOLANGCI_LINT_VERSION := v2.15.0

.PHONY: all
all: build test ## Build and test the application

.PHONY: build
build: ## Build the application
	go build -o main ./cmd/api

.PHONY: run
run: ## Run the application
	go run ./cmd/api

.PHONY: test
test: ## Test the application
	go test ./...

.PHONY: lint
lint: install-golangci-lint ## Lint the application
	$(GOLANGCI_LINT) run

.PHONY: clean
clean: ## Clean the binary
	rm -f main

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
