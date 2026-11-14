.DEFAULT_GOAL := help

export GOBIN := $(CURDIR)/bin

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
	$(GOBIN)/golangci-lint run

.PHONY: clean
clean: ## Clean the binary
	rm -f main

.PHONY: watch
watch: install-air ## Live Reload
	$(GOBIN)/air

.PHONY: help
help: ## Show this help
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@awk 'BEGIN { FS = ":.*## " } /^[a-zA-Z0-9_-]+:.*## / { printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: install-air
install-air:
	@test -x $(GOBIN)/air || go install github.com/air-verse/air@$(AIR_VERSION)

.PHONY: install-golangci-lint
install-golangci-lint:
	@test -x $(GOBIN)/golangci-lint || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
