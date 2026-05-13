//go:build linux

package cgroup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	cg := New("/sys/fs/cgroup/jobworker", "test-job-123")

	if cg.id != "test-job-123" {
		t.Errorf("id = %q, want %q", cg.id, "test-job-123")
	}
	if cg.path != "/sys/fs/cgroup/jobworker/test-job-123" {
		t.Errorf("path = %q, want %q", cg.path, "/sys/fs/cgroup/jobworker/test-job-123")
	}
	if cg.fd != -1 {
		t.Errorf("fd = %d, want -1", cg.fd)
	}
}

func TestPath(t *testing.T) {
	cg := New("/sys/fs/cgroup/jobworker", "job-456")
	if cg.Path() != "/sys/fs/cgroup/jobworker/job-456" {
		t.Errorf("Path() = %q, want %q", cg.Path(), "/sys/fs/cgroup/jobworker/job-456")
	}
}

func TestCreate(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	limits := Limits{
		CPUQuota:    0.5,
		MemoryBytes: 512 * 1024 * 1024, // 512 MiB
		IOWeight:    200,
	}

	if err := cg.Create(limits); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(cg.path); os.IsNotExist(err) {
		t.Fatal("cgroup directory was not created")
	}

	// Verify cpu.max was written
	cpuMax, err := os.ReadFile(filepath.Join(cg.path, "cpu.max"))
	if err != nil {
		t.Fatalf("failed to read cpu.max: %v", err)
	}
	// 0.5 CPU = 50000 quota with 100000 period
	if string(cpuMax) != "50000 100000" {
		t.Errorf("cpu.max = %q, want %q", string(cpuMax), "50000 100000")
	}

	// Verify memory.max was written
	memMax, err := os.ReadFile(filepath.Join(cg.path, "memory.max"))
	if err != nil {
		t.Fatalf("failed to read memory.max: %v", err)
	}
	if string(memMax) != "536870912" {
		t.Errorf("memory.max = %q, want %q", string(memMax), "536870912")
	}

	// Verify io.weight was written
	ioWeight, err := os.ReadFile(filepath.Join(cg.path, "io.weight"))
	if err != nil {
		t.Fatalf("failed to read io.weight: %v", err)
	}
	if string(ioWeight) != "default 200" {
		t.Errorf("io.weight = %q, want %q", string(ioWeight), "default 200")
	}
}

func TestCreate_NoLimits(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	// All zeros = no limits
	limits := Limits{}

	if err := cg.Create(limits); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(cg.path); os.IsNotExist(err) {
		t.Fatal("cgroup directory was not created")
	}

	// Verify no controller files were written
	if _, err := os.Stat(filepath.Join(cg.path, "cpu.max")); !os.IsNotExist(err) {
		t.Error("cpu.max should not exist when CPUQuota = 0")
	}
	if _, err := os.Stat(filepath.Join(cg.path, "memory.max")); !os.IsNotExist(err) {
		t.Error("memory.max should not exist when MemoryBytes = 0")
	}
	if _, err := os.Stat(filepath.Join(cg.path, "io.weight")); !os.IsNotExist(err) {
		t.Error("io.weight should not exist when IOWeight = 0")
	}
}

func TestCreate_PartialLimits(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	// Only CPU limit
	limits := Limits{
		CPUQuota: 1.5, // 150% = 1.5 cores
	}

	if err := cg.Create(limits); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify cpu.max was written
	cpuMax, err := os.ReadFile(filepath.Join(cg.path, "cpu.max"))
	if err != nil {
		t.Fatalf("failed to read cpu.max: %v", err)
	}
	// 1.5 CPU = 150000 quota with 100000 period
	if string(cpuMax) != "150000 100000" {
		t.Errorf("cpu.max = %q, want %q", string(cpuMax), "150000 100000")
	}

	// Other files should not exist
	if _, err := os.Stat(filepath.Join(cg.path, "memory.max")); !os.IsNotExist(err) {
		t.Error("memory.max should not exist")
	}
	if _, err := os.Stat(filepath.Join(cg.path, "io.weight")); !os.IsNotExist(err) {
		t.Error("io.weight should not exist")
	}
}

