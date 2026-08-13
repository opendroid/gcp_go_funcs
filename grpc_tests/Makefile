# Makefile for grpc_tests (Notes gRPC Service on Cloud Run)

GCP_PROJECT ?= gcp-experiments-334602
GCP_REGION  ?= us-central1
IMAGE_TAG   ?= us-west2-docker.pkg.dev/$(GCP_PROJECT)/notes-grpc-server/notes-server:v1
SERVICE     ?= notes

.PHONY: help fmt vet lint test build check update tidy proto docker-build build-submit deploy deploy-all client-local client-cloud

help: ## Show this help message
	@echo "Usage: make [target] [GCP_PROJECT=... GCP_REGION=...]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

update: ## Update Go dependencies to latest versions across grpc_tests modules
	@echo "==> Updating dependencies in notes..."
	cd notes && go get -u ./... && go mod tidy
	@echo "==> Updating dependencies in client..."
	cd client && go get -u ./... && go mod tidy
	@echo "==> Updating dependencies in server..."
	cd server && go get -u ./... && go mod tidy

tidy: ## Run go mod tidy across grpc_tests modules
	@echo "==> Tidying grpc_tests modules..."
	cd notes && go mod tidy
	cd client && go mod tidy
	cd server && go mod tidy

proto: ## Compile notes.proto with protoc to generate Go gRPC stubs
	@echo "==> Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--descriptor_set_out notes/out.pb \
		--include_imports --include_source_info \
		notes/notes.proto

fmt: ## Format Go source files across grpc_tests
	@echo "==> Formatting Go files..."
	gofmt -s -w client notes server

vet: ## Run go vet across grpc_tests modules
	@echo "==> Running go vet..."
	cd client && go vet ./...
	cd notes && go vet ./...
	cd server && go vet ./...

lint: ## Run golangci-lint across grpc_tests modules
	@echo "==> Running golangci-lint..."
	cd client && golangci-lint run ./...
	cd notes && golangci-lint run ./...
	cd server && golangci-lint run ./...

test: ## Run unit tests across grpc_tests modules
	@echo "==> Running unit tests..."
	cd client && go test -v -race ./...
	cd notes && go test -v -race ./...
	cd server && go test -v -race ./...

build: ## Build client, notes, and server Go packages (dry compile)
	@echo "==> Building Go packages..."
	cd client && go build -o /dev/null -buildvcs=false ./...
	cd notes && go build -buildvcs=false ./...
	cd server && go build -o /dev/null -buildvcs=false ./...

check: fmt vet lint test build ## Run all code quality checks (fmt, vet, lint, test, build)

docker-build: check ## Build local Docker container image for server
	@echo "==> Building Docker image for server..."
	cd server && docker build -t notes-server:v1 .

build-submit: check ## Build and push container image to Artifact Registry using Cloud Build
	@echo "==> Submitting Cloud Build for container image..."
	cd server && gcloud builds submit --tag $(IMAGE_TAG) --project=$(GCP_PROJECT)

deploy: ## Deploy containerized gRPC server to Cloud Run with HTTP/2 enabled
	@echo "==> Deploying gRPC service $(SERVICE) to Cloud Run..."
	gcloud run deploy $(SERVICE) \
		--image=$(IMAGE_TAG) \
		--use-http2 \
		--allow-unauthenticated \
		--region=$(GCP_REGION) \
		--project=$(GCP_PROJECT)

deploy-all: check build-submit deploy ## Full cycle: lint/test -> build container -> deploy to Cloud Run

client-local: ## Run demo gRPC client pointing to localhost:8080
	@echo "==> Running demo client against localhost:8080..."
	cd client && NOTES_GRPC_ADDRESS="localhost:8080" go run .

client-cloud: ## Run demo gRPC client pointing to Cloud Run deployment
	@URL=$$(gcloud run services describe $(SERVICE) --project=$(GCP_PROJECT) --region=$(GCP_REGION) --format='value(status.url)'); \
	echo "==> Running demo client against $$URL..."; \
	cd client && NOTES_GRPC_ADDRESS="$$URL" go run .
