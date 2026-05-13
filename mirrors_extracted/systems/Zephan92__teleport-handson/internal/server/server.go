//go:build linux

package server

import (
	"context"
	"io"
	"log/slog"
	"runtime/debug"
	"sync"

	pb "github.com/teleport-handson/api/v1"
	"github.com/teleport-handson/internal/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// WorkerServer implements the gRPC JobWorker service.
type WorkerServer struct {
	// UnimplementedJobWorkerServer must be embedded to have forward compatible implementations.
	pb.UnimplementedJobWorkerServer
	jobs       map[worker.JobID]*worker.Job
	mu         sync.Mutex
	cgroupRoot string
}

// NewWorkerServer initializes a new WorkerServer instance.
func NewWorkerServer(cgroupRoot string) *WorkerServer {
	return &WorkerServer{
		jobs:       make(map[worker.JobID]*worker.Job),
		cgroupRoot: cgroupRoot,
	}
}

// Start handles the creation and execution of a new job.
// It validates the request, initializes the job with resource limits, and starts the process.
func (s *WorkerServer) Start(ctx context.Context, req *pb.StartRequest) (resp *pb.StartResponse, err error) {
	defer s.recoverToError(&err)

	// Security: Do not log req.Args as they may contain sensitive information (secrets/passwords).
	slog.Info("Start request", "cmd", req.GetCommand(), "dir", req.GetWorkingDirectory())

	cn, err := s.extractClientCN(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetCommand() == "" {
		return nil, status.Error(codes.InvalidArgument, "command is required")
	}
	if req.GetWorkingDirectory() == "" {
		return nil, status.Error(codes.InvalidArgument, "working_directory is required")
	}

	specs := worker.ExecutionSpecs{
		Command:          req.GetCommand(),
		Arguments:        req.GetArgs(),
		WorkingDirectory: req.GetWorkingDirectory(),
		OwnerCommonName:  cn,
	}

	job, err := worker.NewJob(specs, s.cgroupRoot)
	if err != nil {
		slog.Error("Failed to create job", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create job: %v", err)
	}
	id := job.ID

	s.mu.Lock()
	s.jobs[id] = job
	s.mu.Unlock()

	if err := job.Start(context.Background()); err != nil {
		slog.Error("Failed to start job", "id", id, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to start job: %v", err)
	}

	slog.Info("Job started successfully", "id", id)

	return &pb.StartResponse{
		JobId:     string(id),
		Status:    "RUNNING",
		CreatedAt: timestamppb.New(job.Status().StartTime),
	}, nil
}

// Stop terminates a running job.
func (s *WorkerServer) Stop(ctx context.Context, req *pb.StopRequest) (resp *pb.StopResponse, err error) {
	defer s.recoverToError(&err)

	slog.Info("Stop request", "id", req.GetJobId())

	job, err := s.authorize(ctx, req.GetJobId())
	if err != nil {
		slog.Warn("Authorization failed or job not found", "id", req.GetJobId(), "error", err)
		return nil, err
	}

	if err := job.Stop(); err != nil {
		slog.Error("Failed to stop job", "id", req.GetJobId(), "error", err)
		return nil, status.Errorf(codes.Internal, "failed to stop job: %v", err)
	}

	slog.Info("Job stopped", "id", req.GetJobId())

	return &pb.StopResponse{Message: "Job stopped successfully"}, nil
}

// Status retrieves the status of a job.
func (s *WorkerServer) Status(ctx context.Context, req *pb.StatusRequest) (resp *pb.StatusResponse, err error) {
	defer s.recoverToError(&err)

	slog.Info("Status request", "id", req.JobId)

	job, err := s.authorize(ctx, req.GetJobId())
	if err != nil {
		slog.Warn("Authorization failed or job not found", "id", req.JobId, "error", err)
		return nil, err
	}

	state := job.Status()

	resp = &pb.StatusResponse{
		JobId:     string(job.ID),
		Owner:     job.Specs.OwnerCommonName,
		Command:   job.Specs.Command,
		Args:      job.Specs.Arguments,
		Status:    s.toProtoStatus(state.Status),
		StartedAt: timestamppb.New(state.StartTime),
		ExitCode:  int32(state.ExitCode),
		Limits:    s.toProtoLimits(job.Limit),
	}

	if state.EndTime != nil {
		resp.EndedAt = timestamppb.New(*state.EndTime)
	}

	return resp, nil
}

// Stream streams the output logs of a job.
func (s *WorkerServer) Stream(req *pb.StreamRequest, stream pb.JobWorker_StreamServer) (err error) {
	defer s.recoverToError(&err)

	slog.Info("Stream request", "id", req.JobId)

	job, err := s.authorize(stream.Context(), req.GetJobId())
	if err != nil {
		slog.Warn("Authorization failed or job not found", "id", req.JobId, "error", err)
		return err
	}

	logStream := job.Stream(stream.Context())
	defer logStream.Close()

	for {
		entry, err := logStream.Next()
		if err != nil {
			if err == io.EOF {
				slog.Info("Stream completed", "id", req.JobId)
				return nil
			}
			return err
		}
		if err := stream.Send(&pb.LogChunk{
			Data:      entry.Data,
			Source:    s.toProtoLogSource(entry.Source),
			Timestamp: timestamppb.New(entry.Timestamp),
		}); err != nil {
			slog.Warn("Failed to send log chunk", "id", req.JobId, "error", err)
			return err
		}
	}
}

// authorize checks if the job exists and if the caller is the owner.
func (s *WorkerServer) authorize(ctx context.Context, jobID string) (*worker.Job, error) {
	cn, err := s.extractClientCN(ctx)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	job, exists := s.jobs[worker.JobID(jobID)]
	s.mu.Unlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "job %s not found", jobID)
	}

	if job.Specs.OwnerCommonName != cn {
		return nil, status.Errorf(codes.NotFound, "job %s not found", jobID)
	}

	return job, nil
}

// extractClientCN extracts the Common Name from the mTLS client certificate.
func (s *WorkerServer) extractClientCN(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no peer found")
	}

	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "unexpected auth info")
	}

	if len(tlsInfo.State.VerifiedChains) == 0 || len(tlsInfo.State.VerifiedChains[0]) == 0 {
		return "", status.Error(codes.Unauthenticated, "no verified client certificate found")
	}

	return tlsInfo.State.VerifiedChains[0][0].Subject.CommonName, nil
}

