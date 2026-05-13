//go:build linux

package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	pb "github.com/manil/job-worker/pkg/proto"
	"github.com/manil/job-worker/pkg/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func TestNewServer(t *testing.T) {
	store := worker.NewJobStore()
	cgroupRoot := "/sys/fs/cgroup/jobworker-test"

	srv := NewServer(store, cgroupRoot)

	if srv == nil {
		t.Fatal("expected non-nil server")
	}
	if srv.store != store {
		t.Error("expected store to be set")
	}
	if srv.cgroupRoot != cgroupRoot {
		t.Errorf("expected cgroupRoot=%s, got %s", cgroupRoot, srv.cgroupRoot)
	}
}

func TestStartJob_NoIdentity(t *testing.T) {
	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	ctx := context.Background() // No identity
	req := &pb.StartJobRequest{Command: "/bin/echo"}

	_, err := srv.StartJob(ctx, req)
	if err == nil {
		t.Error("expected error when no identity in context")
	}
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("expected Unauthenticated, got %v", status.Code(err))
	}
}

func TestStopJob_NoIdentity(t *testing.T) {
	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	ctx := context.Background()
	req := &pb.StopJobRequest{JobId: "some-id"}

	_, err := srv.StopJob(ctx, req)
	if err == nil {
		t.Error("expected error when no identity")
	}
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("expected Unauthenticated, got %v", status.Code(err))
	}
}

func TestStopJob_NotFound(t *testing.T) {
	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	identity := &ClientIdentity{UserID: "client1", IsAdmin: false}
	ctx := ContextWithIdentity(context.Background(), identity)

	req := &pb.StopJobRequest{JobId: "nonexistent"}
	_, err := srv.StopJob(ctx, req)
	if err == nil {
		t.Error("expected error for nonexistent job")
	}
	// Server returns raw domain error; interceptor maps to gRPC status
	if !errors.Is(err, worker.ErrJobNotFound) {
		t.Errorf("expected ErrJobNotFound, got %v", err)
	}
}

func TestGetJobStatus_NoIdentity(t *testing.T) {
	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	ctx := context.Background()
	req := &pb.GetJobStatusRequest{JobId: "some-id"}

	_, err := srv.GetJobStatus(ctx, req)
	if err == nil {
		t.Error("expected error when no identity")
	}
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("expected Unauthenticated, got %v", status.Code(err))
	}
}

func TestGetJobStatus_NotFound(t *testing.T) {
	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	identity := &ClientIdentity{UserID: "client1", IsAdmin: false}
	ctx := ContextWithIdentity(context.Background(), identity)

	_, err := srv.GetJobStatus(ctx, &pb.GetJobStatusRequest{JobId: "nonexistent"})
	// Server returns raw domain error; interceptor maps to gRPC status
	if !errors.Is(err, worker.ErrJobNotFound) {
		t.Errorf("expected ErrJobNotFound, got %v", err)
	}
}

// mockStreamOutputServer implements grpc.ServerStreamingServer for testing.
type mockStreamOutputServer struct {
	grpc.ServerStream
	ctx     context.Context
	sent    [][]byte
	sendErr error
}

func (m *mockStreamOutputServer) Context() context.Context {
	return m.ctx
}

func (m *mockStreamOutputServer) Send(resp *pb.StreamOutputResponse) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	m.sent = append(m.sent, resp.Data)
	return nil
}

func TestStreamOutput_NoIdentity(t *testing.T) {
	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	stream := &mockStreamOutputServer{ctx: context.Background()}
	err := srv.StreamOutput(&pb.StreamOutputRequest{JobId: "some-id"}, stream)
	if err == nil {
		t.Error("expected error when no identity")
	}
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("expected Unauthenticated, got %v", status.Code(err))
	}
}

