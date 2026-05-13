// Package server implements the gRPC server for the JobWorker service.
package server

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiv1 "github.com/benmoss/job-worker-service/pkg/api/jobworker/v1"
	"github.com/benmoss/job-worker-service/pkg/worker"
)

// Server implements the JobWorkerService gRPC service.
type Server struct {
	apiv1.UnimplementedJobWorkerServiceServer
	worker *worker.Worker
}

// New creates a new gRPC server wrapping the given worker.
func New(w *worker.Worker) *Server {
	return &Server{
		worker: w,
	}
}

// StartJob creates and executes a new job.
func (s *Server) StartJob(ctx context.Context, req *apiv1.StartJobRequest) (*apiv1.StartJobResponse, error) {
	if req.Command == "" {
		return nil, status.Error(codes.InvalidArgument, "command is required")
	}

	var limits *worker.ResourceLimits
	if req.Limits != nil {
		limits = &worker.ResourceLimits{
			CPUMax:         req.Limits.CpuMax,
			MemoryMaxBytes: req.Limits.MemoryMaxBytes,
			IOMaxReadBPS:   req.Limits.IoMaxReadBps,
			IOMaxWriteBPS:  req.Limits.IoMaxWriteBps,
		}
	}

	id, err := s.worker.Start(worker.Job{
		Command:        req.Command,
		Args:           req.Args,
		ResourceLimits: limits,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "start job: %v", err)
	}

	return &apiv1.StartJobResponse{
		JobId: id,
	}, nil
}

// StopJob terminates a job and all its child processes.
func (s *Server) StopJob(ctx context.Context, req *apiv1.StopJobRequest) (*apiv1.StopJobResponse, error) {
	if err := s.worker.Stop(req.JobId); err != nil {
		if errors.Is(err, worker.ErrJobNotFound) {
			return nil, status.Error(codes.NotFound, "job not found")
		}
		return nil, status.Errorf(codes.Internal, "stop job: %v", err)
	}

	return &apiv1.StopJobResponse{}, nil
}

// GetStatus returns the current state of a job.
func (s *Server) GetStatus(ctx context.Context, req *apiv1.GetStatusRequest) (*apiv1.GetStatusResponse, error) {
	jobStatus, err := s.worker.Status(req.JobId)
	if err != nil {
		if errors.Is(err, worker.ErrJobNotFound) {
			return nil, status.Error(codes.NotFound, "job not found")
		}
		return nil, status.Errorf(codes.Internal, "get status: %v", err)
	}

	return &apiv1.GetStatusResponse{
		JobId:    jobStatus.ID,
		State:    jobStateToProto(jobStatus.State),
		ExitCode: int32(jobStatus.ExitCode),
		Command:  jobStatus.Command,
		Args:     jobStatus.Args,
	}, nil
}

// GetLogs streams the job's combined stdout/stderr.
func (s *Server) GetLogs(req *apiv1.GetLogsRequest, stream apiv1.JobWorkerService_GetLogsServer) error {
	reader, err := s.worker.Logs(stream.Context(), req.JobId, req.Follow)
	if err != nil {
		if errors.Is(err, worker.ErrJobNotFound) {
			return status.Error(codes.NotFound, "job not found")
		}
		return status.Errorf(codes.Internal, "get logs: %v", err)
	}
	defer reader.Close()

	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if err := stream.Send(&apiv1.GetLogsResponse{
				Data: buf[:n],
			}); err != nil {
				return status.Errorf(codes.Internal, "send logs: %v", err)
			}
		}

		if err != nil {
			if err == io.EOF {
				return nil
			}
			if errors.Is(err, context.Canceled) {
				return status.Error(codes.Canceled, "client disconnected")
			}
			return status.Errorf(codes.Internal, "read logs: %v", err)
		}
	}
}

// jobStateToProto converts a worker.JobState to the proto JobState.
func jobStateToProto(state worker.JobState) apiv1.JobState {
	switch state {
	case worker.JobStateRunning:
		return apiv1.JobState_JOB_STATE_RUNNING
	case worker.JobStateCompleted:
		return apiv1.JobState_JOB_STATE_COMPLETED
	case worker.JobStateStopping:
		return apiv1.JobState_JOB_STATE_STOPPING
	case worker.JobStateStopped:
		return apiv1.JobState_JOB_STATE_STOPPED
	case worker.JobStateKilled:
		return apiv1.JobState_JOB_STATE_KILLED
	default:
		return apiv1.JobState_JOB_STATE_UNSPECIFIED
	}
}
