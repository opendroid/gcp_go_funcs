# Root Makefile for gcp_go_funcs multi-module workspace

GCP_PROJECT ?= gcp-experiments-334602
GCP_REGION  ?= us-central1

MODULES := df-v2 dead_letter_tests grpc_tests

.PHONY: help fmt vet lint test build check all clean \
	deploy-df deploy-dead-letter deploy-grpc

help: ## Show this help message
	@echo "Usage: make [target] [GCP_PROJECT=... GCP_REGION=...]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

all: check ## Run all checks across all modules

fmt: ## Format Go source files across all modules
	@for m in $(MODULES); do \
		echo "==> Formatting $$m..."; \
		$(MAKE) -C $$m fmt; \
	done

vet: ## Run go vet across all modules
	@for m in $(MODULES); do \
		echo "==> Running go vet in $$m..."; \
		$(MAKE) -C $$m vet; \
	done

lint: ## Run golangci-lint across all modules
	@for m in $(MODULES); do \
		echo "==> Running lint in $$m..."; \
		$(MAKE) -C $$m lint; \
	done

test: ## Run unit tests across all modules
	@for m in $(MODULES); do \
		echo "==> Running tests in $$m..."; \
		$(MAKE) -C $$m test; \
	done

build: ## Build all packages across modules
	@for m in $(MODULES); do \
		echo "==> Building in $$m..."; \
		$(MAKE) -C $$m build; \
	done

check: ## Run full code quality pipeline (fmt, vet, lint, test, build) across all modules
	@for m in $(MODULES); do \
		echo "=========================================="; \
		echo "Checking module: $$m"; \
		echo "=========================================="; \
		$(MAKE) -C $$m check || exit 1; \
	done
	@echo ""
	@echo "All workspace modules passed checks successfully!"

deploy-df: ## Pre-check and deploy df-v2 services
	$(MAKE) -C df-v2 deploy-all GCP_PROJECT=$(GCP_PROJECT) GCP_REGION=$(GCP_REGION)

deploy-dead-letter: ## Pre-check and deploy dead_letter_tests functions
	$(MAKE) -C dead_letter_tests deploy-all GCP_PROJECT=$(GCP_PROJECT) GCP_REGION=$(GCP_REGION)

deploy-grpc: ## Pre-check and deploy grpc_tests service
	$(MAKE) -C grpc_tests deploy-all GCP_PROJECT=$(GCP_PROJECT)