func TestRemove(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	if err := cg.Create(Limits{}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := cg.Remove(); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	// Verify directory was removed
	if _, err := os.Stat(cg.path); !os.IsNotExist(err) {
		t.Error("cgroup directory should be removed")
	}
}

func TestOpenFD_CloseFD(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	if err := cg.Create(Limits{}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	fd, err := cg.OpenFD()
	if err != nil {
		t.Fatalf("OpenFD() error = %v", err)
	}
	if fd < 0 {
		t.Errorf("OpenFD() returned invalid fd = %d", fd)
	}
	if cg.fd != fd {
		t.Errorf("cg.fd = %d, want %d", cg.fd, fd)
	}

	// Close should succeed
	if err := cg.CloseFD(); err != nil {
		t.Fatalf("CloseFD() error = %v", err)
	}
	if cg.fd != -1 {
		t.Errorf("cg.fd = %d after close, want -1", cg.fd)
	}

	// Double close should be safe
	if err := cg.CloseFD(); err != nil {
		t.Fatalf("second CloseFD() error = %v", err)
	}
}

func TestKill(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	if err := cg.Create(Limits{}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Kill writes "1" to cgroup.kill
	if err := cg.Kill(); err != nil {
		t.Fatalf("Kill() error = %v", err)
	}

	// Verify cgroup.kill was written
	killFile, err := os.ReadFile(filepath.Join(cg.path, "cgroup.kill"))
	if err != nil {
		t.Fatalf("failed to read cgroup.kill: %v", err)
	}
	if string(killFile) != "1" {
		t.Errorf("cgroup.kill = %q, want %q", string(killFile), "1")
	}
}

func TestCreate_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	// Pre-create the directory
	if err := os.Mkdir(cg.path, 0755); err != nil {
		t.Fatalf("failed to pre-create directory: %v", err)
	}

	// Create should still succeed (idempotent for directory)
	if err := cg.Create(Limits{CPUQuota: 0.5}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
}

// --- Unhappy path tests ---

func TestValidate_InvalidCPUQuota(t *testing.T) {
	limits := Limits{CPUQuota: -0.5}
	err := limits.Validate()
	if err != ErrInvalidCPUQuota {
		t.Errorf("Validate() error = %v, want %v", err, ErrInvalidCPUQuota)
	}
}

func TestValidate_InvalidMemoryLimit(t *testing.T) {
	limits := Limits{MemoryBytes: -1024}
	err := limits.Validate()
	if err != ErrInvalidMemoryLimit {
		t.Errorf("Validate() error = %v, want %v", err, ErrInvalidMemoryLimit)
	}
}

func TestValidate_InvalidIOWeight_TooLow(t *testing.T) {
	limits := Limits{IOWeight: 0} // 0 is valid (means default)
	if err := limits.Validate(); err != nil {
		t.Errorf("Validate() with IOWeight=0 should succeed, got %v", err)
	}

	limits = Limits{IOWeight: -1}
	err := limits.Validate()
	if err != ErrInvalidIOWeight {
		t.Errorf("Validate() error = %v, want %v", err, ErrInvalidIOWeight)
	}
}

func TestValidate_InvalidIOWeight_TooHigh(t *testing.T) {
	limits := Limits{IOWeight: 10001}
	err := limits.Validate()
	if err != ErrInvalidIOWeight {
		t.Errorf("Validate() error = %v, want %v", err, ErrInvalidIOWeight)
	}
}

func TestValidate_ValidIOWeightBoundaries(t *testing.T) {
	// Test valid boundary values
	tests := []struct {
		weight int32
	}{
		{0},     // default
		{1},     // min
		{100},   // typical default
		{10000}, // max
	}

	for _, tt := range tests {
		limits := Limits{IOWeight: tt.weight}
		if err := limits.Validate(); err != nil {
			t.Errorf("Validate() with IOWeight=%d should succeed, got %v", tt.weight, err)
		}
	}
}

func TestCreate_InvalidLimits(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	// Invalid limits should fail validation before creating directory
	limits := Limits{CPUQuota: -1}
	err := cg.Create(limits)
	if err != ErrInvalidCPUQuota {
		t.Errorf("Create() error = %v, want %v", err, ErrInvalidCPUQuota)
	}

	// Directory should not be created
	if _, err := os.Stat(cg.path); !os.IsNotExist(err) {
		t.Error("cgroup directory should not be created for invalid limits")
	}
}

func TestOpenFD_NotExists(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "nonexistent-job")

	// OpenFD should fail for non-existent directory
	_, err := cg.OpenFD()
	if err == nil {
		t.Error("OpenFD() should fail for non-existent directory")
	}
}

func TestKill_NotExists(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "nonexistent-job")

	// Kill should fail for non-existent directory
	err := cg.Kill()
	if err == nil {
		t.Error("Kill() should fail for non-existent directory")
	}
}

func TestRemove_NotExists(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "nonexistent-job")

	// Remove should fail for non-existent directory
	err := cg.Remove()
	if err == nil {
		t.Error("Remove() should fail for non-existent directory")
	}
}

func TestRemove_NotEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	if err := cg.Create(Limits{}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Create a file inside to simulate non-empty cgroup
	dummyFile := filepath.Join(cg.path, "dummy")
	if err := os.WriteFile(dummyFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}

	// Remove should fail because directory is not empty
	err := cg.Remove()
	if err == nil {
		t.Error("Remove() should fail for non-empty directory")
	}
}

func TestIsEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	if err := cg.Create(Limits{}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Simulate cgroup.events file with "populated 0" (empty)
	eventsFile := filepath.Join(cg.path, "cgroup.events")
	if err := os.WriteFile(eventsFile, []byte("populated 0\nfrozen 0\n"), 0644); err != nil {
		t.Fatalf("failed to write cgroup.events: %v", err)
	}

	empty, err := cg.IsEmpty()
	if err != nil {
		t.Fatalf("IsEmpty() error = %v", err)
	}
	if !empty {
		t.Error("IsEmpty() = false, want true")
	}
}

