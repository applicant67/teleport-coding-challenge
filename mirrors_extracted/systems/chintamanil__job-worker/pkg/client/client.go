// Package client provides a high-level client for the Job Worker service.
// It abstracts gRPC and protobuf details, providing a clean domain-focused API.
package client

import (
	"context"
	"fmt"
	"io"

	pb "github.com/manil/job-worker/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client wraps the gRPC connection and provides high-level methods
// for interacting with the Job Worker service.
type Client struct {
	worker pb.WorkerServiceClient
	conn   *grpc.ClientConn
}

// JobStatus represents the status of a job with all relevant details.
type JobStatus struct {
	ID           string
	Owner        string
	Status       string
	ExitCode     int32
	ErrorMessage string
	PID          int32
	Signal       int32 // Populated if stopped by signal (ExitCode > 128)
}

// StartOptions configures job resource limits.
type StartOptions struct {
	CPULimit    float64
	MemoryLimit int64
	IOWeight    int32
}

// StartOption is a functional option for configuring job start.
type StartOption func(*StartOptions)

// WithCPU sets the CPU limit as a fraction of one core.
func WithCPU(limit float64) StartOption {
	return func(o *StartOptions) {
		o.CPULimit = limit
	}
}

// WithMemory sets the memory limit in bytes.
func WithMemory(limit int64) StartOption {
	return func(o *StartOptions) {
		o.MemoryLimit = limit
	}
}

// WithIOWeight sets the I/O weight (1-10000).
func WithIOWeight(weight int32) StartOption {
	return func(o *StartOptions) {
		o.IOWeight = weight
	}
}

// New creates a new client connected to the Job Worker service with mTLS.
func New(addr, certFile, keyFile, caFile string) (*Client, error) {
	tlsConfig, err := NewTLSConfig(certFile, keyFile, caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to configure TLS: %w", err)
	}

	creds := credentials.NewTLS(tlsConfig)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	return &Client{
		worker: pb.NewWorkerServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close closes the underlying gRPC connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Start starts a new job with the given command and arguments.
// Returns the job ID on success.
func (c *Client) Start(ctx context.Context, command string, args []string, opts ...StartOption) (string, error) {
	options := &StartOptions{}
	for _, opt := range opts {
		opt(options)
	}

	req := &pb.StartJobRequest{
		Command:          command,
		Args:             args,
		CpuLimit:         options.CPULimit,
		MemoryLimitBytes: options.MemoryLimit,
		IoWeight:         options.IOWeight,
	}

	resp, err := c.worker.StartJob(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to start job: %w", err)
	}

	return resp.JobId, nil
}

// Stop stops a running job.
func (c *Client) Stop(ctx context.Context, jobID string) error {
	req := &pb.StopJobRequest{
		JobId: jobID,
	}

	_, err := c.worker.StopJob(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to stop job: %w", err)
	}

	return nil
}

// Status returns the current status of a job.
func (c *Client) Status(ctx context.Context, jobID string) (*JobStatus, error) {
	req := &pb.GetJobStatusRequest{
		JobId: jobID,
	}

	resp, err := c.worker.GetJobStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get job status: %w", err)
	}

	status := &JobStatus{
		ID:           resp.JobId,
		Owner:        resp.Owner,
		Status:       formatStatus(resp.Status),
		ExitCode:     resp.ExitCode,
		ErrorMessage: resp.ErrorMessage,
		PID:          resp.Pid,
	}

	// Extract signal if job was stopped
	if resp.Status == pb.JobStatus_JOB_STATUS_STOPPED && resp.ExitCode > 128 {
		status.Signal = resp.ExitCode - 128
	}

	return status, nil
}

// Logs streams the job output to the provided writer.
// It blocks until the job completes or the context is cancelled.
func (c *Client) Logs(ctx context.Context, jobID string, w io.Writer) error {
	req := &pb.StreamOutputRequest{
		JobId: jobID,
	}

	stream, err := c.worker.StreamOutput(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to stream output: %w", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			// Context cancelled - exit cleanly
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("stream error: %w", err)
		}

		if _, err := w.Write(resp.Data); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}
}

// formatStatus converts a protobuf JobStatus to a human-readable string.
func formatStatus(status pb.JobStatus) string {
	switch status {
	case pb.JobStatus_JOB_STATUS_RUNNING:
		return "RUNNING"
	case pb.JobStatus_JOB_STATUS_COMPLETED:
		return "COMPLETED"
	case pb.JobStatus_JOB_STATUS_FAILED:
		return "FAILED"
	case pb.JobStatus_JOB_STATUS_STOPPED:
		return "STOPPED"
	default:
		return "UNKNOWN"
	}
}
