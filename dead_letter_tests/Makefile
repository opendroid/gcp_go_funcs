# Makefile for dead_letter_tests (Pub/Sub Dead-Letter & Retries)

GCP_PROJECT       ?= gcp-experiments-334602
GCP_REGION        ?= us-central1
RUNTIME           ?= go126
TOPIC_MAIN        ?= radio-pluto
TOPIC_DEAD_LETTER ?= pluto-dead-letter
SUB_DEAD_LETTER   ?= dead-letter-reader
CALLING_SA        ?= 576696822186-compute@developer.gserviceaccount.com

.PHONY: help fmt vet lint test build check update tidy deploy-ack deploy-bad-ack deploy-all config-push-sub publish-main publish-dead-letter logs-ack logs-bad-ack

help: ## Show this help message
	@echo "Usage: make [target] [GCP_PROJECT=... GCP_REGION=...]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2}'

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

deploy-ack: check ## Pre-check and deploy AckPubMessage CloudEvent function
	@echo "==> Deploying AckPubMessage (Cloud Function gen2)..."
	gcloud functions deploy AckPubMessage \
		--gen2 \
		--runtime=$(RUNTIME) \
		--region=$(GCP_REGION) \
		--project=$(GCP_PROJECT)

deploy-bad-ack: check ## Pre-check and deploy BadAckFunc HTTP function (returns 500)
	@echo "==> Deploying BadAckFunc (Cloud Function gen2)..."
	gcloud functions deploy BadAckFunc \
		--gen2 \
		--runtime=$(RUNTIME) \
		--trigger-http \
		--region=$(GCP_REGION) \
		--project=$(GCP_PROJECT)

deploy-all: deploy-ack deploy-bad-ack ## Deploy both AckPubMessage and BadAckFunc

config-push-sub: ## Configure dead-letter push subscription with topic path and OIDC auth
	@echo "==> Configuring push subscription $(SUB_DEAD_LETTER)..."
	gcloud pubsub subscriptions modify-push-config $(SUB_DEAD_LETTER) \
		--project=$(GCP_PROJECT) \
		--push-endpoint="https://$(GCP_REGION)-$(GCP_PROJECT).cloudfunctions.net/AckPubMessage/projects/$(GCP_PROJECT)/topics/$(TOPIC_DEAD_LETTER)" \
		--push-auth-service-account="$(CALLING_SA)" \
		--push-auth-token-audience="https://$(GCP_REGION)-$(GCP_PROJECT).cloudfunctions.net/AckPubMessage"

publish-main: ## Publish a test message to radio-pluto to trigger retries and dead-letter
	@echo "==> Publishing message to $(TOPIC_MAIN)..."
	gcloud pubsub topics publish $(TOPIC_MAIN) \
		--project=$(GCP_PROJECT) \
		--message='{"name": "GCP", "rating": "21-star"}'

publish-dead-letter: ## Publish directly to pluto-dead-letter topic to test AckPubMessage
	@echo "==> Publishing message directly to $(TOPIC_DEAD_LETTER)..."
	gcloud pubsub topics publish $(TOPIC_DEAD_LETTER) \
		--project=$(GCP_PROJECT) \
		--message='{"test": "direct-dead-letter"}'

logs-ack: ## Read recent Cloud Run logs for ackpubmessage service
	@gcloud logging read 'resource.type="cloud_run_revision" AND resource.labels.service_name="ackpubmessage"' \
		--project=$(GCP_PROJECT) \
		--limit=20 \
		--freshness=10m \
		--format='table(timestamp,httpRequest.status,textPayload)'

logs-bad-ack: ## Read recent Cloud Run logs for badackfunc service
	@gcloud logging read 'resource.type="cloud_run_revision" AND resource.labels.service_name="badackfunc"' \
		--project=$(GCP_PROJECT) \
		--limit=20 \
		--freshness=10m \
		--format='table(timestamp,httpRequest.status,textPayload)'
