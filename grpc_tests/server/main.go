package main

import (
	"fmt"
	pb "github.com/opendroid/gcp_go_funcs/grpc_tests/notes"
	"google.golang.org/grpc"
	"net"
)

// main host the gRPC server
func main() {
	l, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		fmt.Printf(`{"method": "main", "error": "%s", "text": "Failed to listen"}`, err.Error())
		return
	}

	s := grpc.NewServer()
	pb.RegisterNotesServiceServer(s, &notesServer{})
	fmt.Printf(`{"method": "main", "text": "gRPC listening at %v"}`, l.Addr())
	if err := s.Serve(l); err != nil {
		fmt.Printf(`{"method": "main", "error": "%s", "text": "Failed to serve"}`, err.Error())
	}
}
