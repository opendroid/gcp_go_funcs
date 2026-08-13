# gcp_go_funcs

[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)](https://go.dev/)
[![GCP Cloud Run](https://img.shields.io/badge/Google%20Cloud-Cloud%20Run-4285F4?style=flat&logo=google-cloud)](https://cloud.google.com/run)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A collection of Google Cloud Platform (GCP) Go services, Cloud Run functions (Cloud Functions 2nd gen), Pub/Sub dead-letter topic workflows, Dialogflow CX webhooks, and gRPC implementations.

---

## 📁 Repository Structure

The repository is organized as a Go multi-module workspace (`go.work`):

| Directory | Description | Technologies |
| :--- | :--- | :--- |
| [`df-v2/`](df-v2/) | Dialogflow CX fulfillment webhook and HTTP Cloud Run functions (`HelloShirts`, `Echo`, `LogTestV2`) using `log/slog` structured JSON logging. | Cloud Run functions, `log/slog`, Go 1.26 |
| [`dead_letter_tests/`](dead_letter_tests/) | Pub/Sub retry mechanisms, dead-letter topics (`pluto-dead-letter`), and OIDC push authentication configurations for Cloud Run functions. | Pub/Sub, CloudEvent, Cloud Run, IAM |
| [`grpc_tests/`](grpc_tests/) | Complete gRPC service (`NotesService`) showcasing protobuf definitions, client, and Dockerized server deployed to Cloud Run with HTTP/2. | gRPC, Protobuf, Docker, Cloud Run HTTP/2 |

---

## 🛠️ Prerequisites

- **Go**: Version `1.26` or higher
- **Google Cloud SDK (`gcloud`)**: Authenticated and configured with your project
- **Docker**: For local testing and container image builds (for `grpc_tests`)
- **Protocol Buffers (`protoc`)**: For compiling `.proto` files (for `grpc_tests`)

Ensure the required Google Cloud APIs are enabled:

```shell
gcloud services enable \
  run.googleapis.com \
  cloudfunctions.googleapis.com \
  pubsub.googleapis.com \
  artifactregistry.googleapis.com \
  cloudbuild.googleapis.com \
  logging.googleapis.com
```

Create the Docker Artifact Registry repository (one-time setup for `grpc_tests`):

```shell
gcloud artifacts repositories create notes-grpc-server \
  --repository-format=docker \
  --location=us-west2 \
  --description="Docker repository for Notes gRPC service" \
  --project="$GCP_PROJECT"
```

---

## 🚀 Local Quality Checks & Makefile Automation

Each module includes a `Makefile` that runs formatters, linters (`golangci-lint`, `go vet`), unit tests, and compilation checks before deployment.

### Workspace-wide Checks

```shell
# Run full quality pipeline (fmt, vet, lint, test, build) across all modules
make check

# Update all Go dependencies to latest versions and tidy across all modules
make update

# Or individual checks and maintenance
make fmt      # Format all Go files with gofmt
make vet      # Run go vet across all modules
make lint     # Run golangci-lint across all modules
make test     # Run unit tests across all modules
make build    # Compile all modules
make tidy     # Run go mod tidy across all modules
```

### Module-specific Makefiles & Deployments

Each sub-directory contains its own `Makefile` with dedicated pre-flight checks and `gcloud` deployment shortcuts:

```shell
# Deploy df-v2 Cloud Run functions
make -C df-v2 deploy-all

# Deploy dead_letter_tests Cloud Functions (gen2) & configure push subscription
make -C dead_letter_tests deploy-all
make -C dead_letter_tests config-push-sub

# Build container and deploy gRPC service to Cloud Run with HTTP/2
make -C grpc_tests deploy-all
```

---

## 📦 Modules Overview & Deployment

### 1. Dialogflow CX & Cloud Run Functions (`df-v2`)

Implements Dialogflow CX webhook handling (`HandleWebhookRequest`) and HTTP utility functions.

- **Structured Logging**: Uses standard library `log/slog` configured for JSON output to `os.Stdout`, which Cloud Run and Cloud Logging ingest automatically.
- **Deployment with `gcloud run deploy`**:

```shell
# Deploy Dialogflow CX Webhook
gcloud run deploy cx-webhook \
  --source=df-v2 \
  --function=HandleWebhookRequest \
  --base-image=go126 \
  --project="$GCP_PROJECT" \
  --region="us-central1" \
  --no-allow-unauthenticated

# Deploy HTTP Echo function
gcloud run deploy echo \
  --source=df-v2 \
  --function=Echo \
  --base-image=go126 \
  --project="$GCP_PROJECT" \
  --region="us-central1" \
  --no-allow-unauthenticated
```

See [df-v2/README.md](df-v2/README.md) for detailed deployment patterns, URL resolution, and testing scripts.

---

### 2. Pub/Sub Retries & Dead-Letter Topics (`dead_letter_tests`)

Demonstrates how to configure Pub/Sub push subscriptions with dead-letter forwarding, retry limits, and secure OIDC authentication.

```text
radio-pluto
    │
    ▼
[Source Subscription]
    │
    ├── Retries on HTTP 500 (BadAckFunc)
    └── After max delivery attempts (e.g. 5)
            │
            ▼
     pluto-dead-letter (Topic)
            │
            ▼
     dead-letter-reader (Push Subscription with OIDC)
            │
            ▼
     AckPubMessage (Cloud Run Function -> HTTP 200)
```

- **Push Authentication & Audience**: Proper configuration to ensure Cloud Run accepts OIDC ID tokens without 401/403 errors.
- **Topic URL Routing**: Setting push endpoint path patterns to `/projects/PROJECT_ID/topics/TOPIC_NAME` for proper Functions Framework topic detection.

See [dead_letter_tests/README.md](dead_letter_tests/README.md) for step-by-step IAM permissions, subscription commands, and test verification.

---

### 3. gRPC on Cloud Run (`grpc_tests`)

Demonstrates an end-to-end gRPC service (`NotesService`) running on Cloud Run:

- `grpc_tests/notes`: Protobuf definitions and generated Go stubs.
- `grpc_tests/server`: Multi-stage Dockerfile packaging the gRPC server.
- `grpc_tests/client`: Demo client connecting over TLS / HTTP/2.

```shell
# Deploy gRPC server with HTTP/2 enabled
make -C grpc_tests deploy-all

# Test gRPC client against live Cloud Run deployment
make -C grpc_tests client-cloud

# Test gRPC client locally against localhost:8080
make -C grpc_tests client-local
```

See [grpc_tests/README.md](grpc_tests/README.md) for proto compilation, local Docker testing, and Cloud Run deployment details.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
