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

Note: that http2 need to be enabled.

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

```shell
# Make sure to fetch latest
go get -u github.com/opendroid/gcp_go_funcs/grpc_tests/notes
go get google.golang.org/grpc
# If seeing issues with protobuf
go get -u
go mod tidy
```

### First time deployment
```shell
# Be in server directory
cd server 
# Deploying CloudRun
export GOOGLE_CLOUD_PROJECT=gcp-experiments-334602
# New: New migrate to this
gcloud builds submit --tag us-west2-docker.pkg.dev/$GOOGLE_CLOUD_PROJECT/grpc-notes/notes:v8
# New: Allow UnAuth, use http2
gcloud run deploy notes --image us-west2-docker.pkg.dev/$GOOGLE_CLOUD_PROJECT/grpc-notes/notes:v8 --allow-unauthenticated --use-http2
# Get the host name of the notes service. --format='value(status.url)'
gcloud run services describe notes 

```

### Subsequent deployment

```shell
export GOOGLE_CLOUD_PROJECT=gcp-experiments-334602
gcloud run deploy notes --image us-west2-docker.pkg.dev/$GOOGLE_CLOUD_PROJECT/grpc-notes/notes:v6
```

### Local docker build

```shell
# Build image
docker build -t notes:v2 .

# Run gRPC server, on local docker
docker run -d -p 8080:8080 notes:v2 ./grpc_test_server
# Or, run gRPC server on local machine
go run main.go notes.go
 
# Run client, set NOTES_GRPC_ADDRESS server address to local
NOTES_GRPC_ADDRESS="localhost:8080" go run main.go
# Or, Test client with Cloud run GRPC notes server
NOTES_GRPC_ADDRESS="notes-2dbml6flea-wl.a.run.app:443" go run main.go
```

### Update service API documentation

```shell
# Deploy no auth version
gcloud endpoints services deploy ./notes/out.pb ./notes/api_config.yaml --verbosity=debug

# Deploy auth version
gcloud endpoints services deploy out.pb api_config_auth.yaml
```


