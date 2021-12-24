# tests gRPC API in a cloud run environment

 - [GRPC status codes](https://developers.google.com/maps-booking/reference/grpc-api/status_codes)
 - [Protocol Buffers](https://developers.google.com/protocol-buffers)
 - [Language Guide](https://developers.google.com/protocol-buffers/docs/proto3)
 - [GCP Managed Service Endpoints](https://github.com/GoogleCloudPlatform/golang-samples/tree/main/endpoints/getting-started-grpc)
 - [Cloud Run gRPC](https://cloud.google.com/run/docs/triggering/grpc)
 - [Regenerate gRPC code](https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code)
 - [New Go API for Protobuf](https://go.dev/blog/protobuf-apiv2)

## Installing protobuf

Use this to install the 
```shell
 # Install protoc, installs protoc
brew install protobuf
# Upgrade protobuf 
brew upgrade protobuf 
```

### Deploy Service Definition

Generate a PB file, that will be used with service deploy
```shell
# Generate the client files, note use 'source_relative' to generate .notes.go file in same directory
protoc --go_out=paths=source_relative:. --descriptor_set_out notes/out.pb --include_imports --include_source_info notes/notes.proto

 # generate the services required 'go-grpc_out'
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --descriptor_set_out notes/out.pb \
    --include_imports --include_source_info \
    notes/notes.proto

# Deploy endpoints from directory 'grpc_tests'
gcloud endpoints services deploy notes/out.pb api_config.yaml  --verbosity=debug
```

Note that on GCP we can host gRPC servers on [CloudRun, GKE or Compute Engine](https://cloud.google.com/endpoints/docs/grpc/about-grpc).
It is not available on CloudFunctions or GAE.

## gRPC Service Organization

It has three directories:
 - pb: Service definition in .proto files and generated definitions
 - server: Server side
 - client: Client side example

## Server Deployment

### First time deployment
```shell
# Be in server directory
cd serer 
# Deploying CloudRun
export GOOGLE_CLOUD_PROJECT=gcp-experiments-334602
# Build and push container images.
gcloud builds submit --tag gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-notes

# Deploy notes service for private access.
gcloud run deploy notes-upstream --image gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-notes

# Get the host name of the notes service.
NOTES_URL=$(gcloud run services describe notes-upstream --format='value(status.url)')
NOTES_DOMAIN=${NOTES_URL#https://}

# Deploy notes-relay service for public access.
gcloud run deploy notes --image gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-notes \
    --update-env-vars GRPC_NOTES_HOST=${NOTES_DOMAIN}:443 \
    --allow-unauthenticated
```

### Subsequent deployment

```shell
export GOOGLE_CLOUD_PROJECT=gcp-experiments-334602
gcloud run deploy notes --image gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-notes
gcloud run deploy notes-relay --image gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-notes
```

### Update service API documentation

```shell
# Deploy no auth version
gcloud endpoints services deploy ./notes/out.pb ./notes/api_config.yaml --verbosity=debug

# Deploy auth version
gcloud endpoints services deploy out.pb api_config_auth.yaml
```


