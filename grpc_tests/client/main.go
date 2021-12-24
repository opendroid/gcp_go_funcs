package main

import (
	"context"
	"fmt"
	pb "github.com/opendroid/gcp_go_funcs/grpc_tests/notes"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051")
	if err != nil {
		fmt.Printf(`{"method": "client-main", "error": "%s", "text": "Failed to dial"}`, err.Error())
		return
	}
	defer func() { _ = conn.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := pb.NewNotesServiceClient(conn)

	ts := timestamppb.Now()
	note := &pb.Note{
		Id:          "DD5738EF-0F1D-407F-A1E0-B2F917D2AD08",
		CreatedAt:   ts,
		LastUpdated: ts,
		Author:      "3627FB6E-8F9C-4418-ADEA-E66EFB467ECD",
		Location:    nil,
		Text:        "This is my test note.",
	}
	ans, err := c.CreateNote(ctx,
		&pb.CreateNoteRequest{
			Author: "3627FB6E-8F9C-4418-ADEA-E66EFB467ECD",
			Note:   note,
		})

	fmt.Printf(`{"method": "client-main", "response": "%s"}`,
		ans.GetErrMessage())
}