func TestStreamOutput_NotFound(t *testing.T) {
	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	identity := &ClientIdentity{UserID: "client1", IsAdmin: false}
	ctx := ContextWithIdentity(context.Background(), identity)

	stream := &mockStreamOutputServer{ctx: ctx}
	err := srv.StreamOutput(&pb.StreamOutputRequest{JobId: "nonexistent"}, stream)
	// Server returns raw domain error; interceptor maps to gRPC status
	if !errors.Is(err, worker.ErrJobNotFound) {
		t.Errorf("expected ErrJobNotFound, got %v", err)
	}
}

func TestNewTLSConfig_Success(t *testing.T) {
	certsDir := "../../certs"
	certFile := filepath.Join(certsDir, "server.pem")
	keyFile := filepath.Join(certsDir, "server-key.pem")
	caFile := filepath.Join(certsDir, "ca.pem")

	// Check if certs exist
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		t.Skip("certificates not found, skipping TLS test")
	}

	tlsConfig, err := NewTLSConfig(certFile, keyFile, caFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tlsConfig == nil {
		t.Fatal("expected non-nil TLS config")
	}
	if tlsConfig.MinVersion != tls.VersionTLS13 {
		t.Error("expected TLS 1.3 minimum")
	}
	if tlsConfig.ClientAuth != tls.RequireAndVerifyClientCert {
		t.Error("expected RequireAndVerifyClientCert")
	}
	if len(tlsConfig.Certificates) != 1 {
		t.Error("expected 1 server certificate")
	}
	if tlsConfig.ClientCAs == nil {
		t.Error("expected non-nil ClientCAs pool")
	}
}

func TestNewTLSConfig_InvalidCert(t *testing.T) {
	_, err := NewTLSConfig("nonexistent.pem", "nonexistent-key.pem", "ca.pem")
	if err == nil {
		t.Error("expected error for missing certificate")
	}
}

func loadClientTLSConfig(t *testing.T, certsDir string) *tls.Config {
	t.Helper()

	cert, err := tls.LoadX509KeyPair(
		filepath.Join(certsDir, "client1.pem"),
		filepath.Join(certsDir, "client1-key.pem"),
	)
	if err != nil {
		t.Fatalf("failed to load client cert: %v", err)
	}

	caPEM, err := os.ReadFile(filepath.Join(certsDir, "ca.pem"))
	if err != nil {
		t.Fatalf("failed to read CA: %v", err)
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caPEM)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		MinVersion:   tls.VersionTLS13,
	}
}

func TestServer_ListenAndServe(t *testing.T) {
	certsDir := "../../certs"
	if _, err := os.Stat(filepath.Join(certsDir, "server.pem")); os.IsNotExist(err) {
		t.Skip("certificates not found, skipping server test")
	}

	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	tlsConfig, err := NewTLSConfig(
		filepath.Join(certsDir, "server.pem"),
		filepath.Join(certsDir, "server-key.pem"),
		filepath.Join(certsDir, "ca.pem"),
	)
	if err != nil {
		t.Fatalf("failed to create TLS config: %v", err)
	}

	// Find available port
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	addr := lis.Addr().String()
	lis.Close()

	// Start server in background
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe(addr, tlsConfig)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Verify server is running by connecting
	clientTLS := loadClientTLSConfig(t, certsDir)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(clientTLS)))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	conn.Close()

	// Graceful shutdown
	srv.GracefulStop()

	// Check no error from serve
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("unexpected serve error: %v", err)
		}
	case <-time.After(time.Second):
		t.Error("server did not shut down")
	}
}

