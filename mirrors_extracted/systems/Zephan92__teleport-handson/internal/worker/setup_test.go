//go:build linux

package worker

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/unix"
)

func TestCheckCgroupMount(t *testing.T) {
	tests := []struct {
		name      string
		statError error
		fsType    int64
		wantErr   bool
	}{
		{"Valid", nil, unix.CGROUP2_SUPER_MAGIC, false},
		{"InvalidType", nil, 0xEF53, true}, // EXT4_SUPER_MAGIC
		{"StatError", fmt.Errorf("stat failed"), 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			origStatfs := statfs
			defer func() { statfs = origStatfs }()
			statfs = func(path string, buf *unix.Statfs_t) error {
				if tc.statError != nil {
					return tc.statError
				}
				buf.Type = tc.fsType
				return nil
			}

			err := checkCgroupMount(CgroupMountPoint)
			if (err != nil) != tc.wantErr {
				t.Errorf("checkCgroupMount() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestConfigureControllers(t *testing.T) {
	// Setup temp dir as cgroup root
	tmpDir := t.TempDir()
	workerRoot := filepath.Join(tmpDir, CgroupWorkerDir)

	// Create fake cgroup.controllers
	controllersFile := filepath.Join(tmpDir, "cgroup.controllers")
	if err := os.WriteFile(controllersFile, []byte("cpu memory io cpuset"), 0644); err != nil {
		t.Fatalf("failed to create controllers file: %v", err)
	}

	// Run configuration
	if err := configureControllers(tmpDir); err != nil {
		t.Fatalf("configureControllers() failed: %v", err)
	}

	// Verify root subtree_control
	rootControl := filepath.Join(tmpDir, "cgroup.subtree_control")
	content, err := os.ReadFile(rootControl)
	if err != nil {
		t.Errorf("failed to read root subtree_control: %v", err)
	}
	if string(content) != "+cpu +memory +io" {
		t.Errorf("root subtree_control = %q, want '+cpu +memory +io'", string(content))
	}

	// Verify worker directory created
	if _, err := os.Stat(workerRoot); os.IsNotExist(err) {
		t.Errorf("worker directory not created")
	}

	// Verify worker subtree_control
	workerControl := filepath.Join(workerRoot, "cgroup.subtree_control")
	content, err = os.ReadFile(workerControl)
	if err != nil {
		t.Errorf("failed to read worker subtree_control: %v", err)
	}
	if string(content) != "+cpu +memory +io" {
		t.Errorf("worker subtree_control = %q, want '+cpu +memory +io'", string(content))
	}
}

// TestConfigureControllers_Missing asserts that configureControllers returns an error
// if a required controller (e.g., io) is missing from cgroup.controllers.
func TestConfigureControllers_Missing(t *testing.T) {
	// Setup temp dir as cgroup root
	tmpDir := t.TempDir()

	// Create cgroup.controllers with missing "io"
	controllersFile := filepath.Join(tmpDir, "cgroup.controllers")
	if err := os.WriteFile(controllersFile, []byte("cpu memory"), 0644); err != nil {
		t.Fatalf("failed to create controllers file: %v", err)
	}

	// Run configuration
	err := configureControllers(tmpDir)
	if err == nil {
		t.Error("configureControllers() should fail when 'io' controller is missing")
	} else if err.Error() != fmt.Sprintf("missing controller io in %s", controllersFile) {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestConfigureControllers_ReadError asserts that configureControllers returns an error
// if the cgroup.controllers file cannot be read.
func TestConfigureControllers_ReadError(t *testing.T) {
	tmpDir := t.TempDir()

	// Do not create cgroup.controllers to force a read error
	if err := configureControllers(tmpDir); err == nil {
		t.Error("configureControllers() should fail when cgroup.controllers is missing")
	}
}

// TestSetupCgroups_RootCheck asserts that SetupCgroups fails fast if the process
// is not running with root privileges.
func TestSetupCgroups_RootCheck(t *testing.T) {
	origGeteuid := geteuid
	defer func() { geteuid = origGeteuid }()

	// Simulate non-root
	geteuid = func() int { return 1000 }

	if err := SetupCgroups(); err == nil {
		t.Error("SetupCgroups() should fail when not root")
	}
}

// TestSetupCgroups_Integration runs the full SetupCgroups flow on a real system
// (requires root) to verify end-to-end configuration.
func TestSetupCgroups_Integration(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("Skipping integration test; requires root")
	}
	if err := SetupCgroups(); err != nil {
		t.Errorf("SetupCgroups() failed: %v", err)
	}
}
