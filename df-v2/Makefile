# Makefile for df-v2 (Dialogflow CX & Cloud Run functions)

GCP_PROJECT ?= gcp-experiments-334602
GCP_REGION  ?= us-central1
RUNTIME     ?= go126

.PHONY: help fmt vet lint test build check update tidy deploy-cx-webhook deploy-hello-shirts deploy-echo deploy-log-test deploy-all describe-cx describe-hello test-cx

help: ## Show this help message
	@echo "Usage: make [target] [GCP_PROJECT=... GCP_REGION=...]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

update: ## Update Go dependencies to latest versions and tidy
	@echo "==> Updating dependencies..."
	go get -u ./...
	go mod tidy

tidy: ## Run go mod tidy
	@echo "==> Tidying module..."
	go mod tidy

fmt: ## Format Go source files
	@echo "==> Formatting Go files..."
	gofmt -s -w .

vet: ## Run go vet
	@echo "==> Running go vet..."
	go vet ./...

lint: ## Run golangci-lint
	@echo "==> Running golangci-lint..."
	golangci-lint run ./...

test: ## Run unit tests
	@echo "==> Running unit tests..."
	go test -v -race ./...

build: ## Compile Go package
	@echo "==> Compiling package..."
	go build -buildvcs=false ./...

check: fmt vet lint test build ## Run all code quality checks (fmt, vet, lint, test, build)

deploy-cx-webhook: check ## Pre-check and deploy Dialogflow CX webhook (HandleWebhookRequest)
	@echo "==> Deploying cx-webhook to Cloud Run..."
	gcloud run deploy cx-webhook \
		--source=. \
		--function=HandleWebhookRequest \
		--base-image=$(RUNTIME) \
		--project=$(GCP_PROJECT) \
		--region=$(GCP_REGION) \
		--no-allow-unauthenticated

deploy-hello-shirts: check ## Pre-check and deploy HelloShirts HTTP function
	@echo "==> Deploying hello-shirts to Cloud Run..."
	gcloud run deploy hello-shirts \
		--source=. \
		--function=HelloShirts \
		--base-image=$(RUNTIME) \
		--project=$(GCP_PROJECT) \
		--region=$(GCP_REGION) \
		--no-allow-unauthenticated

deploy-echo: check ## Pre-check and deploy Echo HTTP function
	@echo "==> Deploying echo to Cloud Run..."
	gcloud run deploy echo \
		--source=. \
		--function=Echo \
		--base-image=$(RUNTIME) \
		--project=$(GCP_PROJECT) \
		--region=$(GCP_REGION) \
		--no-allow-unauthenticated

deploy-log-test: check ## Pre-check and deploy LogTestV2 HTTP function
	@echo "==> Deploying log-test-v2 to Cloud Run..."
	gcloud run deploy log-test-v2 \
		--source=. \
		--function=LogTestV2 \
		--base-image=$(RUNTIME) \
		--project=$(GCP_PROJECT) \
		--region=$(GCP_REGION) \
		--no-allow-unauthenticated

deploy-all: deploy-cx-webhook deploy-hello-shirts deploy-echo ## Deploy all df-v2 functions to Cloud Run

describe-cx: ## Get cx-webhook service status URL
	@gcloud run services describe cx-webhook --project=$(GCP_PROJECT) --region=$(GCP_REGION) --format='value(status.url)'

describe-hello: ## Get hello-shirts service status URL
	@gcloud run services describe hello-shirts --project=$(GCP_PROJECT) --region=$(GCP_REGION) --format='value(status.url)'

test-cx: ## Test cx-webhook with an authenticated sample payload
	@URL=$$(gcloud run services describe cx-webhook --project=$(GCP_PROJECT) --region=$(GCP_REGION) --format='value(status.url)'); \
	echo "Posting to $$URL..."; \
	curl -m 70 -X POST "$$URL" \
		-H "Authorization: Bearer $$(gcloud auth print-identity-token)" \
		-H "Content-Type: application/json" \
		-d '{"fulfillmentInfo": {"tag": "welcome"}}'
