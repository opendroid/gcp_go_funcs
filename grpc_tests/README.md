# tests gRPC API in a cloud run environment

The `grpc_tests` demonstrates how to implement a [gRPC](https://grpc.io/docs/what-is-grpc/introduction/) service, and 
deploy it on [GCP Cloud Run](https://cloud.google.com/endpoints/docs/grpc/about-grpc).
We can also host gRPC servers on [CloudRun, GKE or Compute Engine](https://cloud.google.com/endpoints/docs/grpc/about-grpc).
It is not available on CloudFunctions or GAE.

The steps to implement a gRPC service are:
1. Define service contract in a `.proto` file.
2. Generate client and server proto definitions from `.proto` file
3. Write the server code, and
4. Write a sample demo client, to access the server (deployed on local docker or Cloud Run).

## Sample `Notes` Service

This example implements the `notes service`. The code is organized as:
1. `notes`: Contains `notes.proto` file and `protoc` generated code
2. `server`: Implements the server side contract. It is deployed on GCP Cloud Run, or a local docker
3. `client`: Is a sample client code that tests the server code.

The service allows a client to:
1. Create a Note
2. Fetch all notes created by an Author, and
3. Fetch a specific note.

### Defining `notes` Contract

The service is defined in `notes` directory. You may choose to create the client, server and .proto definitions as 
separate repos. It will make dependency management easier. The sample `notes` service contract is:
```protobuf
message CreateNoteRequest {} // rpc request format to create a note
message CreateNoteResponse {} // rpc response of creation of a note
message GetNoteRequest {} // rpc request, to get a specific note by ID
message GetNoteResponse {} // rpc response, with the note contents
message GetNotesByAuthorRequest {} // rpc request, to fetch all notes by a specific author
message GetNotesByAuthorResponse {}  // rpc response, containing all notes
service NotesService { // NotesService is the main service.
  rpc CreateNote(CreateNoteRequest) returns(CreateNoteResponse);
  rpc GetNote(GetNoteRequest) returns (GetNoteResponse);
  rpc GetNotesByAuthor(GetNotesByAuthorRequest) returns (GetNotesByAuthorResponse);
}
```
The contract is implemented in [notes.proto](https://github.com/opendroid/gcp_go_funcs/blob/main/grpc_tests/notes/notes.proto)
file. Take time and thought to create the contract. Once a contract is created, keeping server and
client side of code in sync becomes harder.

After writing the gRPC service `notes.proto` definition file, generate the message definition [notes.pb.go](https://github.com/opendroid/gcp_go_funcs/blob/main/grpc_tests/notes/notes.pb.go)
and server interface implementation [notes_grpc.pb.go](https://github.com/opendroid/gcp_go_funcs/blob/main/grpc_tests/notes/notes_grpc.pb.go), using command:
```shell
 # generate the services required 'go-grpc_out'
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --descriptor_set_out notes/out.pb \
    --include_imports --include_source_info \
    notes/notes.proto
```
The `notes/out.pb` file is used to deploy and/or document the endpoints.

#### GCP Endpoints
We can also document the Service definitions in [GCP Endpoints](https://cloud.google.com/endpoints). 
We have not configured this yet. The sample `endpoint` commands are:
```shell
# Deploy no auth version
gcloud endpoints services deploy ./notes/out.pb ./notes/api_config.yaml --verbosity=debug
# Deploy auth version
gcloud endpoints services deploy out.pb api_config_auth.yaml
```

#### Installing protobuf
We need to install `protoc` on your computer. On macOS use `brew`:
```shell
 # Install protoc, installs protoc
brew install protobuf
# Upgrade protobuf 
brew upgrade protobuf 
```
### Build Server Component
Goto the `server` directory, and build it (properly), fetch the right dependencies.
Note that the dependencies may not the code that is cloned.
```shell
# Make sure to fetch latest
go get -u github.com/opendroid/gcp_go_funcs/grpc_tests/notes
go get google.golang.org/grpc
# If dependencies issues try these
go get -u
go mod tidy
```
Once we can successfully build, deploy to __Local docker__ or __Cloud Run__ environments.

Note that if you are making changes to `notes` definition, during development you may set the:
```shell
# Temporarily use local notes directory.
go mod edit -replace github.com/opendroid/gcp_go_funcs/grpc_tests/notes=../notes
```
However, the relative backward path won't work for crating Docker image as it works in current directory.
To do that you will need to [clone git-repo into Docker](https://janakerman.co.uk/docker-git-clone/).

#### Deploying to Local docker
Before we deploy to GCP Cloud Run, test the server and client code locally. Once, it works on a local docker
proceed with Cloud Run deploy. Build local `server` docker using these commands:

```shell
# Run these in 'server' directory
docker build -t notes:v1 . # Build docker image as v1
```
Once deployed, test the client and server. The server will expose gRPC on port __localhost:8080__.
```shell
# Run Notes gRPC server, on local docker, expose :8080
docker run --rm -p 8080:8080 notes:v1 ./grpc_test_server
# Run client, (in 'client' dir), set NOTES_GRPC_ADDRESS server address to local
NOTES_GRPC_ADDRESS="localhost:8080" go run main.go
```

#### Deploying to GCP Cloud Run

| :exclamation: | For streaming gRPC, enable http2 |
|:----------:|:---------------------------------|

To deploy the server in __Cloud Run__ be in `server` directory. First make sure that the auth and GCP projects are set up appropriately:
```shell
gcloud config set account ajaythakur1972@gmail.com # Make sure you have right login
gcloud config configurations activate gcp-experiments # Activate right project
```
The Cloud Run exposes the gRPC on [port 443](https://ahmet.im/blog/grpc-auth-cloud-run/). 
Use these commands to deploy the Cloud Run version of server:
```shell
cd server # Be in server directory 
export GOOGLE_CLOUD_PROJECT=gcp-experiments-334602
# Build the image and keep it in Artifact Repository (not GCR)
gcloud builds submit --tag us-west2-docker.pkg.dev/$GOOGLE_CLOUD_PROJECT/grpc-notes/notes:v9
# Deploy: Allow UnAuthenticated, use http2
gcloud run deploy notes --image us-west2-docker.pkg.dev/$GOOGLE_CLOUD_PROJECT/grpc-notes/notes:v9 \
  --allow-unauthenticated --use-http2
gcloud run services describe notes  # Check the service configuration, if HTTP2 is enabled
# Client: Test with Cloud run GRPC notes server using command: (in 'client' dir)
NOTES_GRPC_ADDRESS="notes-2dbml6flea-wl.a.run.app:443" go run main.go
```

## References
- [GRPC status codes](https://developers.google.com/maps-booking/reference/grpc-api/status_codes)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [Language Guide](https://developers.google.com/protocol-buffers/docs/proto3)
- [GCP Managed Service Endpoints](https://github.com/GoogleCloudPlatform/golang-samples/tree/main/endpoints/getting-started-grpc)
- [Cloud Run gRPC](https://cloud.google.com/run/docs/triggering/grpc)
- [Regenerate gRPC code](https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code)
- [New Go API for Protobuf](https://go.dev/blog/protobuf-apiv2)
- [Artifact Registry](https://cloud.google.com/artifact-registry/docs)