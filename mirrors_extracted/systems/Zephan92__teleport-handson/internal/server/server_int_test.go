//go:build linux

package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	pb "github.com/teleport-handson/api/v1"
	"github.com/teleport-handson/internal/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// ctxWithUser creates a context with a mock mTLS peer certificate.
// This simulates a request coming from an authenticated user with the given Common Name.
func ctxWithUser(ctx context.Context, cn string) context.Context {
	cert := &x509.Certificate{Subject: pkix.Name{CommonName: cn}}
	return peer.NewContext(ctx, &peer.Peer{
		AuthInfo: credentials.TLSInfo{
			State: tls.ConnectionState{
				VerifiedChains:    [][]*x509.Certificate{{cert}},
				HandshakeComplete: true,
			},
		},
	})
}

// mockStream implements pb.JobWorker_StreamServer for testing.
type mockStream struct {
	grpc.ServerStream
	ctx      context.Context
	messages []*pb.LogChunk
	sendErr  error
}

func (m *mockStream) Context() context.Context {
	return m.ctx
}

func (m *mockStream) Send(msg *pb.LogChunk) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	m.messages = append(m.messages, msg)
	return nil
}

func TestWorkerServer_StartStop(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())
	ctx := ctxWithUser(t.Context(), "alice")

	// 1. Start a simple job
	startReq := &pb.StartRequest{
		Command:          "sleep",
		Args:             []string{"0.1"},
		WorkingDirectory: "/tmp",
	}
	startResp, err := s.Start(ctx, startReq)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if startResp.JobId == "" {
		t.Fatal("Expected JobId, got empty")
	}

	// 2. Verify Status
	statusReq := &pb.StatusRequest{JobId: startResp.JobId}
	statusResp, err := s.Status(ctx, statusReq)
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if statusResp.Owner != "alice" {
		t.Errorf("Expected owner alice, got %s", statusResp.Owner)
	}

	// 3. Stop the job
	stopReq := &pb.StopRequest{JobId: startResp.JobId}
	_, err = s.Stop(ctx, stopReq)
	// Stop might fail if the job finished quickly (sleep 0.1), which is acceptable in this test.
	// We primarily want to ensure it doesn't return Unauthenticated.
	if status.Code(err) == codes.Unauthenticated {
		t.Errorf("Stop returned Unauthenticated")
	}
}

func TestWorkerServer_Authorization(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())
	aliceCtx := ctxWithUser(t.Context(), "alice")
	bobCtx := ctxWithUser(t.Context(), "bob")

	// Alice starts a job
	resp, err := s.Start(aliceCtx, &pb.StartRequest{
		Command: "sleep", Args: []string{"10"}, WorkingDirectory: "/tmp",
	})
	if err != nil {
		t.Fatalf("Failed to start job: %v", err)
	}

	// Bob tries to stop it (Should fail)
	_, err = s.Stop(bobCtx, &pb.StopRequest{JobId: resp.JobId})
	if status.Code(err) != codes.PermissionDenied {
		t.Errorf("Expected PermissionDenied for Bob stopping job, got: %v", err)
	}

	// Bob tries to check status (Should fail)
	_, err = s.Status(bobCtx, &pb.StatusRequest{JobId: resp.JobId})
	if status.Code(err) != codes.PermissionDenied {
		t.Errorf("Expected PermissionDenied for Bob checking status, got: %v", err)
	}

	// Bob tries to stream (Should fail)
	stream := &mockStream{ctx: bobCtx}
	err = s.Stream(&pb.StreamRequest{JobId: resp.JobId}, stream)
	if status.Code(err) != codes.PermissionDenied {
		t.Errorf("Expected PermissionDenied for Bob streaming job, got: %v", err)
	}

	// Alice stops it (Should succeed)
	_, err = s.Stop(aliceCtx, &pb.StopRequest{JobId: resp.JobId})
	if err != nil {
		t.Errorf("Alice failed to stop job: %v", err)
	}
}

