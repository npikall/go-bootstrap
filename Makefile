.DEFAULT_GOAL := help

export GOBIN := $(CURDIR)/bin

GOLANGCI_LINT_VERSION := v2.15.0

.PHONY: all
all: lint test build ## Lint, test and build

.PHONY: build
build: ## Build the binary
	go build -o go-bootstrap .

.PHONY: test
test: ## Run the unit tests
	go test ./...

.PHONY: e2e
e2e: install-golangci-lint ## Generate real projects and build, test and lint them (needs network)
	PATH="$(GOBIN):$$PATH" GO_BOOTSTRAP_E2E=1 go test ./internal/scaffold -run TestEndToEnd -timeout 30m -v

.PHONY: golden
golden: ## Update the golden files after changing the template
	go test ./internal/scaffold -run TestGolden -update

.PHONY: lint
lint: install-golangci-lint ## Lint the code
	$(GOBIN)/golangci-lint run

.PHONY: clean
clean: ## Clean the binary
	rm -f go-bootstrap

.PHONY: help
help: ## Show this help
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@awk 'BEGIN { FS = ":.*## " } /^[a-zA-Z0-9_-]+:.*## / { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: install-golangci-lint
install-golangci-lint:
	@test -x $(GOBIN)/golangci-lint || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
