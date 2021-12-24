package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/opendroid/gcp_go_funcs/grpc_tests/notes"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const (
	AuthorID = "3627fb6e-8f9c-4418-adea-e66efb467ecd"
)

// main tests a Notes Client
func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf(`{"method": "client-main", "error": "%s", "text": "Failed to dial"}`, err.Error())
		return
	}
	defer func() { _ = conn.Close() }()

	c := pb.NewNotesServiceClient(conn)
	createNote(c)
	getNotesByAuthor(c, AuthorID)
}

// createNote helper method to test GRPC note creation
func createNote(c pb.NotesServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ans, err := c.CreateNote(ctx, createNoteRequest())

	if err != nil {
		fmt.Printf(`{"method": "createNote", "response": "%s"}`, err.Error())
		return
	}

	fmt.Printf(`{"method": "createNote", "response": "%s"}`,
		ans.GetErrMessage())
}

// getNotesByAuthor fetches all notes by an author
func getNotesByAuthor(c pb.NotesServiceClient, author string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ans, err := c.GetNotesByAuthor(ctx, &pb.GetNotesByAuthorRequest{Author: AuthorID})
	if err != nil {
		fmt.Printf(`{"method": "getNotesByAuthor", "response": "%s"}`, err.Error())
		return
	}
	// print all notes fetched
	notes := ans.GetNotes()
	if len(notes) == 0 {
		fmt.Printf(`{"method": "getNotesByAuthor", "text": "no notes by author", "author": "%s"}`, author)
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
		fmt.Printf(`{"method": "getNotesByAuthor", "note": %d, "text": "%s", "at": "%s", "locations": %s}`,
			i+1, n.GetText(), n.CreatedAt.AsTime(), locations)
	}
}

// createNoteRequest a test note request
func createNoteRequest() *pb.CreateNoteRequest {
	// Create a note request
	ts := timestamppb.Now()
	location := &pb.Location{
		Latitude:  37.773972,
		Longitude: -122.431297,
		At:        ts,
	}
	note := &pb.Note{
		Id:          uuid.NewString(), // RFC-4122
		CreatedAt:   ts,
		LastUpdated: ts,
		Author:      AuthorID,
		Locations:   []*pb.Location{location},
		Text:        "This is my test note, created in San Fran.",
	}

	return &pb.CreateNoteRequest{
		Author: AuthorID,
		Note:   note,
	}
}
