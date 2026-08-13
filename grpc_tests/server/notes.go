package main

import (
	"context"
	"fmt"
	notespb "github.com/opendroid/gcp_go_funcs/grpc_tests/notes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
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
	notes sync.Map // Stores [string][]Note
)

const (
	InvalidID = "00000000-0000-0000-0000-000000000000"
)

// CreateNote add a note to the local map
func (s *notesServer) CreateNote(ctx context.Context, request *notespb.CreateNoteRequest) (*notespb.CreateNoteResponse, error) {
	locations := request.Note.GetLocations()
	if len(locations) > 0 {
		for _, loc := range locations {
			lat := loc.Latitude
			long := loc.Longitude
			at := loc.At.AsTime()
			logger.DebugContext(
				ctx,
				"created note with location",
				"method", "CreateNote",
				"note_id", request.GetNote().GetId(),
				"author", request.GetAuthor(),
				"lat", lat,
				"long", long,
				"at", at.String(),
			)
		}
	} else {
		logger.DebugContext(
			ctx,
			"created note without location",
			"method", "CreateNote",
			"note_id", request.GetNote().GetId(),
			"author", request.GetAuthor(),
		)
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
	// Save data in the syncMap
	if n, ok := notes.Load(author); !ok {
		notes.Store(author, []Note{note})
	} else {
		if nd, valid := n.([]Note); valid {
			nu := append(nd, note)
			notes.Store(author, nu)
		}
	}
	msg := "OK"
	response := &notespb.CreateNoteResponse{ErrMessage: &msg}
	return response, nil
}

// GetNote that is a specific UUID and by Author
func (s *notesServer) GetNote(ctx context.Context, _ *notespb.GetNoteRequest) (*notespb.GetNoteResponse, error) {
	logger.DebugContext(ctx, "GetNote stub called", "method", "GetNote")
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
func (s *notesServer) GetNotesByAuthor(ctx context.Context, request *notespb.GetNotesByAuthorRequest) (*notespb.GetNotesByAuthorResponse, error) {
	author := request.GetAuthor()
	if author == "" {
		logger.WarnContext(ctx, "missing author UUID", "method", "GetNotesByAuthor")
		return nil, fmt.Errorf("GetNotesByAuthor: need author UUID")
	}

	// Copy all notes data.
	if n, ok := notes.Load(author); ok {
		if nd, valid := n.([]Note); valid && len(nd) > 0 { // nd, Notes data
			nptrs := make([]*notespb.Note, len(nd))
			for i, ni := range nd { // ni is ith-note by Author
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
			logger.DebugContext(
				ctx,
				"retrieved notes by author",
				"method", "GetNotesByAuthor",
				"count", len(nd),
				"author", author,
			)
			return &notespb.GetNotesByAuthorResponse{Notes: nptrs}, nil
		}
	}
	logger.DebugContext(
		ctx,
		"no notes found for author",
		"method", "GetNotesByAuthor",
		"author", author,
	)
	return nil, fmt.Errorf("GetNotesByAuthor: no notes by author %s", author)
}
