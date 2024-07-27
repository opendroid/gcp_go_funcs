# Dockerfile References: https://docs.docker.com/engine/reference/builder/
# Sample https://cloud.google.com/run/docs/quickstarts/build-and-deploy?_ga=2.91290522.-1679093051.1593441137

# Start from the latest golang base image
FROM golang:1.22 as builder

# Add Maintainer Info
LABEL maintainer="Open Web <plutoapps@outlook.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and  go.sum files
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o grpc_test_server

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:latest
# RUN apk add --no-cache ca-certificates
# RUN apk update && apk add ca-certificates && update-ca-certificates && rm -rf /var/cache/apk/*
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/grpc_test_server /grpc_test_server

# Expose port 8080 to the outside world
# EXPOSE 8080

# Command to run the executable
CMD ["/grpc_test_server"]