// recoverToError captures panics and sets the return error to Internal.
func (s *WorkerServer) recoverToError(err *error) {
	if r := recover(); r != nil {
		slog.Error("Panic recovered", "panic", r, "stack", string(debug.Stack()))
		*err = status.Error(codes.Internal, "internal server error")
	}
}

func (s *WorkerServer) toProtoStatus(status worker.JobStatus) pb.StatusResponse_JobStatus {
	switch status {
	case worker.StatusRunning:
		return pb.StatusResponse_RUNNING
	case worker.StatusStopped:
		return pb.StatusResponse_STOPPED
	case worker.StatusCompleted:
		return pb.StatusResponse_COMPLETED
	case worker.StatusFailed:
		return pb.StatusResponse_FAILED
	default:
		return pb.StatusResponse_STATUS_UNSPECIFIED
	}
}

func (s *WorkerServer) toProtoLogSource(src worker.OutputSource) pb.LogChunk_Source {
	switch src {
	case worker.OutputSourceStdout:
		return pb.LogChunk_STDOUT
	case worker.OutputSourceStderr:
		return pb.LogChunk_STDERR
	default:
		return pb.LogChunk_SOURCE_UNSPECIFIED
	}
}

func (s *WorkerServer) toProtoLimits(l worker.ResourceLimits) *pb.ResourceLimits {
	return &pb.ResourceLimits{
		MemoryBytes:  uint64(l.MemoryLimitBytes),
		CpuPercent:   l.CPULimitPercent,
		DiskReadBps:  l.IOLimit.ReadBPS,
		DiskWriteBps: l.IOLimit.WriteBPS,
	}
}
