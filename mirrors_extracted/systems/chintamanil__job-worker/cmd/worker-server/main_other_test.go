//go:build !linux

package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func buildTestBinary(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("go", "build", "-o", "worker-server-test", ".")
	cmd.Dir = "."
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build: %v\n%s", err, output)
	}
	return "worker-server-test"
}

func TestNonLinuxStub_Help(t *testing.T) {
	binary := buildTestBinary(t)
	defer os.Remove(binary)

	// Test --help works
	helpCmd := exec.Command("./"+binary, "--help")
	output, _ := helpCmd.CombinedOutput() // --help may exit non-zero in some cases

	// Verify help contains expected content
	if !bytes.Contains(output, []byte("Job Worker Server")) {
		t.Errorf("help output missing 'Job Worker Server', got: %s", output)
	}
	if !bytes.Contains(output, []byte("--listen")) {
		t.Errorf("help output missing '--listen' flag, got: %s", output)
	}
	if !bytes.Contains(output, []byte("--cert")) {
		t.Errorf("help output missing '--cert' flag, got: %s", output)
	}
	if !bytes.Contains(output, []byte("--cgroup-root")) {
		t.Errorf("help output missing '--cgroup-root' flag, got: %s", output)
	}
}

func TestNonLinuxStub_Version(t *testing.T) {
	binary := buildTestBinary(t)
	defer os.Remove(binary)

	versionCmd := exec.Command("./"+binary, "--version")
	output, _ := versionCmd.CombinedOutput()

	if !bytes.Contains(output, []byte("dev")) {
		t.Errorf("version output missing 'dev': %s", output)
	}
}

func TestNonLinuxStub_RunError(t *testing.T) {
	binary := buildTestBinary(t)
	defer os.Remove(binary)

	// Running without args should fail with Linux requirement error
	runCmd := exec.Command("./" + binary)
	output, err := runCmd.CombinedOutput()

	if err == nil {
		t.Error("expected error when running on non-Linux")
	}

	if !bytes.Contains(output, []byte("requires Linux")) {
		t.Errorf("expected 'requires Linux' error, got: %s", output)
	}
}
