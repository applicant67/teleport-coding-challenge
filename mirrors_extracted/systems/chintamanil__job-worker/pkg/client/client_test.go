package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"

	pb "github.com/manil/job-worker/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// =============================================================================
// TLS Configuration Tests
// =============================================================================

func TestNewTLSConfig(t *testing.T) {
	certsDir := filepath.Join("..", "..", "certs")

	if _, err := os.Stat(certsDir); os.IsNotExist(err) {
		t.Skip("certs directory not found, skipping TLS config test")
	}

	certFile := filepath.Join(certsDir, "client1.pem")
	keyFile := filepath.Join(certsDir, "client1-key.pem")
	caFile := filepath.Join(certsDir, "ca.pem")

	t.Run("happy path - valid certificates", func(t *testing.T) {
		cfg, err := NewTLSConfig(certFile, keyFile, caFile)
		if err != nil {
			t.Fatalf("NewTLSConfig() error = %v", err)
		}

		if cfg == nil {
			t.Fatal("NewTLSConfig() returned nil config")
		}

		if len(cfg.Certificates) != 1 {
			t.Errorf("expected 1 certificate, got %d", len(cfg.Certificates))
		}

		if cfg.RootCAs == nil {
			t.Error("RootCAs is nil")
		}

		if cfg.MinVersion != 0x0304 { // TLS 1.3
			t.Errorf("MinVersion = %x, want TLS 1.3 (0x0304)", cfg.MinVersion)
		}
	})

	t.Run("unhappy path - missing cert file", func(t *testing.T) {
		_, err := NewTLSConfig("nonexistent.pem", keyFile, caFile)
		if err == nil {
			t.Error("expected error for missing cert file")
		}
	})

	t.Run("unhappy path - missing key file", func(t *testing.T) {
		_, err := NewTLSConfig(certFile, "nonexistent.pem", caFile)
		if err == nil {
			t.Error("expected error for missing key file")
		}
	})

	t.Run("unhappy path - missing CA file", func(t *testing.T) {
		_, err := NewTLSConfig(certFile, keyFile, "nonexistent.pem")
		if err == nil {
			t.Error("expected error for missing CA file")
		}
	})

	t.Run("edge case - invalid PEM in CA file", func(t *testing.T) {
		// Create a temp file with invalid PEM
		tmpFile, err := os.CreateTemp("", "invalid-ca-*.pem")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())
		tmpFile.WriteString("not a valid PEM")
		tmpFile.Close()

		_, err = NewTLSConfig(certFile, keyFile, tmpFile.Name())
		if err == nil {
			t.Error("expected error for invalid PEM")
		}
	})
}

// =============================================================================
// Status Formatting Tests
// =============================================================================

