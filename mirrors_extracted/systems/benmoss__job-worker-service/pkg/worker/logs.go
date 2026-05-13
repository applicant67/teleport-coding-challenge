package worker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const defaultLogDir = "/var/lib/jobworker"

// logfilePath returns the path for a job's log file.
func logfilePath(logDir string, id uint64) string {
	if logDir == "" {
		logDir = defaultLogDir
	}
	return filepath.Join(logDir, fmt.Sprintf("job-%d", id), "output")
}

// createLogfileInput contains parameters for creating a log file.
type createLogfileInput struct {
	LogDir string
	JobID  uint64
}

// createLogfile creates the job-specific directory and log file.
// Returns an outputWriter, a cleanup function that closes the file and removes the job directory, and an error.
// If an error is returned, no cleanup is needed.
func createLogfile(input createLogfileInput) (*outputWriter, func() error, error) {
	logfilePath := logfilePath(input.LogDir, input.JobID)
	dir := filepath.Dir(logfilePath)

	// Use MkdirAll to create both the root log dir and job-specific dir
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("create directory %s: %w", dir, err)
	}

	logfile, err := os.OpenFile(logfilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		os.Remove(dir)
		return nil, nil, fmt.Errorf("create file %s: %w", logfilePath, err)
	}

	output := newOutputWriter(logfile)
	cleanup := func() error {
		closeErr := output.Close()
		removeErr := output.RemoveJobDir()
		return errors.Join(closeErr, removeErr)
	}

	return output, cleanup, nil
}

// outputWriter wraps a file and provides broadcast notifications for followers.
type outputWriter struct {
	file    *os.File
	mu      sync.Mutex
	notify  chan struct{} // closed and replaced on each write
	version uint64        // increments on each write
	done    bool          // true when process exits
}

// newOutputWriter creates a new outputWriter wrapping the given file.
func newOutputWriter(file *os.File) *outputWriter {
	return &outputWriter{
		file:   file,
		notify: make(chan struct{}),
	}
}

// Write writes data to the underlying file and broadcasts to followers.
func (w *outputWriter) Write(p []byte) (n int, err error) {
	n, err = w.file.Write(p)
	if n > 0 {
		w.broadcast()
	}
	return n, err
}

// broadcast notifies followers that new data is available.
func (w *outputWriter) broadcast() {
	w.mu.Lock()
	defer w.mu.Unlock()
	close(w.notify)
	w.notify = make(chan struct{})
	w.version++
}

// getNotifyState returns the current notification channel and version.
func (w *outputWriter) getNotifyState() (chan struct{}, uint64, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.notify, w.version, w.done
}

// markDone marks the output as complete and broadcasts one final time.
func (w *outputWriter) markDone() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.done = true
	close(w.notify)
}

// Close closes the underlying file.
func (w *outputWriter) Close() error {
	return w.file.Close()
}

// Name returns the path of the underlying file.
func (w *outputWriter) Name() string {
	return w.file.Name()
}

// RemoveJobDir removes the job directory containing the log file.
func (w *outputWriter) RemoveJobDir() error {
	return os.RemoveAll(filepath.Dir(w.file.Name()))
}

// followingReader reads from a file and blocks waiting for new data when follow mode is enabled.
type followingReader struct {
	ctx         context.Context
	file        *os.File
	output      *outputWriter
	lastVersion uint64
}

// Read implements io.Reader. It reads available data and blocks waiting for more if in follow mode.
func (r *followingReader) Read(p []byte) (int, error) {
	for {
		// Try to read from file
		n, err := r.file.Read(p)
		if n > 0 {
			return n, nil
		}

		// Check for non-EOF errors
		if err != nil && err != io.EOF {
			return 0, err
		}

		// We hit EOF. Check if new data is available or if output is done.
		notifyChan, version, done := r.output.getNotifyState()

		// If version changed since last check, new data may be available - retry immediately
		if version != r.lastVersion {
			r.lastVersion = version
			continue
		}

		// If done and we've read everything, return EOF
		if done {
			return 0, io.EOF
		}

		// No new data yet. Wait for notification or context cancellation.
		select {
		case <-notifyChan:
			// New data available, retry read
			continue
		case <-r.ctx.Done():
			return 0, r.ctx.Err()
		}
	}
}

// Close closes the underlying file.
func (r *followingReader) Close() error {
	return r.file.Close()
}
