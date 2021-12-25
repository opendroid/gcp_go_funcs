package main

import (
	"context"
	"fmt"
	notespb "github.com/opendroid/gcp_go_funcs/grpc_tests/notes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// Define notes server
type notesServer struct {
	notespb.UnimplementedNotesServiceServer
}

// Location a mapper for notespb.Location
type Location struct {
	Latitude  float64   `json:"latitude,omitempty"`
	Longitude float64   `json:"longitude,omitempty"`
	At        time.Time `json:"at,omitempty"`
}

// Note extracted fields from notespb.Note. All fields are copied and saved here
type Note struct {
	Id          string     `json:"id,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	LastUpdated time.Time  `json:"last_updated,omitempty"`
	Author      string     `json:"author,omitempty"`
	Text        string     `json:"text,omitempty"`
	Locations   []Location `json:"locations,omitempty"`
}

var (
	// notes is poor mans demo-data store.
	notes map[string][]Note
)

const (
	InvalidID = "00000000-0000-0000-0000-000000000000"
)

func init() {
	notes = make(map[string][]Note)
}

// CreateNote add a note to the local map
func (s *notesServer) CreateNote(_ context.Context, request *notespb.CreateNoteRequest) (*notespb.CreateNoteResponse, error) {
	locations := request.Note.GetLocations()
	if len(locations) > 0 {
		for _, loc := range locations {
			lat := loc.Latitude
			long := loc.Longitude
			at := loc.At.AsTime()
			m := fmt.Sprintf(`{"severity": "DEBUG", "method": "CreateNote", "noteID": "%s", "author": "%s", "lat": %f, "long": %f, "at": "%s"}`,
				request.GetNote().GetId(), request.GetAuthor(), lat, long, at.String())
			fmt.Println(m)
		}
	} else {
		m := fmt.Sprintf(`{"severity": "DEBUG", "method": "CreateNote", "noteID": "%s", "author": "%s"}`,
			request.GetNote().GetId(), request.GetAuthor())
		fmt.Println(m)
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
		n = make([]Note, 0)
		notes[author] = append(n, note)
	} else {
		notes[author] = append(n, note)
	}
	msg := "OK"
	response := &notespb.CreateNoteResponse{ErrMessage: &msg}
	return response, nil
}

// GetNote that is a specific UUID and by Author
func (s *notesServer) GetNote(_ context.Context, _ *notespb.GetNoteRequest) (*notespb.GetNoteResponse, error) {
	m := fmt.Sprintf(`{"severity": "DEBUG", "method": "GetNote", "text": "implement me"}`)
	fmt.Println(m)
	now := timestamppb.Now()
	return &notespb.GetNoteResponse{Note: &notespb.Note{
		Id:          InvalidID,
		CreatedAt:   now,
		LastUpdated: now,
		Author:      InvalidID,
		Locations:   []*notespb.Location{{Latitude: 37.773972, Longitude: -122.431297, At: now}},
		Text:        "Implement me",
	}}, nil
}

// GetNotesByAuthor all notes by the Author
func (s *notesServer) GetNotesByAuthor(_ context.Context, request *notespb.GetNotesByAuthorRequest) (*notespb.GetNotesByAuthorResponse, error) {
	author := request.GetAuthor()
	if author == "" {
		fmt.Println(`{"severity": "WARNING", "method": "GetNotesByAuthor", "text": "need author UUID"}`)
		return nil, fmt.Errorf("GetNotesByAuthor: need author UUID")
	}
	m := fmt.Sprintf(`{"severity": "DEBUG", "method": "GetNotesByAuthor", "author": "%s"}`, author)
	fmt.Println(m)
	// Copy all notes data. No sync operation
	if n, ok := notes[author]; ok && len(n) > 0 {
		nptrs := make([]*notespb.Note, len(n))
		for i, ni := range n { // ni is ith-note by Author
			nptrs[i] = new(notespb.Note) // Create a new note
			nptrs[i].Id = ni.Id
			nptrs[i].Author = ni.Author
			nptrs[i].Text = ni.Text
			nptrs[i].CreatedAt = timestamppb.New(ni.CreatedAt)
			nptrs[i].LastUpdated = timestamppb.New(ni.LastUpdated)
			nptrs[i].Locations = make([]*notespb.Location, len(ni.Locations)) // create space for locations
			for j, loc := range ni.Locations {                                // Add each location
				nptrs[i].Locations[j] = &notespb.Location{
					Latitude:  loc.Latitude,
					Longitude: loc.Longitude,
					At:        timestamppb.New(loc.At),
				} // Add each location
			}
		}
		return &notespb.GetNotesByAuthorResponse{Notes: nptrs}, nil
	}
	return nil, fmt.Errorf("GetNotesByAuthor: no notes by author %s", author)
}