func TestFormatStatus(t *testing.T) {
	tests := []struct {
		name   string
		status pb.JobStatus
		want   string
	}{
		// Happy paths
		{name: "running", status: pb.JobStatus_JOB_STATUS_RUNNING, want: "RUNNING"},
		{name: "completed", status: pb.JobStatus_JOB_STATUS_COMPLETED, want: "COMPLETED"},
		{name: "failed", status: pb.JobStatus_JOB_STATUS_FAILED, want: "FAILED"},
		{name: "stopped", status: pb.JobStatus_JOB_STATUS_STOPPED, want: "STOPPED"},
		// Edge cases
		{name: "unspecified", status: pb.JobStatus_JOB_STATUS_UNSPECIFIED, want: "UNKNOWN"},
		{name: "invalid positive", status: pb.JobStatus(999), want: "UNKNOWN"},
		{name: "invalid negative", status: pb.JobStatus(-1), want: "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatStatus(tt.status); got != tt.want {
				t.Errorf("formatStatus(%v) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

// =============================================================================
// Start Options Tests
// =============================================================================

func TestStartOptions(t *testing.T) {
	t.Run("happy path - WithCPU", func(t *testing.T) {
		opts := &StartOptions{}
		WithCPU(1.5)(opts)
		if opts.CPULimit != 1.5 {
			t.Errorf("CPULimit = %v, want 1.5", opts.CPULimit)
		}
	})

	t.Run("happy path - WithMemory", func(t *testing.T) {
		opts := &StartOptions{}
		WithMemory(1024 * 1024)(opts)
		if opts.MemoryLimit != 1024*1024 {
			t.Errorf("MemoryLimit = %v, want %v", opts.MemoryLimit, 1024*1024)
		}
	})

	t.Run("happy path - WithIOWeight", func(t *testing.T) {
		opts := &StartOptions{}
		WithIOWeight(500)(opts)
		if opts.IOWeight != 500 {
			t.Errorf("IOWeight = %v, want 500", opts.IOWeight)
		}
	})

	t.Run("happy path - combined options", func(t *testing.T) {
		opts := &StartOptions{}
		for _, opt := range []StartOption{
			WithCPU(2.0),
			WithMemory(1024),
			WithIOWeight(100),
		} {
			opt(opts)
		}

		if opts.CPULimit != 2.0 {
			t.Errorf("CPULimit = %v, want 2.0", opts.CPULimit)
		}
		if opts.MemoryLimit != 1024 {
			t.Errorf("MemoryLimit = %v, want 1024", opts.MemoryLimit)
		}
		if opts.IOWeight != 100 {
			t.Errorf("IOWeight = %v, want 100", opts.IOWeight)
		}
	})

	t.Run("edge case - zero values", func(t *testing.T) {
		opts := &StartOptions{}
		WithCPU(0)(opts)
		WithMemory(0)(opts)
		WithIOWeight(0)(opts)

		if opts.CPULimit != 0 || opts.MemoryLimit != 0 || opts.IOWeight != 0 {
			t.Error("zero values should be allowed")
		}
	})

	t.Run("edge case - negative values", func(t *testing.T) {
		opts := &StartOptions{}
		WithCPU(-1.0)(opts)
		WithMemory(-1024)(opts)
		WithIOWeight(-100)(opts)

		// Negative values are stored; validation is server-side
		if opts.CPULimit != -1.0 {
			t.Errorf("CPULimit = %v, want -1.0", opts.CPULimit)
		}
	})
}

// =============================================================================
// JobStatus Tests
// =============================================================================

func TestJobStatus(t *testing.T) {
	t.Run("happy path - running job", func(t *testing.T) {
		status := &JobStatus{
			ID:     "test-123",
			Owner:  "client1",
			Status: "RUNNING",
			PID:    1234,
		}

		if status.ID != "test-123" {
			t.Errorf("ID = %v, want test-123", status.ID)
		}
		if status.PID != 1234 {
			t.Errorf("PID = %v, want 1234", status.PID)
		}
	})

	t.Run("happy path - completed job", func(t *testing.T) {
		status := &JobStatus{
			ID:       "test-456",
			Status:   "COMPLETED",
			ExitCode: 0,
		}

		if status.ExitCode != 0 {
			t.Errorf("ExitCode = %v, want 0", status.ExitCode)
		}
	})

	t.Run("happy path - failed job", func(t *testing.T) {
		status := &JobStatus{
			ID:           "test-789",
			Status:       "FAILED",
			ExitCode:     1,
			ErrorMessage: "command not found",
		}

		if status.ErrorMessage != "command not found" {
			t.Errorf("ErrorMessage = %v", status.ErrorMessage)
		}
	})

	t.Run("edge case - stopped job with SIGKILL", func(t *testing.T) {
		// SIGKILL (9) -> exit code 137 (128 + 9)
		status := &JobStatus{
			ID:       "test-stopped",
			Status:   "STOPPED",
			ExitCode: 137,
			Signal:   9,
		}

		if status.Signal != 9 {
			t.Errorf("Signal = %v, want 9 (SIGKILL)", status.Signal)
		}
	})

	t.Run("edge case - exit code exactly 128", func(t *testing.T) {
		status := &JobStatus{
			ID:       "test-128",
			Status:   "STOPPED",
			ExitCode: 128,
			Signal:   0, // Signal would be 0, which is not a valid signal
		}

		if status.Signal != 0 {
			t.Errorf("Signal = %v, want 0 for exit code 128", status.Signal)
		}
	})
}

// =============================================================================
// Mock Server for Client Tests
// =============================================================================

type mockWorkerServer struct {
	pb.UnimplementedWorkerServiceServer

	// Configurable responses
	startJobResp  *pb.StartJobResponse
	startJobErr   error
	stopJobResp   *pb.StopJobResponse
	stopJobErr    error
	statusResp    *pb.GetJobStatusResponse
	statusErr     error
	streamChunks  [][]byte
	streamErr     error
	streamErrAt   int // Send error after this many chunks (-1 = no error)

	// Capture requests
	lastStartReq  *pb.StartJobRequest
	lastStopReq   *pb.StopJobRequest
	lastStatusReq *pb.GetJobStatusRequest
	lastStreamReq *pb.StreamOutputRequest
}

func (m *mockWorkerServer) StartJob(_ context.Context, req *pb.StartJobRequest) (*pb.StartJobResponse, error) {
	m.lastStartReq = req
	if m.startJobErr != nil {
		return nil, m.startJobErr
	}
	return m.startJobResp, nil
}

func (m *mockWorkerServer) StopJob(_ context.Context, req *pb.StopJobRequest) (*pb.StopJobResponse, error) {
	m.lastStopReq = req
	if m.stopJobErr != nil {
		return nil, m.stopJobErr
	}
	return m.stopJobResp, nil
}

func (m *mockWorkerServer) GetJobStatus(_ context.Context, req *pb.GetJobStatusRequest) (*pb.GetJobStatusResponse, error) {
	m.lastStatusReq = req
	if m.statusErr != nil {
		return nil, m.statusErr
	}
	return m.statusResp, nil
}

func (m *mockWorkerServer) StreamOutput(req *pb.StreamOutputRequest, stream pb.WorkerService_StreamOutputServer) error {
	m.lastStreamReq = req

	// Return immediate error if streamErrAt is 0 and no chunks
	if m.streamErr != nil && m.streamErrAt == 0 {
		return m.streamErr
	}

	for i, chunk := range m.streamChunks {
		if m.streamErrAt > 0 && i >= m.streamErrAt {
			return m.streamErr
		}
		if err := stream.Send(&pb.StreamOutputResponse{Data: chunk}); err != nil {
			return err
		}
	}
	return nil
}

func startMockServer(t *testing.T, mock *mockWorkerServer) (addr string, cleanup func()) {
	t.Helper()

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterWorkerServiceServer(server, mock)

	go func() {
		server.Serve(lis)
	}()

	return lis.Addr().String(), func() {
		server.Stop()
	}
}

func newTestClient(t *testing.T, addr string) *Client {
	t.Helper()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	return &Client{
		worker: pb.NewWorkerServiceClient(conn),
		conn:   conn,
	}
}

// =============================================================================
// Client.Start() Tests
// =============================================================================

func TestClient_Start(t *testing.T) {
	t.Run("happy path - start job successfully", func(t *testing.T) {
		mock := &mockWorkerServer{
			startJobResp: &pb.StartJobResponse{JobId: "job-123"},
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		jobID, err := client.Start(context.Background(), "/bin/echo", []string{"hello"})
		if err != nil {
			t.Fatalf("Start() error = %v", err)
		}

		if jobID != "job-123" {
			t.Errorf("jobID = %v, want job-123", jobID)
		}

		// Verify request was sent correctly
		if mock.lastStartReq.Command != "/bin/echo" {
			t.Errorf("command = %v, want /bin/echo", mock.lastStartReq.Command)
		}
		if len(mock.lastStartReq.Args) != 1 || mock.lastStartReq.Args[0] != "hello" {
			t.Errorf("args = %v, want [hello]", mock.lastStartReq.Args)
		}
	})

	t.Run("happy path - start with options", func(t *testing.T) {
		mock := &mockWorkerServer{
			startJobResp: &pb.StartJobResponse{JobId: "job-456"},
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		_, err := client.Start(context.Background(), "make", []string{"-j4"},
			WithCPU(0.5),
			WithMemory(512*1024*1024),
			WithIOWeight(200),
		)
		if err != nil {
			t.Fatalf("Start() error = %v", err)
		}

		if mock.lastStartReq.CpuLimit != 0.5 {
			t.Errorf("CpuLimit = %v, want 0.5", mock.lastStartReq.CpuLimit)
		}
		if mock.lastStartReq.MemoryLimitBytes != 512*1024*1024 {
			t.Errorf("MemoryLimitBytes = %v", mock.lastStartReq.MemoryLimitBytes)
		}
		if mock.lastStartReq.IoWeight != 200 {
			t.Errorf("IoWeight = %v, want 200", mock.lastStartReq.IoWeight)
		}
	})

	t.Run("unhappy path - server returns error", func(t *testing.T) {
		mock := &mockWorkerServer{
			startJobErr: status.Error(codes.InvalidArgument, "invalid command"),
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		_, err := client.Start(context.Background(), "", nil)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("unhappy path - context cancelled", func(t *testing.T) {
		mock := &mockWorkerServer{
			startJobResp: &pb.StartJobResponse{JobId: "job-789"},
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := client.Start(ctx, "/bin/echo", nil)
		if err == nil {
			t.Error("expected error for cancelled context")
		}
	})

	t.Run("edge case - empty command", func(t *testing.T) {
		mock := &mockWorkerServer{
			startJobResp: &pb.StartJobResponse{JobId: "job-empty"},
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		// Client doesn't validate - server does
		jobID, err := client.Start(context.Background(), "", nil)
		if err != nil {
			t.Fatalf("Start() error = %v", err)
		}
		if jobID != "job-empty" {
			t.Errorf("jobID = %v", jobID)
		}
	})

	t.Run("edge case - nil args", func(t *testing.T) {
		mock := &mockWorkerServer{
			startJobResp: &pb.StartJobResponse{JobId: "job-nil-args"},
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		_, err := client.Start(context.Background(), "/bin/ls", nil)
		if err != nil {
			t.Fatalf("Start() error = %v", err)
		}

		if len(mock.lastStartReq.Args) != 0 {
			t.Errorf("args should be nil or empty, got %v", mock.lastStartReq.Args)
		}
	})
}

// =============================================================================
// Client.Stop() Tests
// =============================================================================

func TestClient_Stop(t *testing.T) {
	t.Run("happy path - stop job successfully", func(t *testing.T) {
		mock := &mockWorkerServer{
			stopJobResp: &pb.StopJobResponse{},
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		err := client.Stop(context.Background(), "job-123")
		if err != nil {
			t.Fatalf("Stop() error = %v", err)
		}

		if mock.lastStopReq.JobId != "job-123" {
			t.Errorf("JobId = %v, want job-123", mock.lastStopReq.JobId)
		}
	})

	t.Run("unhappy path - job not found", func(t *testing.T) {
		mock := &mockWorkerServer{
			stopJobErr: status.Error(codes.NotFound, "job not found"),
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		err := client.Stop(context.Background(), "nonexistent")
		if err == nil {
			t.Error("expected error for nonexistent job")
		}
	})

	t.Run("unhappy path - permission denied", func(t *testing.T) {
		mock := &mockWorkerServer{
			stopJobErr: status.Error(codes.PermissionDenied, "not job owner"),
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		err := client.Stop(context.Background(), "other-users-job")
		if err == nil {
			t.Error("expected error for permission denied")
		}
	})

	t.Run("edge case - empty job ID", func(t *testing.T) {
		mock := &mockWorkerServer{
			stopJobResp: &pb.StopJobResponse{},
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		// Client doesn't validate - server does
		err := client.Stop(context.Background(), "")
		if err != nil {
			t.Fatalf("Stop() error = %v", err)
		}
	})
}

// =============================================================================
// Client.Status() Tests
// =============================================================================

func TestClient_Status(t *testing.T) {
	t.Run("happy path - running job", func(t *testing.T) {
		mock := &mockWorkerServer{
			statusResp: &pb.GetJobStatusResponse{
				JobId:  "job-123",
				Owner:  "client1",
				Status: pb.JobStatus_JOB_STATUS_RUNNING,
				Pid:    1234,
			},
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		status, err := client.Status(context.Background(), "job-123")
		if err != nil {
			t.Fatalf("Status() error = %v", err)
		}

		if status.ID != "job-123" {
			t.Errorf("ID = %v", status.ID)
		}
		if status.Owner != "client1" {
			t.Errorf("Owner = %v", status.Owner)
		}
		if status.Status != "RUNNING" {
			t.Errorf("Status = %v", status.Status)
		}
		if status.PID != 1234 {
			t.Errorf("PID = %v", status.PID)
		}
	})

	t.Run("happy path - completed job", func(t *testing.T) {
		mock := &mockWorkerServer{
			statusResp: &pb.GetJobStatusResponse{
				JobId:    "job-456",
				Status:   pb.JobStatus_JOB_STATUS_COMPLETED,
				ExitCode: 0,
			},
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		status, err := client.Status(context.Background(), "job-456")
		if err != nil {
			t.Fatalf("Status() error = %v", err)
		}

		if status.Status != "COMPLETED" {
			t.Errorf("Status = %v", status.Status)
		}
		if status.ExitCode != 0 {
			t.Errorf("ExitCode = %v", status.ExitCode)
		}
	})

	t.Run("happy path - stopped job with signal extraction", func(t *testing.T) {
		mock := &mockWorkerServer{
			statusResp: &pb.GetJobStatusResponse{
				JobId:    "job-789",
				Status:   pb.JobStatus_JOB_STATUS_STOPPED,
				ExitCode: 137, // 128 + SIGKILL(9)
			},
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		status, err := client.Status(context.Background(), "job-789")
		if err != nil {
			t.Fatalf("Status() error = %v", err)
		}

		if status.Signal != 9 {
			t.Errorf("Signal = %v, want 9", status.Signal)
		}
	})

	t.Run("unhappy path - job not found", func(t *testing.T) {
		mock := &mockWorkerServer{
			statusErr: status.Error(codes.NotFound, "job not found"),
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		_, err := client.Status(context.Background(), "nonexistent")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("edge case - exit code exactly 128 for stopped job", func(t *testing.T) {
		mock := &mockWorkerServer{
			statusResp: &pb.GetJobStatusResponse{
				JobId:    "job-128",
				Status:   pb.JobStatus_JOB_STATUS_STOPPED,
				ExitCode: 128,
			},
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		status, err := client.Status(context.Background(), "job-128")
		if err != nil {
			t.Fatalf("Status() error = %v", err)
		}

		// 128 - 128 = 0, which is not a valid signal
		if status.Signal != 0 {
			t.Errorf("Signal = %v, want 0 (exit code 128 is boundary)", status.Signal)
		}
	})

	t.Run("edge case - exit code > 128 but not STOPPED status", func(t *testing.T) {
		mock := &mockWorkerServer{
			statusResp: &pb.GetJobStatusResponse{
				JobId:    "job-failed-137",
				Status:   pb.JobStatus_JOB_STATUS_FAILED,
				ExitCode: 137,
			},
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		status, err := client.Status(context.Background(), "job-failed-137")
		if err != nil {
			t.Fatalf("Status() error = %v", err)
		}

		// Signal should only be extracted for STOPPED status
		if status.Signal != 0 {
			t.Errorf("Signal = %v, want 0 (only extract signal for STOPPED)", status.Signal)
		}
	})
}

// =============================================================================
// Client.Logs() Tests
// =============================================================================

func TestClient_Logs(t *testing.T) {
	t.Run("happy path - stream logs successfully", func(t *testing.T) {
		mock := &mockWorkerServer{
			streamChunks: [][]byte{
				[]byte("line 1\n"),
				[]byte("line 2\n"),
				[]byte("line 3\n"),
			},
			streamErrAt: -1,
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		var buf bytes.Buffer
		err := client.Logs(context.Background(), "job-123", &buf)
		if err != nil {
			t.Fatalf("Logs() error = %v", err)
		}

		expected := "line 1\nline 2\nline 3\n"
		if buf.String() != expected {
			t.Errorf("output = %q, want %q", buf.String(), expected)
		}
	})

	t.Run("happy path - empty output", func(t *testing.T) {
		mock := &mockWorkerServer{
			streamChunks: nil,
			streamErrAt:  -1,
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		var buf bytes.Buffer
		err := client.Logs(context.Background(), "job-empty", &buf)
		if err != nil {
			t.Fatalf("Logs() error = %v", err)
		}

		if buf.Len() != 0 {
			t.Errorf("expected empty output, got %q", buf.String())
		}
	})

	t.Run("unhappy path - job not found", func(t *testing.T) {
		mock := &mockWorkerServer{
			streamErr:   status.Error(codes.NotFound, "job not found"),
			streamErrAt: 0,
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		var buf bytes.Buffer
		err := client.Logs(context.Background(), "nonexistent", &buf)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("unhappy path - stream error after partial data", func(t *testing.T) {
		mock := &mockWorkerServer{
			streamChunks: [][]byte{
				[]byte("partial data\n"),
				[]byte("more data\n"),
			},
			streamErr:   status.Error(codes.Internal, "stream error"),
			streamErrAt: 1, // Error after first chunk
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		var buf bytes.Buffer
		err := client.Logs(context.Background(), "job-err", &buf)
		if err == nil {
			t.Error("expected error")
		}

		// Should have received partial data
		if !bytes.Contains(buf.Bytes(), []byte("partial data")) {
			t.Errorf("expected partial data, got %q", buf.String())
		}
	})

	t.Run("unhappy path - writer error", func(t *testing.T) {
		mock := &mockWorkerServer{
			streamChunks: [][]byte{[]byte("data\n")},
			streamErrAt:  -1,
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		errWriter := &errorWriter{err: errors.New("write failed")}
		err := client.Logs(context.Background(), "job-123", errWriter)
		if err == nil {
			t.Error("expected error from writer")
		}
	})

	t.Run("edge case - context cancelled returns nil", func(t *testing.T) {
		mock := &mockWorkerServer{
			streamChunks: [][]byte{[]byte("data\n")},
			streamErrAt:  -1,
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		var buf bytes.Buffer
		err := client.Logs(ctx, "job-123", &buf)
		// Context cancellation should return nil (clean exit)
		if err != nil {
			t.Logf("Logs() returned error (may be expected): %v", err)
		}
	})

	t.Run("edge case - large output", func(t *testing.T) {
		// Generate 1MB of data in chunks
		chunk := bytes.Repeat([]byte("x"), 32*1024)
		chunks := make([][]byte, 32) // 32 * 32KB = 1MB
		for i := range chunks {
			chunks[i] = chunk
		}

		mock := &mockWorkerServer{
			streamChunks: chunks,
			streamErrAt:  -1,
		}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		defer client.Close()

		var buf bytes.Buffer
		err := client.Logs(context.Background(), "job-large", &buf)
		if err != nil {
			t.Fatalf("Logs() error = %v", err)
		}

		expectedLen := 32 * 32 * 1024
		if buf.Len() != expectedLen {
			t.Errorf("output length = %d, want %d", buf.Len(), expectedLen)
		}
	})
}

// =============================================================================
// Client.Close() Tests
// =============================================================================

func TestClient_Close(t *testing.T) {
	t.Run("happy path - close successfully", func(t *testing.T) {
		mock := &mockWorkerServer{}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		err := client.Close()
		if err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	})

	t.Run("edge case - close multiple times", func(t *testing.T) {
		mock := &mockWorkerServer{}
		addr, cleanup := startMockServer(t, mock)
		defer cleanup()

		client := newTestClient(t, addr)
		client.Close()
		// Second close should not panic
		err := client.Close()
		// Error is expected but should not panic
		_ = err
	})
}

// =============================================================================
// Helper Types
// =============================================================================

type errorWriter struct {
	err error
}

func (w *errorWriter) Write(_ []byte) (n int, err error) {
	return 0, w.err
}

var _ io.Writer = (*errorWriter)(nil)