func TestMapJobStatus(t *testing.T) {
	tests := []struct {
		name   string
		input  worker.JobStatus
		expect pb.JobStatus
	}{
		{
			name:   "running",
			input:  worker.JobStatusRunning,
			expect: pb.JobStatus_JOB_STATUS_RUNNING,
		},
		{
			name:   "completed",
			input:  worker.JobStatusCompleted,
			expect: pb.JobStatus_JOB_STATUS_COMPLETED,
		},
		{
			name:   "failed",
			input:  worker.JobStatusFailed,
			expect: pb.JobStatus_JOB_STATUS_FAILED,
		},
		{
			name:   "stopped",
			input:  worker.JobStatusStopped,
			expect: pb.JobStatus_JOB_STATUS_STOPPED,
		},
		{
			name:   "unknown",
			input:  worker.JobStatus(999),
			expect: pb.JobStatus_JOB_STATUS_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapJobStatus(tt.input)
			if got != tt.expect {
				t.Errorf("mapJobStatus(%v) = %v, want %v", tt.input, got, tt.expect)
			}
		})
	}
}

func TestStopJob_PermissionDenied(t *testing.T) {
	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	// Add a job owned by a different user
	otherJob := &worker.Job{
		ID:    "other-job",
		Owner: "other-user",
	}
	store.Add(otherJob)

	// Try to stop it as client1 (non-admin)
	identity := &ClientIdentity{UserID: "client1", IsAdmin: false}
	ctx := ContextWithIdentity(context.Background(), identity)

	req := &pb.StopJobRequest{JobId: "other-job"}
	_, err := srv.StopJob(ctx, req)
	if err == nil {
		t.Error("expected error when stopping another user's job")
	}
	// ErrJobNotFound is returned to hide existence of job from unauthorized users
	if !errors.Is(err, worker.ErrJobNotFound) {
		t.Errorf("expected ErrJobNotFound (hiding job existence), got %v", err)
	}
}

func TestGetJobStatus_AdminCanAccessOthersJob(t *testing.T) {
	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	// Add a job owned by a different user
	otherJob := &worker.Job{
		ID:    "other-job",
		Owner: "other-user",
	}
	store.Add(otherJob)

	// Admin should be able to see/access another user's job
	identity := &ClientIdentity{UserID: "admin", IsAdmin: true}
	ctx := ContextWithIdentity(context.Background(), identity)

	req := &pb.GetJobStatusRequest{JobId: "other-job"}
	resp, err := srv.GetJobStatus(ctx, req)
	if err != nil {
		t.Errorf("admin should be able to access other user's job, got error: %v", err)
	}
	if resp != nil && resp.Owner != "other-user" {
		t.Errorf("expected owner=other-user, got %s", resp.Owner)
	}
}

func TestGetJobStatus_PermissionDenied(t *testing.T) {
	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	// Add a job owned by a different user
	otherJob := &worker.Job{
		ID:    "other-job",
		Owner: "other-user",
	}
	store.Add(otherJob)

	// Try to get status as client1 (non-admin)
	identity := &ClientIdentity{UserID: "client1", IsAdmin: false}
	ctx := ContextWithIdentity(context.Background(), identity)

	req := &pb.GetJobStatusRequest{JobId: "other-job"}
	_, err := srv.GetJobStatus(ctx, req)
	if err == nil {
		t.Error("expected error when getting status of another user's job")
	}
	// ErrJobNotFound is returned to hide existence of job from unauthorized users
	if !errors.Is(err, worker.ErrJobNotFound) {
		t.Errorf("expected ErrJobNotFound, got %v", err)
	}
}

func TestStreamOutput_PermissionDenied(t *testing.T) {
	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	// Add a job owned by a different user
	otherJob := &worker.Job{
		ID:    "other-job",
		Owner: "other-user",
	}
	store.Add(otherJob)

	// Try to stream as client1 (non-admin)
	identity := &ClientIdentity{UserID: "client1", IsAdmin: false}
	ctx := ContextWithIdentity(context.Background(), identity)

	stream := &mockStreamOutputServer{ctx: ctx}
	err := srv.StreamOutput(&pb.StreamOutputRequest{JobId: "other-job"}, stream)
	// ErrJobNotFound is returned to hide existence of job from unauthorized users
	if !errors.Is(err, worker.ErrJobNotFound) {
		t.Errorf("expected ErrJobNotFound, got %v", err)
	}
}

