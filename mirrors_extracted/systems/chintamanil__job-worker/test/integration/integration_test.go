//go:build integration

// Package integration contains end-to-end integration tests for the job-worker.
// These tests require Linux with cgroups v2 and root privileges.
//
// Run with: go test -v -tags=integration -count=1 ./test/integration/...
package integration

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/manil/job-worker/internal/server"
	"github.com/manil/job-worker/pkg/client"
	"github.com/manil/job-worker/pkg/worker"
)

const (
	// cgroupRoot is the base path for test cgroups.
	cgroupRoot = "/sys/fs/cgroup/jobworker-integration-test"

	// certsDir is the path to the certificates directory.
	certsDir = "../../certs"
)

// testEnv holds the test environment with server and client.
type testEnv struct {
	t          *testing.T
	server     *server.Server
	serverAddr string
	cgroupRoot string
}

// newTestEnv creates a new test environment with a running server.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	skipIfNotLinux(t)
	skipIfNotRoot(t)
	skipIfNoCgroups(t)
	skipIfNoCerts(t)

	// Create cgroup root for this test
	testCgroupRoot := filepath.Join(cgroupRoot, t.Name())
	if err := os.MkdirAll(testCgroupRoot, 0755); err != nil {
		t.Fatalf("failed to create cgroup root: %v", err)
	}

	// Enable controllers for child cgroups (best effort)
	subtreeControl := filepath.Join(testCgroupRoot, "cgroup.subtree_control")
	_ = os.WriteFile(subtreeControl, []byte("+cpu +memory +io"), 0644)

	// Create job store and server
	store := worker.NewJobStore()
	srv := server.NewServer(store, testCgroupRoot)

	// Load TLS config
	tlsConfig, err := server.NewTLSConfig(
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

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)

	// Check for startup errors
	select {
	case err := <-errCh:
		t.Fatalf("server failed to start: %v", err)
	default:
	}

	env := &testEnv{
		t:          t,
		server:     srv,
		serverAddr: addr,
		cgroupRoot: testCgroupRoot,
	}

	t.Cleanup(func() {
		srv.GracefulStop()
		os.RemoveAll(testCgroupRoot)
	})

	return env
}

// newClient creates a new client connected to the test server.
func (e *testEnv) newClient(certName string) *client.Client {
	e.t.Helper()

	c, err := client.New(
		e.serverAddr,
		filepath.Join(certsDir, certName+".pem"),
		filepath.Join(certsDir, certName+"-key.pem"),
		filepath.Join(certsDir, "ca.pem"),
	)
	if err != nil {
		e.t.Fatalf("failed to create client: %v", err)
	}

	e.t.Cleanup(func() {
		c.Close()
	})

	return c
}

// skipIfNotLinux skips the test if not running on Linux.
func skipIfNotLinux(t *testing.T) {
	t.Helper()
	if os.Getenv("GOOS") == "darwin" || (os.Getenv("GOOS") == "" && !isLinux()) {
		t.Skip("integration tests require Linux")
	}
}

// isLinux checks if running on Linux by checking for /proc/version.
func isLinux() bool {
	_, err := os.Stat("/proc/version")
	return err == nil
}

// skipIfNotRoot skips the test if not running as root.
func skipIfNotRoot(t *testing.T) {
	t.Helper()
	if os.Geteuid() != 0 {
		t.Skip("integration tests require root privileges")
	}
}

// skipIfNoCgroups skips the test if cgroups v2 is not available.
func skipIfNoCgroups(t *testing.T) {
	t.Helper()
	if _, err := os.Stat("/sys/fs/cgroup/cgroup.controllers"); os.IsNotExist(err) {
		t.Skip("integration tests require cgroups v2")
	}
}

// skipIfNoCerts skips the test if certificates are not available.
func skipIfNoCerts(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(certsDir, "server.pem")); os.IsNotExist(err) {
		t.Skip("certificates not found - run 'make certs' first")
	}
}

// skipIfNoStressNg skips the test if stress-ng is not available.
func skipIfNoStressNg(t *testing.T) {
	t.Helper()
	if _, err := os.Stat("/usr/bin/stress-ng"); os.IsNotExist(err) {
		t.Skip("stress-ng not installed - required for resource limit tests")
	}
}

// skipIfControllersNotDelegated skips if cgroup controllers aren't delegated.
// In containers/VMs, controllers like cpu, memory, io may not be writable.
func skipIfControllersNotDelegated(t *testing.T, controllers ...string) {
	t.Helper()

	// Create a test cgroup to check if we can write controller files
	testDir := filepath.Join(cgroupRoot, "controller-test-"+t.Name())
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Skipf("cannot create test cgroup: %v", err)
	}
	defer os.Remove(testDir)

	for _, ctrl := range controllers {
		var testFile, testValue string
		switch ctrl {
		case "cpu":
			testFile = "cpu.max"
			testValue = "50000 100000"
		case "memory":
			testFile = "memory.max"
			testValue = "134217728" // 128MB
		case "io":
			testFile = "io.weight"
			testValue = "default 100"
		default:
			continue
		}

		path := filepath.Join(testDir, testFile)
		if err := os.WriteFile(path, []byte(testValue), 0644); err != nil {
			t.Skipf("%s controller not delegated (cannot write %s): %v", ctrl, testFile, err)
		}
	}
}

// waitForStatus polls job status until it matches the expected status or times out.
func waitForStatus(c *client.Client, jobID, expectedStatus string, timeout time.Duration) error {
	ctx := context.Background()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		status, err := c.Status(ctx, jobID)
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}
		if status.Status == expectedStatus {
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for status %s", expectedStatus)
}
