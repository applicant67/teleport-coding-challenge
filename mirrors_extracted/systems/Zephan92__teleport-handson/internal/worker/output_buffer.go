package worker

import (
	"context"
	"io"
	"os"
	"sync"
	"time"
)

// streamWriter adapts OutputBuffer to io.Writer with a fixed source.
type streamWriter struct {
	buf *OutputBuffer
	src OutputSource
}

// Write implements io.Writer by delegating to the underlying buffer with the configured source.
func (w *streamWriter) Write(p []byte) (int, error) {
	return w.buf.write(p, w.src)
}

// OutputBuffer implements a thread-safe, multi-consumer broadcaster for process output.
// It uses a pull-based iterator pattern (LogStream) to allow clients to read logs
// at their own pace without blocking the writer or allocating per-client goroutines.
type OutputBuffer struct {
	// mu protects the internal state of the buffer (entries, closed).
	mu sync.Mutex

	// cond is a condition variable used to wake up dormant client streamers when new data arrives.
	cond sync.Cond
	// entries is the in-memory store of all stdout/stderr produced by the process.
	// Note: This is currently unbounded. In a production environment, this should be
	// implemented as a ring buffer or have a maximum size limit to prevent OOM.
	entries []LogEntry
	// closed indicates the process has finished writing to the buffer.
	closed bool
}

// NewOutputBuffer initializes a new OutputBuffer.
func NewOutputBuffer() *OutputBuffer {
	b := &OutputBuffer{
		entries: make([]LogEntry, 0),
	}
	b.cond.L = &b.mu
	return b
}

// Stdout returns an io.Writer that tags all writes as OutputSourceStdout.
func (b *OutputBuffer) Stdout() io.Writer {
	return &streamWriter{buf: b, src: OutputSourceStdout}
}

// Stderr returns an io.Writer that tags all writes as OutputSourceStderr.
func (b *OutputBuffer) Stderr() io.Writer {
	return &streamWriter{buf: b, src: OutputSourceStderr}
}

// write is the internal implementation that handles tagging and broadcasting.
func (b *OutputBuffer) write(p []byte, src OutputSource) (int, error) {
	// Lock early to ensure the "closed" check and the append are atomic.
	// This prevents a race where the buffer is closed between the check and the write.
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return 0, os.ErrClosed
	}

	// We must copy the slice because the caller (os/exec) may reuse the backing array.
	data := make([]byte, len(p))
	copy(data, p)

	entry := LogEntry{
		Data:      data,
		Timestamp: time.Now().UTC(),
		Source:    src,
	}

	b.entries = append(b.entries, entry)
	b.cond.Broadcast()

	return len(p), nil
}

// Close marks the buffer as closed and wakes up all waiting streamers.
func (b *OutputBuffer) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.closed {
		b.closed = true
		b.cond.Broadcast()
	}
}

// LogStream is an iterator for reading log entries from an OutputBuffer.
// It maintains a cursor (offset) into the buffer and blocks in Next() until new data is available.
type LogStream struct {
	b      *OutputBuffer
	offset int
	ctx    context.Context
	stop   func() bool
}

// Close cleans up resources associated with the stream.
func (s *LogStream) Close() {
	if s.stop != nil {
		s.stop()
	}
}

// Next returns the next log entry, blocking if necessary.
// Returns io.EOF when the buffer is closed and drained.
func (s *LogStream) Next() (LogEntry, error) {
	s.b.mu.Lock()
	defer s.b.mu.Unlock()

	for {
		if err := s.ctx.Err(); err != nil {
			return LogEntry{}, err
		}

		if s.offset < len(s.b.entries) {
			entry := s.b.entries[s.offset]
			s.offset++
			return entry, nil
		}

		if s.b.closed {
			return LogEntry{}, io.EOF
		}

		s.b.cond.Wait()
	}
}

// Stream returns a LogStream iterator that allows clients to read the buffer's output.
func (b *OutputBuffer) Stream(ctx context.Context) *LogStream {
	s := &LogStream{b: b, ctx: ctx}
	// Register AfterFunc to wake up Wait() on context cancellation.
	s.stop = context.AfterFunc(ctx, func() {
		b.mu.Lock()
		b.cond.Broadcast()
		b.mu.Unlock()
	})
	return s
}
