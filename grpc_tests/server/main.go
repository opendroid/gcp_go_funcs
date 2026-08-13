package main

import (
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
		logger.Warn("PORT environment variable not set, using default port", "method", "server-main", "port", port)
	}

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("failed to listen on port", "method", "server-main", "error", err, "port", port)
		return
	}

	s := grpc.NewServer()
	pb.RegisterNotesServiceServer(s, &notesServer{})
	logger.Info("gRPC server listening", "method", "server-main", "address", l.Addr().String())

	// Trap SIGTERM, test by: docker kill --signal="SIGTERM"
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
		logger.Error("server exited on failure to serve", "method", "server-main", "error", err)
	}
	wg.Wait()
	logger.Info("shutdown signal received, exiting", "method", "server-main", "signal", sig.String())
}