func TestIsEmpty_NotEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	if err := cg.Create(Limits{}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Simulate cgroup.events file with "populated 1" (has processes)
	eventsFile := filepath.Join(cg.path, "cgroup.events")
	if err := os.WriteFile(eventsFile, []byte("populated 1\nfrozen 0\n"), 0644); err != nil {
		t.Fatalf("failed to write cgroup.events: %v", err)
	}

	empty, err := cg.IsEmpty()
	if err != nil {
		t.Fatalf("IsEmpty() error = %v", err)
	}
	if empty {
		t.Error("IsEmpty() = true, want false")
	}
}

func TestIsEmpty_NoEventsFile(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	if err := cg.Create(Limits{}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Don't create cgroup.events file
	_, err := cg.IsEmpty()
	if err == nil {
		t.Error("IsEmpty() should fail when cgroup.events doesn't exist")
	}
}

// --- Nil-safety tests ---

func TestNilCgroup_Path(t *testing.T) {
	var cg *Cgroup
	if cg.Path() != "" {
		t.Errorf("Path() on nil = %q, want empty string", cg.Path())
	}
}

func TestNilCgroup_Create(t *testing.T) {
	var cg *Cgroup
	err := cg.Create(Limits{})
	if err != ErrNilCgroup {
		t.Errorf("Create() on nil = %v, want %v", err, ErrNilCgroup)
	}
}

func TestNilCgroup_OpenFD(t *testing.T) {
	var cg *Cgroup
	fd, err := cg.OpenFD()
	if err != ErrNilCgroup {
		t.Errorf("OpenFD() on nil error = %v, want %v", err, ErrNilCgroup)
	}
	if fd != -1 {
		t.Errorf("OpenFD() on nil fd = %d, want -1", fd)
	}
}

func TestNilCgroup_CloseFD(t *testing.T) {
	var cg *Cgroup
	// CloseFD on nil should be a no-op (safe)
	if err := cg.CloseFD(); err != nil {
		t.Errorf("CloseFD() on nil = %v, want nil", err)
	}
}

func TestNilCgroup_Kill(t *testing.T) {
	var cg *Cgroup
	err := cg.Kill()
	if err != ErrNilCgroup {
		t.Errorf("Kill() on nil = %v, want %v", err, ErrNilCgroup)
	}
}

func TestNilCgroup_Remove(t *testing.T) {
	var cg *Cgroup
	err := cg.Remove()
	if err != ErrNilCgroup {
		t.Errorf("Remove() on nil = %v, want %v", err, ErrNilCgroup)
	}
}

func TestNilCgroup_IsEmpty(t *testing.T) {
	var cg *Cgroup
	_, err := cg.IsEmpty()
	if err != ErrNilCgroup {
		t.Errorf("IsEmpty() on nil = %v, want %v", err, ErrNilCgroup)
	}
}

func TestNilCgroup_Cleanup(t *testing.T) {
	var cg *Cgroup
	// Cleanup on nil should be a no-op (safe)
	if err := cg.Cleanup(); err != nil {
		t.Errorf("Cleanup() on nil = %v, want nil", err)
	}
}

// --- Cleanup tests ---

func TestCleanup(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	if err := cg.Create(Limits{}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Open FD to test that Cleanup closes it
	_, err := cg.OpenFD()
	if err != nil {
		t.Fatalf("OpenFD() error = %v", err)
	}

	// Cleanup should close FD and remove directory
	if err := cg.Cleanup(); err != nil {
		t.Fatalf("Cleanup() error = %v", err)
	}

	// Verify FD is closed
	if cg.fd != -1 {
		t.Errorf("cg.fd = %d after Cleanup, want -1", cg.fd)
	}

	// Verify directory is removed
	if _, err := os.Stat(cg.path); !os.IsNotExist(err) {
		t.Error("cgroup directory should be removed after Cleanup")
	}
}

func TestCleanup_CollectsErrors(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	if err := cg.Create(Limits{}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Create a file inside to make Remove fail
	dummyFile := filepath.Join(cg.path, "dummy")
	if err := os.WriteFile(dummyFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}

	// Cleanup should return error (remove fails)
	err := cg.Cleanup()
	if err == nil {
		t.Error("Cleanup() should return error when Remove fails")
	}

	// FD should still be closed (best-effort cleanup)
	if cg.fd != -1 {
		t.Errorf("cg.fd = %d after Cleanup, want -1 (FD should be closed even if Remove fails)", cg.fd)
	}
}

func TestCleanup_AlreadyClosed(t *testing.T) {
	tmpDir := t.TempDir()
	cg := New(tmpDir, "test-job")

	if err := cg.Create(Limits{}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Cleanup without ever opening FD should still work
	if err := cg.Cleanup(); err != nil {
		t.Fatalf("Cleanup() error = %v", err)
	}

	// Directory should be removed
	if _, err := os.Stat(cg.path); !os.IsNotExist(err) {
		t.Error("cgroup directory should be removed after Cleanup")
	}
}
