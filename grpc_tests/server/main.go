package main

import (
	"fmt"
	pb "github.com/opendroid/gcp_go_funcs/grpc_tests/notes"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// main host the gRPC server
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		m := fmt.Sprintf(`{"severity": "WARNING", "method": "server-main", "port": "%s", "message": "default port"}`, port)
		fmt.Println(m)
	}

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		m := fmt.Sprintf(`{"severity": "ERROR", "method": "server-main", "error": "%s", "message": "message on faileing to listen"}`, err.Error())
		fmt.Println(m)
		return
	}

	s := grpc.NewServer()
	pb.RegisterNotesServiceServer(s, &notesServer{})
	m := fmt.Sprintf(`{"severity": "DEBUG", "method": "server-main", "message": "gRPC listening at %v"}`, l.Addr())
	fmt.Println(m)

	// Trap SIGTERM,  test  by: docker kill --signal="SIGTERM"
	// https://cloud.google.com/run/docs/samples/cloudrun-sigterm-handler
	var wg sync.WaitGroup
	var sig os.Signal
	wg.Add(1)
	go func() {
		defer wg.Done()
		term := make(chan os.Signal, 1) // don't block the notifier
		signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
		sig = <-term     // Wait on term
		s.GracefulStop() // Gracefully shutdown
	}()

	// Start the server
	if err := s.Serve(l); err != nil {
		m := fmt.Sprintf(`{"severity": "ERROR", "method": "server-main", "error": "%s", "message": "exiting on failure to serve"}`, err.Error())
		fmt.Println(m)
	}
	wg.Wait()
	m = fmt.Sprintf(`{"severity": "INFO", "method": "server-main", "message": "%v signal received exiting"}`, sig)
	fmt.Println(m)
}
