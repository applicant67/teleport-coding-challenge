package main

import (
	"flag"
	"log"

	pb "github.com/gstelang/job-worker-service/service/proto/v1"
	"google.golang.org/grpc"
)

func main() {
	// 1. Define command-line flags
	command := flag.String("command", "", "Command to execute (start, stop, query, stream)")
	jobID := flag.String("jobid", "", "Job ID for stop, query, or stream")
	flag.Parse()

	// 2. Connection setup
	creds := setupTLS()
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()
	client := pb.NewJobWorkerClient(conn)

	// 3. Handle commands
	switch *command {
	case "start":
		remainingArgs := flag.Args()
		commandName := remainingArgs[0]
		commandArgs := remainingArgs[1:]
		startJob(client, commandName, commandArgs)
	case "stop":
		stopJob(client, *jobID)
	case "query":
		queryJob(client, *jobID)
	case "stream":
		streamLogs(client, *jobID)
	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}
