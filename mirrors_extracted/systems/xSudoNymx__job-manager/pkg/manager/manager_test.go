package manager

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/xSudoNymx/job-manager/pkg/job"
	"github.com/xSudoNymx/job-manager/pkg/output"
)

func TestManager_StartAndStatus(t *testing.T) {
	m := newTestManager(t)

	id, err := m.Start("/bin/echo", []string{"hello"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if id == "" {
		t.Error("Start() returned empty ID")
	}

	// Wait for completion
	time.Sleep(100 * time.Millisecond)

	status, err := m.Status(id)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	if status.ID != id {
		t.Errorf("Status.ID = %q, want %q", status.ID, id)
	}
	if status.State != job.Exited {
		t.Errorf("Status.State = %v, want Exited", status.State)
	}
}

func TestManager_Stop(t *testing.T) {
	m := newTestManager(t)

	id, err := m.Start("/bin/sleep", []string{"60"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if err := m.Stop(id); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	// Wait for stop to complete
	time.Sleep(100 * time.Millisecond)

	status, _ := m.Status(id)
	if status.State != job.Stopped {
		t.Errorf("State = %v, want Stopped", status.State)
	}
}

func TestManager_Stop_NotFound(t *testing.T) {
	m := newTestManager(t)

	err := m.Stop("nonexistent-id")
	if !errors.Is(err, ErrJobNotFound) {
		t.Errorf("Stop() = %v, want ErrJobNotFound", err)
	}
}

func TestManager_Status_NotFound(t *testing.T) {
	m := newTestManager(t)

	_, err := m.Status("nonexistent-id")
	if !errors.Is(err, ErrJobNotFound) {
		t.Errorf("Status() = %v, want ErrJobNotFound", err)
	}
}

func TestManager_Subscribe(t *testing.T) {
	m := newTestManager(t)

	id, err := m.Start("/bin/echo", []string{"test output"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Wait for output
	time.Sleep(100 * time.Millisecond)

	sub, err := m.Subscribe(id)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	defer sub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	frame, err := sub.Read(ctx)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if frame.Tag != output.TagStdout {
		t.Errorf("Tag = %v, want stdout", frame.Tag)
	}
	if len(frame.Data) == 0 {
		t.Error("Data should not be empty")
	}
}

func TestManager_Subscribe_NotFound(t *testing.T) {
	m := newTestManager(t)

	_, err := m.Subscribe("nonexistent-id")
	if !errors.Is(err, ErrJobNotFound) {
		t.Errorf("Subscribe() = %v, want ErrJobNotFound", err)
	}
}

func TestManager_OutputStreaming(t *testing.T) {
	m := newTestManager(t)

	// Start a job that produces multiple lines
	id, err := m.Start("/bin/sh", []string{"-c", "echo line1; echo line2; echo line3"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Wait for completion
	time.Sleep(200 * time.Millisecond)

	sub, err := m.Subscribe(id)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	defer sub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Read all output
	var frames []*output.Frame
	for {
		frame, err := sub.Read(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
		frames = append(frames, frame)
	}

	// Should have at least one frame
	if len(frames) == 0 {
		t.Error("should have received output frames")
	}
}

func TestManager_Shutdown(t *testing.T) {
	m := newTestManager(t)

	_, err := m.Start("/bin/sleep", []string{"60"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	_, err = m.Start("/bin/sleep", []string{"60"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Shutdown stops all jobs - no error means success
	if err := m.Shutdown(); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}
}

func newTestManager(t *testing.T) *Manager {
	t.Helper()

	dir := t.TempDir()
	m, err := New(Config{
		DataDir: dir,
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Cleanup(func() {
		m.Shutdown()
	})

	return m
}
