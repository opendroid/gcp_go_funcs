package main

import (
	"context"
	"fmt"
	pb "github.com/opendroid/gcp_go_funcs/grpc_tests/notes"
)

// Define notes server
type notesServer struct {
	pb.UnimplementedNotesServiceServer
}

var (
	notes map[string][]pb.Note
)

func init() {
	notes = make(map[string][]pb.Note)
}

// CreateNote add a note to the local map
func (s *notesServer) CreateNote(ctx context.Context, request *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error) {
	fmt.Printf(`{"method": "CreateNote", "noteID": "%s", "author": "%s"}`,
		request.GetNote().GetId(), request.GetAuthor())
	if n, ok := notes[request.GetAuthor()]; !ok {
		n = make([]pb.Note, 1)
		notes[request.GetAuthor()] = append(n, *request.GetNote())
	} else {
		notes[request.GetAuthor()] = append(n, *request.GetNote())
	}
	msg := "OK"
	response := &pb.CreateNoteResponse{ErrMessage: &msg}
	return response, nil
}

func (s *notesServer) GetNote(ctx context.Context, request *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	fmt.Printf(`{"method": "GetNote"}`)
	return nil, nil
}

func (s *notesServer) GetNotesByAuthor(ctx context.Context, request *pb.GetNotesByAuthorRequest) (*pb.GetNotesByAuthorResponse, error) {
	fmt.Printf(`{"method": "GetNote"}`)
	return nil, nil
}
