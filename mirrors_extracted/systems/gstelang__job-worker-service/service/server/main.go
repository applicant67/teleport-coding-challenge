package main

import (
	"log"
	"net"
	"time"

	pb "github.com/gstelang/job-worker-service/service/proto/v1"
	"github.com/gstelang/job-worker-service/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {

	// TLS setup
	creds := setupTLS()

	// Create a gRPC server object with TLS creds
	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 30 * time.Minute, // Disconnect after 30 minutes of inactivity
			Time:              10 * time.Minute, // Send keepalive pings every 10 minutes
			Timeout:           20 * time.Second, // Wait 20 seconds for keepalive ping ack
		}),
	)

	// Initialize Server
	pb.RegisterJobWorkerServer(s, &server{
		// use worker library's JobManager
		jobManager: worker.NewJobManager(),
		clients:    make(map[string][]pb.JobWorker_StreamLogsServer),
	})

	// Create a listener on TCP port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// Serve
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
