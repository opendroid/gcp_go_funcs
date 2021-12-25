package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/google/uuid"
	notespb "github.com/opendroid/gcp_go_funcs/grpc_tests/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"strings"
	"time"
)

const (
	AuthorID            = "3627fb6e-8f9c-4418-adea-e66efb467ecd"
	TimeOut             = time.Second * 120
	Host                = "notes-2dbml6flea-wl.a.run.app" // notes-2dbml6flea-wl.a.run.app
	HostPort            = Host + ":443"
	GCPCloudRunEndpoint = "run.app"
)

// main tests a Notes Client
func main() {
	// Set up a connection to the server.
	var opts []grpc.DialOption
	hostPort := HostPort
	if addr := os.Getenv("NOTES_GRPC_ADDRESS"); addr != "" {
		hostPort = addr
	}
	fmt.Printf(`{"severity": "DEBUG", "method": "client-main", "text": "trying host", "host": "%s"}`, hostPort)
	// Note: gRPC client app must handle TLS, per https://ahmet.im/blog/grpc-auth-cloud-run/
	if strings.Contains(hostPort, GCPCloudRunEndpoint) {
		cred := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	} else {
		// For local host testing.
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	// opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(hostPort, opts...)
	if err != nil {
		fmt.Printf(`{"severity": "ERROR", "method": "client-main", "error": %q, "text": "Failed to dial"}`, err.Error())
		return
	}
	defer func() { _ = conn.Close() }()

	c := notespb.NewNotesServiceClient(conn)
	createNote(c)
	getNotesByAuthor(c, AuthorID)
}

// createNote helper method to test GRPC note creation
func createNote(c notespb.NotesServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	ans, err := c.CreateNote(ctx, createNoteRequest())

	if err != nil {
		fmt.Printf(`{"severity": "DEBUG", "ERROR": "createNote", "response": %q}`, err.Error())
		return
	}

	fmt.Printf(`{"severity": "DEBUG", "method": "createNote", "text": "%s"}`,
		ans.GetErrMessage())
}

// getNotesByAuthor fetches all notes by an author
func getNotesByAuthor(c notespb.NotesServiceClient, author string) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	ans, err := c.GetNotesByAuthor(ctx, &notespb.GetNotesByAuthorRequest{Author: AuthorID})
	if err != nil {
		fmt.Printf(`{"severity": "ERROR", "method": "getNotesByAuthor", "response": %q}`, err.Error())
		return
	}
	// print all notes fetched
	notes := ans.GetNotes()
	if len(notes) == 0 {
		fmt.Printf(`{"severity": "DEBUG", "method": "getNotesByAuthor", "text": "no notes by author", "author": "%s"}`, author)
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
		fmt.Printf(`{"severity": "DEBUG", "method": "getNotesByAuthor", "note": %d, "id": "%s", "text": "%s", "at": "%s", "locations": %s}`,
			i+1, n.GetId(), n.GetText(), n.CreatedAt.AsTime(), locations)
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
