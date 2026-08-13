package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	notespb "github.com/opendroid/gcp_go_funcs/grpc_tests/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	AuthorID            = "3627fb6e-8f9c-4418-adea-e66efb467ecd"
	TimeOut             = time.Second * 10
	GCPCloudRunHost     = "notes-2dbml6flea-uc.a.run.app"
	GCPCloudRunEndpoint = "run.app"
)

// main tests a Notes Client
func main() {
	// Set up a connection to the server.
	var opts []grpc.DialOption
	hostPort := GCPCloudRunHost
	if addr := os.Getenv("NOTES_GRPC_ADDRESS"); addr != "" {
		hostPort = addr
	}
	logger.Debug("attempting connection to host", "method", "client-main", "host", hostPort)

	// Note: gRPC client app must handle TLS, per https://ahmet.im/blog/grpc-auth-cloud-run/
	// Check if run.app supplied TLS certificate is trusted
	if strings.Contains(hostPort, GCPCloudRunEndpoint) {
		opts = append(opts, grpc.WithAuthority(hostPort))
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			logger.Error("failed to load system root CA cert pool", "method", "client-main", "error", err)
			return
		}
		cred := credentials.NewTLS(&tls.Config{RootCAs: systemRoots})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	} else {
		// Insecure for localhost:8080 testing.
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		logger.Info("proceeding without TLS", "method", "client-main")
	}

	conn, err := grpc.NewClient(hostPort, opts...)
	if err != nil {
		logger.Error("failed to dial gRPC host", "method", "client-main", "error", err, "host", hostPort)
		return
	}
	defer func() { _ = conn.Close() }()
	c := notespb.NewNotesServiceClient(conn) // Create a client
	createNote(c)                            // Call CreateNote API on 'c'
	getNotesByAuthor(c, AuthorID)
}

// createNote helper method to test GRPC note creation
func createNote(c notespb.NotesServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	ans, err := c.CreateNote(ctx, createNoteRequest())

	if err != nil {
		logger.ErrorContext(ctx, "failed to create note", "method", "createNote", "error", err)
		return
	}
	logger.InfoContext(ctx, "createNote response", "method", "createNote", "response", ans.GetErrMessage())
}

// getNotesByAuthor fetches all notes by an author
func getNotesByAuthor(c notespb.NotesServiceClient, author string) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	ans, err := c.GetNotesByAuthor(ctx, &notespb.GetNotesByAuthorRequest{Author: AuthorID})
	if err != nil {
		logger.ErrorContext(ctx, "failed to get notes by author", "method", "getNotesByAuthor", "error", err, "author", author)
		return
	}
	// print all notes fetched
	notes := ans.GetNotes()
	if len(notes) == 0 {
		logger.InfoContext(ctx, "no notes found by author", "method", "getNotesByAuthor", "author", author)
		return
	}
	for i, n := range notes {
		logger.InfoContext(
			ctx,
			"retrieved note",
			"method", "getNotesByAuthor",
			"index", i+1,
			"id", n.GetId(),
			"text", n.GetText(),
			"created_at", n.CreatedAt.AsTime().String(),
			"location_count", len(n.GetLocations()),
		)
	}
}

// createNoteRequest a test note request
func createNoteRequest() *notespb.CreateNoteRequest {
	// Create a note request
	ts := timestamppb.Now()
	sf := &notespb.Location{
		Latitude:  37.773972, // San Fran
		Longitude: -122.431297,
		At:        ts,
	}
	note := &notespb.Note{
		Id:          uuid.NewString(), // RFC-4122
		CreatedAt:   ts,
		LastUpdated: ts,
		Author:      AuthorID,
		Locations:   []*notespb.Location{sf},
		Text:        "This is my test note, created in San Fran.",
	}

	return &notespb.CreateNoteRequest{
		Author: AuthorID,
		Note:   note,
	}
}
