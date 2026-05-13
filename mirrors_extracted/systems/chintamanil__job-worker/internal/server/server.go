//go:build linux

package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	pb "github.com/manil/job-worker/pkg/proto"
	"github.com/manil/job-worker/pkg/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

// Server implements the WorkerService gRPC server.
type Server struct {
	pb.UnimplementedWorkerServiceServer
	store      *worker.JobStore
	cgroupRoot string

	mu         sync.Mutex
	grpcServer *grpc.Server
}

// NewServer creates a new gRPC server instance.
func NewServer(store *worker.JobStore, cgroupRoot string) *Server {
	return &Server{
		store:      store,
		cgroupRoot: cgroupRoot,
	}
}

// StartJob creates and starts a new job.
// Errors are mapped to gRPC status codes by the interceptor.
func (s *Server) StartJob(ctx context.Context, req *pb.StartJobRequest) (*pb.StartJobResponse, error) {
	identity, ok := IdentityFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no identity in context")
	}

	job := &worker.Job{
		ID:          uuid.New().String(),
		Owner:       identity.UserID,
		Command:     req.Command,
		Args:        req.Args,
		CPULimit:    req.CpuLimit,
		MemoryLimit: req.MemoryLimitBytes,
		IOWeight:    req.IoWeight,
	}

	if err := job.Start(s.cgroupRoot); err != nil {
		return nil, err // Interceptor maps to gRPC status
	}

	s.store.Add(job)

	return &pb.StartJobResponse{JobId: job.ID}, nil
}

// StopJob terminates a running job.
// Errors are mapped to gRPC status codes by the interceptor.
func (s *Server) StopJob(ctx context.Context, req *pb.StopJobRequest) (*pb.StopJobResponse, error) {
	identity, ok := IdentityFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no identity in context")
	}

	job, err := s.store.Get(identity.UserID, req.JobId, identity.IsAdmin)
	if err != nil {
		return nil, err // Interceptor maps to gRPC status
	}

	if err := job.Stop(); err != nil {
		return nil, err // Interceptor maps to gRPC status
	}

	return &pb.StopJobResponse{}, nil
}

// GetJobStatus returns the current status of a job.
// Errors are mapped to gRPC status codes by the interceptor.
func (s *Server) GetJobStatus(ctx context.Context, req *pb.GetJobStatusRequest) (*pb.GetJobStatusResponse, error) {
	identity, ok := IdentityFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no identity in context")
	}

	job, err := s.store.Get(identity.UserID, req.JobId, identity.IsAdmin)
	if err != nil {
		return nil, err // Interceptor maps to gRPC status
	}

	state := job.State()
	return &pb.GetJobStatusResponse{
		JobId:        job.ID,
		Status:       mapJobStatus(state.Status),
		ExitCode:     int32(state.ExitCode),
		ErrorMessage: state.Error,
		Owner:        job.Owner,
		Pid:          int32(state.Pid),
	}, nil
}

// StreamOutput streams job output to the client.
// Errors are mapped to gRPC status codes by the interceptor.
func (s *Server) StreamOutput(req *pb.StreamOutputRequest, stream pb.WorkerService_StreamOutputServer) error {
	identity, ok := IdentityFromContext(stream.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "no identity in context")
	}

	job, err := s.store.Get(identity.UserID, req.JobId, identity.IsAdmin)
	if err != nil {
		return err // Interceptor maps to gRPC status
	}

	reader := job.Output.NewReader(stream.Context())
	defer reader.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			return nil // Job finished, all output sent
		}
		if err != nil {
			return nil // Context cancelled, clean exit
		}
		if err := stream.Send(&pb.StreamOutputResponse{Data: buf[:n]}); err != nil {
			return nil // Client disconnected
		}
	}
}

// NewTLSConfig creates a TLS configuration for mTLS with TLS 1.3.
func NewTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	// Load server certificate
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}

	// Load CA certificate for client verification
	caPEM, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

// ListenAndServe starts the gRPC server with mTLS.
func (s *Server) ListenAndServe(addr string, tlsConfig *tls.Config) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.mu.Lock()
	s.grpcServer = grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.UnaryInterceptor(UnaryAuthInterceptor),
		grpc.StreamInterceptor(StreamAuthInterceptor),
	)
	pb.RegisterWorkerServiceServer(s.grpcServer, s)
	s.mu.Unlock()

	return s.grpcServer.Serve(lis)
}

// ListenAndServeWithSignals starts the gRPC server and handles OS signals for graceful shutdown.
// Blocks until SIGINT or SIGTERM is received, then gracefully stops the server.
func (s *Server) ListenAndServeWithSignals(addr string, tlsConfig *tls.Config) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.ListenAndServe(addr, tlsConfig)
	}()

	select {
	case <-ctx.Done():
		s.GracefulStop()
		return nil
	case err := <-errCh:
		return err
	}
}

// GracefulStop gracefully stops the gRPC server.
func (s *Server) GracefulStop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

// mapJobStatus converts worker.JobStatus to proto JobStatus.
func mapJobStatus(s worker.JobStatus) pb.JobStatus {
	switch s {
	case worker.JobStatusRunning:
		return pb.JobStatus_JOB_STATUS_RUNNING
	case worker.JobStatusCompleted:
		return pb.JobStatus_JOB_STATUS_COMPLETED
	case worker.JobStatusFailed:
		return pb.JobStatus_JOB_STATUS_FAILED
	case worker.JobStatusStopped:
		return pb.JobStatus_JOB_STATUS_STOPPED
	default:
		return pb.JobStatus_JOB_STATUS_UNSPECIFIED
	}
}
