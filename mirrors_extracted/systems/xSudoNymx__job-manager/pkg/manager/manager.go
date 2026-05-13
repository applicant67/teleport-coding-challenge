package manager

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/xSudoNymx/job-manager/pkg/job"
	"github.com/xSudoNymx/job-manager/pkg/output"
)

var (
	// ErrJobNotFound is returned when a job ID does not exist.
	ErrJobNotFound = errors.New("job not found")

	// ErrStorageUnavailable indicates required local storage is unavailable.
	ErrStorageUnavailable = errors.New("storage unavailable")

	// ErrStartFailed indicates a job failed to start.
	ErrStartFailed = errors.New("failed to start job")

	// ErrStopFailed indicates a job failed to stop.
	ErrStopFailed = errors.New("failed to stop job")

	// ErrSubscribeFailed indicates a job output subscription failed.
	ErrSubscribeFailed = errors.New("failed to subscribe job")

	// ErrShutdownFailed indicates the manager failed to shut down cleanly.
	ErrShutdownFailed = errors.New("failed to shutdown manager")
)

type entry struct {
	job *job.Job
	log *output.Log
}

// Manager coordinates multiple jobs and their output logs.
type Manager struct {
	mu sync.RWMutex

	jobs    map[string]*entry
	dataDir string
	logger  *slog.Logger
}

type Config struct {
	// DataDir is where job output logs are stored.
	// If empty, uses a temp directory.
	DataDir string
	Logger  *slog.Logger
}

func New(cfg Config) (*Manager, error) {
	dataDir := cfg.DataDir
	if dataDir == "" {
		dataDir = filepath.Join(os.TempDir(), "job-worker")
	}

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	logger := cfg.Logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	return &Manager{
		jobs:    make(map[string]*entry),
		dataDir: dataDir,
		logger:  logger,
	}, nil
}

// Start creates and starts a new job.
// Returns the job ID on success.
func (m *Manager) Start(exe string, args []string) (string, error) {
	id := uuid.New().String()

	logPath := filepath.Join(m.dataDir, id+".log")
	log, err := output.New(output.Config{Path: logPath})
	if err != nil {
		m.logger.Error("create output log failed", "job_id", id, "path", logPath, "error", err)
		return "", ErrStorageUnavailable
	}

	j := job.NewWithID(id, job.Config{
		Exe:    exe,
		Args:   args,
		Stdout: output.NewFrameWriter(log, output.TagStdout),
		Stderr: output.NewFrameWriter(log, output.TagStderr),
	})

	if err := j.Start(); err != nil {
		if closeErr := log.Close(); closeErr != nil {
			m.logger.Warn("failed to close log during cleanup", "job_id", id, "error", closeErr)
		}
		if removeErr := os.Remove(logPath); removeErr != nil {
			m.logger.Warn("failed to remove log file during cleanup", "path", logPath, "error", removeErr)
		}
		m.logger.Error("failed to start job", "error", err)
		return "", ErrStartFailed
	}

	e := &entry{job: j, log: log}

	m.mu.Lock()
	m.jobs[id] = e
	m.mu.Unlock()

	go func() {
		<-j.Done()
		if err := log.Close(); err != nil {
			m.logger.Warn("failed to close log", "job_id", id, "error", err)
		}
	}()

	return id, nil
}

// Stop terminates a running job.
func (m *Manager) Stop(id string) error {
	m.mu.RLock()
	e, ok := m.jobs[id]
	m.mu.RUnlock()

	if !ok {
		return ErrJobNotFound
	}

	err := e.job.Stop()
	if err != nil {
		m.logger.Error("failed to stop job", "job_id", id, "error", err)
		return ErrStopFailed
	}

	return nil
}

// Status returns the current status of a job.
func (m *Manager) Status(id string) (*job.Status, error) {
	m.mu.RLock()
	e, ok := m.jobs[id]
	m.mu.RUnlock()

	if !ok {
		return nil, ErrJobNotFound
	}

	return e.job.Status(), nil
}

// Subscribe creates an output subscriber for a job.
func (m *Manager) Subscribe(id string) (*output.Subscriber, error) {
	m.mu.RLock()
	e, ok := m.jobs[id]
	m.mu.RUnlock()

	if !ok {
		return nil, ErrJobNotFound
	}

	sub, err := e.log.Subscribe()
	if err != nil {
		m.logger.Error("failed to subscribe job", "job_id", id, "error", err)
		return nil, ErrSubscribeFailed
	}

	return sub, nil
}

// Shutdown stops all running jobs.
func (m *Manager) Shutdown() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for _, e := range m.jobs {
		if e.job.State() == job.Running {
			if err := e.job.Stop(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	m.dataDir = ""
	clear(m.jobs)

	if len(errs) > 0 {
		m.logger.Error("failed to shutdown jobs", "jobs", errs)
		return ErrShutdownFailed
	}
	return nil
}
