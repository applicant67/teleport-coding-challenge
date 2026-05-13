package main

import (
	"context"
	"fmt"
	"slices"
	"sync"

	pb "github.com/gstelang/job-worker-service/service/proto/v1"
	"github.com/gstelang/job-worker-service/worker"
)

type server struct {
	pb.UnimplementedJobWorkerServer
	jobManager *worker.JobManager
	clients    map[string][]pb.JobWorker_StreamLogsServer // jobID -> client streams
	streamsMu  sync.Mutex
}

func (s *server) StartJob(ctx context.Context, req *pb.StartJobRequest) (*pb.StartJobResponse, error) {

	if err := authorize(ctx, "StopJob"); err != nil {
		return nil, err
	}
	command := convertToCommand(req)

	jobID, err := s.jobManager.Start(command)

	if err == nil {
		return &pb.StartJobResponse{JobId: jobID, Message: "start_job"}, nil
	} else {
		return &pb.StartJobResponse{JobId: "", Message: fmt.Sprintf("%v", err)}, err
	}
}

func (s *server) StopJob(ctx context.Context, req *pb.StopJobRequest) (*pb.StopJobResponse, error) {

	if err := authorize(ctx, "StopJob"); err != nil {
		return nil, err
	}
	jobID := req.GetJobId()

	ok, err := s.jobManager.Stop(jobID)
	if err != nil {
		return &pb.StopJobResponse{Success: false, Message: fmt.Sprintf("Failed to stop job: %v", err)}, err
	}

	return &pb.StopJobResponse{Success: ok, Message: "Stopped job!"}, nil
}

func (s *server) QueryJob(ctx context.Context, req *pb.QueryJobRequest) (*pb.QueryJobResponse, error) {
	if err := authorize(ctx, "QueryJob"); err != nil {
		return nil, err
	}
	jobID := req.GetJobId()

	jobSummary, err := s.jobManager.Query(jobID)
	if err != nil {
		return &pb.QueryJobResponse{Message: fmt.Sprintf("Failed to query job: %v", err)}, err
	}

	return &pb.QueryJobResponse{
		Status:   jobSummary.Status.String(),
		Pid:      int32(jobSummary.PID),
		ExitCode: int32(jobSummary.ExitCode),
		Signal:   int32(jobSummary.Signal),
		Message:  "Sucess!",
	}, nil
}

func (s *server) StreamLogs(req *pb.StreamLogsRequest, stream pb.JobWorker_StreamLogsServer) error {

	if err := authorize(stream.Context(), "StreamLogs"); err != nil {
		return err
	}

	jobID := req.GetJobId()

	// Register this client
	s.addClient(jobID, stream)
	defer s.removeClient(jobID, stream)

	// Fetch existing logs and create/get the log channel
	existingLogs, logChannel, err := s.jobManager.Stream(stream.Context(), jobID)
	if err != nil {
		return fmt.Errorf("failed to stream logs: %w", err)
	}

	// Send existing logs to the client
	// TODO: Test this with high volume and large concurrent connections with real time updates.
	for _, logEntry := range existingLogs {
		if err := stream.Send(&pb.StreamLogsResponse{Message: logEntry}); err != nil {
			return fmt.Errorf("failed to send existing log entry: %w", err)
		}
	}

	for {
		select {
		case logEntry, ok := <-logChannel:
			if !ok {
				fmt.Printf("Log channel closed for jobID %s", jobID)
				return nil
			}
			if err := s.sendLogToClient(jobID, logEntry); err != nil {
				fmt.Printf("Failed to send log entry for jobID %s: %v", jobID, err)
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (s *server) addClient(jobID string, stream pb.JobWorker_StreamLogsServer) {
	s.streamsMu.Lock()
	defer s.streamsMu.Unlock()
	if s.clients == nil {
		s.clients = make(map[string][]pb.JobWorker_StreamLogsServer)
	}
	s.clients[jobID] = append(s.clients[jobID], stream)
}

func (s *server) removeClient(jobID string, stream pb.JobWorker_StreamLogsServer) {
	s.streamsMu.Lock()
	defer s.streamsMu.Unlock()
	if clients, ok := s.clients[jobID]; ok {
		s.clients[jobID] = slices.DeleteFunc(clients, func(currentStream pb.JobWorker_StreamLogsServer) bool {
			return currentStream == stream
		})
	}
}

func (s *server) sendLogToClient(jobID string, logEntry []byte) error {
	s.streamsMu.Lock()
	defer s.streamsMu.Unlock()
	clients, ok := s.clients[jobID]
	if !ok {
		return fmt.Errorf("no clients found for jobID %s", jobID)
	}
	for _, client := range clients {
		if err := client.Send(&pb.StreamLogsResponse{Message: logEntry}); err != nil {
			fmt.Printf("Error sending log to client for jobID %s: %v", jobID, err)
			s.removeClient(jobID, client)
			return err
		}
	}
	return nil
}
