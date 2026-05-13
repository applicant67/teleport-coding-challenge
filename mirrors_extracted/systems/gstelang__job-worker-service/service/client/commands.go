package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	pb "github.com/gstelang/job-worker-service/service/proto/v1"
)

func startJob(client pb.JobWorkerClient, commandName string, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.StartJobRequest{
		CommandName: commandName,
		CommandArgs: args,
	}

	res, err := client.StartJob(ctx, req)
	if err != nil {
		log.Fatalf("Failed to start job: %v", err)
		return err
	}

	fmt.Printf("Job started with ID: %s\n", res.JobId)
	return nil
}

func stopJob(client pb.JobWorkerClient, jobID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.StopJobRequest{
		JobId: jobID,
	}

	res, err := client.StopJob(ctx, req)
	if err != nil {
		log.Fatalf("Failed to stop job: %v", err)
		return err
	}

	fmt.Printf("Job stopped: %v\n", res.Message)
	return nil
}

func queryJob(client pb.JobWorkerClient, jobID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.QueryJobRequest{
		JobId: jobID,
	}

	res, err := client.QueryJob(ctx, req)
	if err != nil {
		log.Fatalf("Failed to query job: %v", err)
		return err
	}

	// beautify output
	resJSON, _ := json.MarshalIndent(res, "", "  ")
	fmt.Printf("Job query result: %s\n", resJSON)
	return nil
}

func streamLogs(client pb.JobWorkerClient, jobID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	req := &pb.StreamLogsRequest{
		JobId: jobID,
	}

	stream, err := client.StreamLogs(ctx, req)
	if err != nil {
		log.Fatalf("Failed to stream logs: %v", err)
		return err
	}

	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				log.Fatalf("Error receiving log: %v", err)
				return
			}

			// printing string representation.
			// adjust if binary
			fmt.Printf("Log entry: %s\n", res.Message)
		}
	}()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	defer func() {
		cancel()
		signal.Stop(sigchan)
	}()
	<-sigchan
	return nil
}