func TestWorkerServer_MissingPeer(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())
	ctx := context.Background() // No peer info

	_, err := s.Start(ctx, &pb.StartRequest{Command: "echo"})
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("Expected Unauthenticated for missing peer, got: %v", err)
	}
}

func TestWorkerServer_Auth_NoCert(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())

	// Context with TLS info but no certificates
	ctx := peer.NewContext(context.Background(), &peer.Peer{
		AuthInfo: credentials.TLSInfo{
			State: tls.ConnectionState{
				PeerCertificates:  []*x509.Certificate{},
				HandshakeComplete: true,
			},
		},
	})

	_, err := s.Start(ctx, &pb.StartRequest{Command: "echo"})
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("Expected Unauthenticated for missing client cert, got: %v", err)
	}
}

func TestWorkerServer_Stream_Success(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())
	ctx := ctxWithUser(t.Context(), "alice")

	// Start a job that prints to stdout and stderr
	resp, err := s.Start(ctx, &pb.StartRequest{
		Command:          "sh",
		Args:             []string{"-c", "echo stdout-msg; echo stderr-msg >&2"},
		WorkingDirectory: "/tmp",
	})
	if err != nil {
		t.Fatalf("Failed to start job: %v", err)
	}

	stream := &mockStream{ctx: ctx}
	err = s.Stream(&pb.StreamRequest{JobId: resp.JobId}, stream)
	if err != nil {
		t.Fatalf("Stream failed: %v", err)
	}

	// Verify output
	foundStdout := false
	foundStderr := false
	for _, msg := range stream.messages {
		txt := string(msg.Data)
		if msg.Source == pb.LogChunk_STDOUT && strings.Contains(txt, "stdout-msg") {
			foundStdout = true
		}
		if msg.Source == pb.LogChunk_STDERR && strings.Contains(txt, "stderr-msg") {
			foundStderr = true
		}
	}

	if !foundStdout {
		t.Error("Did not find expected stdout message")
	}
	if !foundStderr {
		t.Error("Did not find expected stderr message")
	}
}

func TestWorkerServer_Validation(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())
	ctx := ctxWithUser(t.Context(), "alice")

	// Empty Command
	_, err := s.Start(ctx, &pb.StartRequest{Command: "", WorkingDirectory: "/tmp"})
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument for empty command, got: %v", err)
	}

	// Empty WorkingDirectory
	_, err = s.Start(ctx, &pb.StartRequest{Command: "echo", WorkingDirectory: ""})
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument for empty working directory, got: %v", err)
	}

	// Nil Request
	_, err = s.Start(ctx, nil)
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument for nil request, got: %v", err)
	}

	// Nil Stop Request
	if _, err := s.Stop(ctx, nil); status.Code(err) != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument for nil stop request, got: %v", err)
	}

	// Nil Status Request
	if _, err := s.Status(ctx, nil); status.Code(err) != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument for nil status request, got: %v", err)
	}

	// Nil Stream Request
	if err := s.Stream(nil, &mockStream{ctx: ctx}); status.Code(err) != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument for nil stream request, got: %v", err)
	}
}

func TestWorkerServer_NotFound(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())
	ctx := ctxWithUser(t.Context(), "alice")
	badID := "00000000-0000-0000-0000-000000000000"

	if _, err := s.Stop(ctx, &pb.StopRequest{JobId: badID}); status.Code(err) != codes.NotFound {
		t.Errorf("Expected NotFound for Stop, got: %v", err)
	}

	if _, err := s.Status(ctx, &pb.StatusRequest{JobId: badID}); status.Code(err) != codes.NotFound {
		t.Errorf("Expected NotFound for Status, got: %v", err)
	}

	stream := &mockStream{ctx: ctx}
	if err := s.Stream(&pb.StreamRequest{JobId: badID}, stream); status.Code(err) != codes.NotFound {
		t.Errorf("Expected NotFound for Stream, got: %v", err)
	}
}

