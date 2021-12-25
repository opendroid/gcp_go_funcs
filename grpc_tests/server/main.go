package main

import (
	"fmt"
	pb "github.com/opendroid/gcp_go_funcs/grpc_tests/notes"
	"google.golang.org/grpc"
	"net"
	"os"
)

// main host the gRPC server
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		m := fmt.Sprintf(`{"severity": "WARNING", "method": "server-main", "port": "%s", "text": "default port"}`, port)
		fmt.Println(m)
	}

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		m := fmt.Sprintf(`{"severity": "ERROR", "method": "main", "error": "%s", "text": "Exit on faileing to listen"}`, err.Error())
		fmt.Println(m)
		return
	}

	s := grpc.NewServer()
	pb.RegisterNotesServiceServer(s, &notesServer{})
	m := fmt.Sprintf(`{"severity": "DEBUG", "method": "main", "text": "gRPC listening at %v"}`, l.Addr())
	fmt.Println(m)
	if err := s.Serve(l); err != nil {
		m := fmt.Sprintf(`{"severity": "ERROR", "method": "main", "error": "%s", "text": "exiting on failure to serve"}`, err.Error())
		fmt.Println(m)
	}
}
