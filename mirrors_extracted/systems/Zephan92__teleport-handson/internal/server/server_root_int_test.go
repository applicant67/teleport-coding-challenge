//go:build linux

package server

import (
	"os"
	"path/filepath"
	"testing"

	pb "github.com/teleport-handson/api/v1"
)

// TestWorkerServer_Root_EndToEnd verifies the server works with real Cgroups.
// It runs only when the test process has root privileges.
func TestWorkerServer_Root_EndToEnd(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("Skipping root-mode end-to-end test")
	}

	// Use a dedicated cgroup root for server integration tests to avoid conflicts
	cgroupRoot := "/sys/fs/cgroup/teleport-worker-server-test"

	// Setup: Create root and enable controllers
	if err := os.MkdirAll(cgroupRoot, 0755); err != nil {
		t.Fatalf("Failed to create cgroup root: %v", err)
	}
	// Best effort to enable controllers (cpu, memory, io)
	// We ignore errors here because they might already be enabled or inherited.
	_ = os.WriteFile(filepath.Join(cgroupRoot, "cgroup.subtree_control"), []byte("+cpu +memory +io"), 0644)

	// Teardown: Remove the cgroup root after the test
	defer func() {
		_ = os.Remove(cgroupRoot)
	}()

	s := NewWorkerServer(cgroupRoot)
	ctx := ctxWithUser(t.Context(), "alice")

	// 1. Start a long-running job
	startReq := &pb.StartRequest{
		Command:          "sleep",
		Args:             []string{"100"},
		WorkingDirectory: "/tmp",
	}
	resp, err := s.Start(ctx, startReq)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// 2. Verify it is running
	st, err := s.Status(ctx, &pb.StatusRequest{JobId: resp.JobId})
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if st.Status != pb.StatusResponse_RUNNING {
		t.Errorf("Expected RUNNING, got %v", st.Status)
	}

	// 3. Stop the job
	// This triggers the real cgroup.kill mechanism via the worker library
	if _, err := s.Stop(ctx, &pb.StopRequest{JobId: resp.JobId}); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	// 4. Verify it is stopped
	st, err = s.Status(ctx, &pb.StatusRequest{JobId: resp.JobId})
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if st.Status != pb.StatusResponse_STOPPED {
		t.Errorf("Expected STOPPED, got %v", st.Status)
	}
}