func TestWorkerServer_Lifecycle_Failed(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())
	ctx := ctxWithUser(t.Context(), "alice")

	// Run a command that exits with error code
	resp, err := s.Start(ctx, &pb.StartRequest{
		Command:          "sh",
		Args:             []string{"-c", "exit 42"},
		WorkingDirectory: "/tmp",
	})
	if err != nil {
		t.Fatalf("Failed to start job: %v", err)
	}

	// Poll status until failed
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for job failure")
		case <-ticker.C:
			st, _ := s.Status(ctx, &pb.StatusRequest{JobId: resp.JobId})
			if st.Status == pb.StatusResponse_FAILED {
				if st.ExitCode != 42 {
					t.Errorf("Expected exit code 42, got %d", st.ExitCode)
				}
				return
			}
		}
	}
}

func TestWorkerServer_Stream_ContextCancel(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())
	ctx := ctxWithUser(t.Context(), "alice")

	// Start a long running job
	resp, err := s.Start(ctx, &pb.StartRequest{
		Command:          "sleep",
		Args:             []string{"10"},
		WorkingDirectory: "/tmp",
	})
	if err != nil {
		t.Fatalf("Failed to start job: %v", err)
	}

	// Create a context that cancels quickly
	cancelCtx, cancel := context.WithCancel(ctx)

	// Start streaming in a goroutine
	errChan := make(chan error)
	go func() {
		stream := &mockStream{ctx: cancelCtx}
		errChan <- s.Stream(&pb.StreamRequest{JobId: resp.JobId}, stream)
	}()

	// Cancel the context
	cancel()

	select {
	case <-errChan:
		// Success: Stream returned (likely nil) upon context cancellation
	case <-time.After(1 * time.Second):
		t.Fatal("Stream did not return after context cancellation")
	}
}

func TestWorkerServer_Lifecycle_Completed(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())
	ctx := ctxWithUser(t.Context(), "alice")

	// Run a command that exits successfully immediately
	resp, err := s.Start(ctx, &pb.StartRequest{
		Command:          "true",
		WorkingDirectory: "/tmp",
	})
	if err != nil {
		t.Fatalf("Failed to start job: %v", err)
	}

	// Poll status until completed
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for job completion")
		case <-ticker.C:
			st, _ := s.Status(ctx, &pb.StatusRequest{JobId: resp.JobId})
			if st.Status == pb.StatusResponse_COMPLETED {
				if st.ExitCode != 0 {
					t.Errorf("Expected exit code 0, got %d", st.ExitCode)
				}
				return
			}
		}
	}
}

func TestWorkerServer_Lifecycle_Stopped(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())
	ctx := ctxWithUser(t.Context(), "alice")

	// Start a long running job
	resp, err := s.Start(ctx, &pb.StartRequest{
		Command:          "sleep",
		Args:             []string{"10"},
		WorkingDirectory: "/tmp",
	})
	if err != nil {
		t.Fatalf("Failed to start job: %v", err)
	}

	// Stop the job
	if _, err := s.Stop(ctx, &pb.StopRequest{JobId: resp.JobId}); err != nil {
		t.Fatalf("Failed to stop job: %v", err)
	}

	// Verify status is STOPPED
	st, err := s.Status(ctx, &pb.StatusRequest{JobId: resp.JobId})
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}
	if st.Status != pb.StatusResponse_STOPPED {
		t.Errorf("Expected status STOPPED, got %v", st.Status)
	}
}

func TestWorkerServer_Start_ExecError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())
	ctx := ctxWithUser(t.Context(), "alice")

	// Try to start a non-existent command
	_, err := s.Start(ctx, &pb.StartRequest{
		Command:          "/bin/this_does_not_exist_12345",
		WorkingDirectory: "/tmp",
	})

	if status.Code(err) != codes.Internal {
		t.Errorf("Expected Internal error for exec failure, got: %v", err)
	}
}