func TestGracefulStop_NotStarted(t *testing.T) {
	store := worker.NewJobStore()
	srv := NewServer(store, t.TempDir())

	// Should not panic when stopping server that was never started
	srv.GracefulStop()
}

// Integration test requires root and cgroups - run with: sudo go test -v -run TestIntegration
func TestIntegration_FullJobLifecycle(t *testing.T) {
	// Check if running as root (required for cgroups)
	if os.Geteuid() != 0 {
		t.Skip("skipping integration test: requires root for cgroups")
	}

	certsDir := "../../certs"
	if _, err := os.Stat(filepath.Join(certsDir, "server.pem")); os.IsNotExist(err) {
		t.Skip("certificates not found")
	}

	// Check cgroups v2 is available
	if _, err := os.Stat("/sys/fs/cgroup/cgroup.controllers"); os.IsNotExist(err) {
		t.Skip("cgroups v2 not available")
	}

	cgroupRoot := "/sys/fs/cgroup/jobworker-test"
	if err := os.MkdirAll(cgroupRoot, 0755); err != nil {
		t.Skipf("cannot create cgroup root (read-only filesystem?): %v", err)
	}
	defer os.RemoveAll(cgroupRoot)

	store := worker.NewJobStore()
	srv := NewServer(store, cgroupRoot)

	tlsConfig, err := NewTLSConfig(
		filepath.Join(certsDir, "server.pem"),
		filepath.Join(certsDir, "server-key.pem"),
		filepath.Join(certsDir, "ca.pem"),
	)
	if err != nil {
		t.Fatalf("failed to create TLS config: %v", err)
	}

	// Start server
	lis, _ := net.Listen("tcp", "localhost:0")
	addr := lis.Addr().String()
	lis.Close()

	go srv.ListenAndServe(addr, tlsConfig)
	defer srv.GracefulStop()
	time.Sleep(100 * time.Millisecond)

	// Create client
	clientTLS := loadClientTLSConfig(t, certsDir)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(clientTLS)))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer conn.Close()

	client := pb.NewWorkerServiceClient(conn)
	ctx := context.Background()

	// Start job
	startResp, err := client.StartJob(ctx, &pb.StartJobRequest{
		Command: "/bin/echo",
		Args:    []string{"integration test"},
	})
	if err != nil {
		t.Fatalf("StartJob failed: %v", err)
	}
	t.Logf("Started job: %s", startResp.JobId)

	// Get status
	statusResp, err := client.GetJobStatus(ctx, &pb.GetJobStatusRequest{JobId: startResp.JobId})
	if err != nil {
		t.Fatalf("GetJobStatus failed: %v", err)
	}
	t.Logf("Job status: %v, owner: %s", statusResp.Status, statusResp.Owner)

	// Verify owner is extracted from certificate CN
	if statusResp.Owner != "client1" {
		t.Errorf("expected owner=client1, got %s", statusResp.Owner)
	}

	// Wait for job to complete
	time.Sleep(200 * time.Millisecond)

	// Stream output
	stream, err := client.StreamOutput(ctx, &pb.StreamOutputRequest{JobId: startResp.JobId})
	if err != nil {
		t.Fatalf("StreamOutput failed: %v", err)
	}

	var output []byte
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv failed: %v", err)
		}
		output = append(output, resp.Data...)
	}

	expected := "integration test\n"
	if string(output) != expected {
		t.Errorf("expected output %q, got %q", expected, string(output))
	}

	// Verify final status
	finalStatus, err := client.GetJobStatus(ctx, &pb.GetJobStatusRequest{JobId: startResp.JobId})
	if err != nil {
		t.Fatalf("GetJobStatus failed: %v", err)
	}
	if finalStatus.Status != pb.JobStatus_JOB_STATUS_COMPLETED {
		t.Errorf("expected COMPLETED, got %v", finalStatus.Status)
	}
}
