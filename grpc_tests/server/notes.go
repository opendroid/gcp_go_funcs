package main

import (
	"context"
	"fmt"
	pb "github.com/opendroid/gcp_go_funcs/grpc_tests/notes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// Define notes server
type notesServer struct {
	pb.UnimplementedNotesServiceServer
}

type Location struct {
	Latitude  float64
	Longitude float64
	At        time.Time
}

// Note extracted fields from pb.Note
type Note struct {
	Id          string
	CreatedAt   time.Time
	LastUpdated time.Time
	Author      string
	Text        string
	Locations   []Location
}

var (
	notes map[string][]Note
)

func init() {
	notes = make(map[string][]Note)
}

// CreateNote add a note to the local map
func (s *notesServer) CreateNote(_ context.Context, request *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error) {
	locations := request.Note.GetLocations()
	if len(locations) > 0 {
		for _, loc := range locations {
			lat := loc.Latitude
			long := loc.Longitude
			at := loc.At.AsTime()
			fmt.Printf(`{"method": "CreateNote", "noteID": "%s", "author": "%s", "lat": %f, "long": %f, "at": "%s"}`,
				request.GetNote().GetId(), request.GetAuthor(), lat, long, at.String())
		}
	} else {
		fmt.Printf(`{"method": "CreateNote", "noteID": "%s", "author": "%s"}`,
			request.GetNote().GetId(), request.GetAuthor())
	}

	// Get notes fields, copy from protobuf to local map
	author := request.GetAuthor()
	note := Note{
		Id:          request.GetNote().GetId(),
		CreatedAt:   request.Note.CreatedAt.AsTime(),
		LastUpdated: request.Note.LastUpdated.AsTime(),
		Author:      request.GetAuthor(),
		Text:        request.Note.GetText(),
		Locations:   make([]Location, len(request.Note.GetLocations())),
	}
	for i, loc := range request.Note.Locations {
		note.Locations[i].At = loc.At.AsTime()
		note.Locations[i].Latitude = loc.GetLatitude()
		note.Locations[i].Longitude = loc.GetLongitude()
	}

	if n, ok := notes[author]; !ok {
		n = make([]Note, 1)
		notes[author] = append(n, note)
	} else {
		notes[author] = append(n, note)
	}
	msg := "OK"
	response := &pb.CreateNoteResponse{ErrMessage: &msg}
	return response, nil
}

// GetNote that is a specific UUID and by Author
func (s *notesServer) GetNote(_ context.Context, request *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	fmt.Printf(`{"method": "GetNote"}`)
	return nil, nil
}

// GetNotesByAuthor all notes by the Author
func (s *notesServer) GetNotesByAuthor(_ context.Context, request *pb.GetNotesByAuthorRequest) (*pb.GetNotesByAuthorResponse, error) {
	author := request.GetAuthor()
	if author == "" {
		fmt.Printf(`{"method": "GetNotesByAuthor", "text": "need author UUID"}`)
		return nil, fmt.Errorf("GetNotesByAuthor: need author UUID")
	}
	fmt.Printf(`{"method": "GetNotesByAuthor", "author": "%s"}`, author)
	// Copy all notes data. No sync operation
	if n, ok := notes[author]; ok && len(n) > 0 {
		nptrs := make([]*pb.Note, len(n))
		for i := 0; i < len(n); i++ {
			nptrs[i] = new(pb.Note) // Create a new note
			nptrs[i].Id = n[i].Id
			nptrs[i].Author = n[i].Author
			nptrs[i].Text = n[i].Text
			nptrs[i].CreatedAt = timestamppb.New(n[i].CreatedAt)
			nptrs[i].LastUpdated = timestamppb.New(n[i].LastUpdated)
			nptrs[i].Locations = make([]*pb.Location, len(n[i].Locations)) // create space for locations
			for j := 0; j < len(n[i].Locations); j++ {
				nptrs[i].Locations[j] = &pb.Location{
					Latitude:  n[i].Locations[j].Latitude,
					Longitude: n[i].Locations[j].Longitude,
					At:        timestamppb.New(n[i].Locations[j].At),
				} // Add each location
			}
		}
		return &pb.GetNotesByAuthorResponse{Notes: nptrs}, nil
	}
	return nil, fmt.Errorf("GetNotesByAuthor: no notes by author %s", author)
}