func TestWorkerServer_Concurrency(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping user-mode test (simulated cgroups) when running as root")
	}
	s := NewWorkerServer(t.TempDir())
	ctx := ctxWithUser(t.Context(), "alice")

	// 1. Start a long-running job that produces output
	resp, err := s.Start(ctx, &pb.StartRequest{
		Command:          "sh",
		Args:             []string{"-c", "for i in $(seq 1 10); do echo $i; sleep 0.01; done"},
		WorkingDirectory: "/tmp",
	})
	if err != nil {
		t.Fatalf("Failed to start job: %v", err)
	}

	// 2. Concurrent Streaming: 10 clients streaming the same job
	var wg sync.WaitGroup
	clients := 10
	errs := make(chan error, clients)

	for i := 0; i < clients; i++ {
		wg.Go(func() {
			stream := &mockStream{ctx: ctx}
			if err := s.Stream(&pb.StreamRequest{JobId: resp.JobId}, stream); err != nil {
				errs <- err
			}
			// Verify we got some data
			if len(stream.messages) == 0 {
				errs <- status.Error(codes.Internal, "received no messages")
			}
		})
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("Concurrent stream failed: %v", err)
	}

	// 3. Verify Job is still accessible
	if _, err := s.Status(ctx, &pb.StatusRequest{JobId: resp.JobId}); err != nil {
		t.Errorf("Failed to get status after concurrent streaming: %v", err)
	}
}

func TestWorkerServer_ProtoMapping(t *testing.T) {
	s := &WorkerServer{}

	// Test toProtoStatus
	testsStatus := []struct {
		in  worker.JobStatus
		out pb.StatusResponse_JobStatus
	}{
		{worker.StatusRunning, pb.StatusResponse_RUNNING},
		{worker.StatusStopped, pb.StatusResponse_STOPPED},
		{worker.StatusCompleted, pb.StatusResponse_COMPLETED},
		{worker.StatusFailed, pb.StatusResponse_FAILED},
		{worker.StatusUnspecified, pb.StatusResponse_STATUS_UNSPECIFIED},
		{worker.JobStatus(999), pb.StatusResponse_STATUS_UNSPECIFIED},
	}

	for _, tc := range testsStatus {
		if got := s.toProtoStatus(tc.in); got != tc.out {
			t.Errorf("toProtoStatus(%v) = %v, want %v", tc.in, got, tc.out)
		}
	}

	// Test toProtoLogSource
	testsSource := []struct {
		in  worker.OutputSource
		out pb.LogChunk_Source
	}{
		{worker.OutputSourceStdout, pb.LogChunk_STDOUT},
		{worker.OutputSourceStderr, pb.LogChunk_STDERR},
		{worker.OutputSourceUnspecified, pb.LogChunk_SOURCE_UNSPECIFIED},
		{worker.OutputSource(999), pb.LogChunk_SOURCE_UNSPECIFIED},
	}

	for _, tc := range testsSource {
		if got := s.toProtoLogSource(tc.in); got != tc.out {
			t.Errorf("toProtoLogSource(%v) = %v, want %v", tc.in, got, tc.out)
		}
	}

	// Test toProtoLimits
	limits := worker.ResourceLimits{
		MemoryLimitBytes: 1024,
		CPULimitPercent:  0.5,
		IOLimit: worker.LinuxIOLimit{
			ReadBPS:  100,
			WriteBPS: 200,
		},
	}
	protoLimits := s.toProtoLimits(limits)
	if protoLimits.MemoryBytes != 1024 ||
		protoLimits.CpuPercent != 0.5 ||
		protoLimits.DiskReadBps != 100 ||
		protoLimits.DiskWriteBps != 200 {
		t.Errorf("toProtoLimits mismatch: got %v", protoLimits)
	}
}
