package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
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
	m := fmt.Sprintf(`{"severity": "DEBUG", "method": "client-main", "text": "trying host", "host": "%s"}`, hostPort)
	fmt.Println(m)
	// Note: gRPC client app must handle TLS, per https://ahmet.im/blog/grpc-auth-cloud-run/
	// Check if  run.app supplied TLS certificate is trusted
	if strings.Contains(hostPort, GCPCloudRunEndpoint) {
		opts = append(opts, grpc.WithAuthority(hostPort))
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			m := fmt.Sprintf(`{"severity": "ERROR", "method": "client-main", "error": %q, "text": "Failed to load system root CA cert pool"}`, err.Error())
			fmt.Println(m)
			return
		}
		cred := credentials.NewTLS(&tls.Config{RootCAs: systemRoots})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	} else {
		// Insecure for localhost:8080 testing.
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		fmt.Println(`{"severity": "INFO", "method": "client-main", "message": "Proceeding without TLS"}`)
	}
	// opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.NewClient(hostPort, opts...)
	if err != nil {
		m := fmt.Sprintf(`{"severity": "ERROR", "method": "client-main", "error": %q, "text": "Failed to dial"}`, err.Error())
		fmt.Println(m)
		return
	}
	defer func() { _ = conn.Close() }()
	c := notespb.NewNotesServiceClient(conn) // Create a client
	createNote(c)                            // Call CreatesNot API on 'c'
	getNotesByAuthor(c, AuthorID)
}

// createNote helper method to test GRPC note creation
func createNote(c notespb.NotesServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	ans, err := c.CreateNote(ctx, createNoteRequest())

	if err != nil {
		m := fmt.Sprintf(`{"severity": "DEBUG", "ERROR": "createNote", "response": %q}`, err.Error())
		fmt.Println(m)
		return
	}
	m := fmt.Sprintf(`{"severity": "DEBUG", "method": "createNote", "text": "%s"}`, ans.GetErrMessage())
	fmt.Println(m)
}

// getNotesByAuthor fetches all notes by an author
func getNotesByAuthor(c notespb.NotesServiceClient, author string) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	ans, err := c.GetNotesByAuthor(ctx, &notespb.GetNotesByAuthorRequest{Author: AuthorID})
	if err != nil {
		m := fmt.Sprintf(`{"severity": "ERROR", "method": "getNotesByAuthor", "message": %q}`, err.Error())
		fmt.Println(m)
		return
	}
	// print all notes fetched
	notes := ans.GetNotes()
	if len(notes) == 0 {
		m := fmt.Sprintf(`{"severity": "DEBUG", "method": "getNotesByAuthor", "message": "no notes by author", "author": "%s"}`, author)
		fmt.Println(m)
		return
	}
	for i, n := range notes {
		locations := "["
		if loc := n.GetLocations(); len(loc) > 0 {
			for _, l := range loc {
				locations += fmt.Sprintf(`{"lat": %f, "long": %f, "at": "%s"}`, l.Latitude, l.Longitude, l.At.AsTime())
			}
		}
		locations += "]"
		m := fmt.Sprintf(`{"severity": "DEBUG", "method": "getNotesByAuthor", "note": %d, "id": "%s", "text": "%s", "at": "%s", "locations": %s}`,
			i+1, n.GetId(), n.GetText(), n.CreatedAt.AsTime(), locations)
		fmt.Println(m)
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
