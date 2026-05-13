// Package worker provides job execution management with process isolation.
package worker

import (
	"context"
	"io"
	"sync"
)

// ErrReaderClosed is returned when reading from a closed reader.
var ErrReaderClosed = io.ErrClosedPipe

// MemoryBuffer is a thread-safe buffer that supports multiple concurrent readers.
// Each reader starts from offset 0 and is notified when new data is written.
// Uses per-reader channels for efficient notification.
//
// Thread Safety:
//   - Write() and Close() are safe for concurrent use
//   - Multiple readers can read concurrently
//   - Each reader must be used by a single goroutine (io.Reader convention)
type MemoryBuffer struct {
	mu        sync.Mutex
	content   []byte
	closed    bool
	closeOnce sync.Once

	// Per-reader notification channels
	readers map[*MemoryBufferReader]chan struct{}
}

// NewMemoryBuffer creates a new MemoryBuffer.
func NewMemoryBuffer() *MemoryBuffer {
	return &MemoryBuffer{
		content: make([]byte, 0, 4096),
		readers: make(map[*MemoryBufferReader]chan struct{}),
	}
}

// Write appends data to the buffer and notifies all waiting readers.
// Each reader is notified via its own channel.
// Lock is released before sending notifications to reduce contention.
// It implements io.Writer.
func (mb *MemoryBuffer) Write(p []byte) (int, error) {
	mb.mu.Lock()
	if mb.closed {
		mb.mu.Unlock()
		return 0, io.ErrClosedPipe
	}

	mb.content = append(mb.content, p...)

	// Copy channels to notify (release lock before sending)
	toNotify := make([]chan struct{}, 0, len(mb.readers))
	for _, ch := range mb.readers {
		toNotify = append(toNotify, ch)
	}
	mb.mu.Unlock()

	// Notify without holding lock (reduces contention)
	for _, ch := range toNotify {
		select {
		case ch <- struct{}{}:
		default: // Already has pending notification
		}
	}

	return len(p), nil
}

// Close marks the buffer as complete (process exited) and wakes all readers.
// Safe to call multiple times.
func (mb *MemoryBuffer) Close() error {
	mb.closeOnce.Do(func() {
		mb.mu.Lock()
		mb.closed = true

		// Copy channels to notify (release lock before sending)
		toNotify := make([]chan struct{}, 0, len(mb.readers))
		for _, ch := range mb.readers {
			toNotify = append(toNotify, ch)
		}
		mb.mu.Unlock()

		// Notify without holding lock
		for _, ch := range toNotify {
			select {
			case ch <- struct{}{}:
			default:
			}
		}
	})
	return nil
}

// readAt reads len(p) bytes from the buffer starting at the given offset.
// Internal method used by MemoryBufferReader.
func (mb *MemoryBuffer) readAt(p []byte, offset int64) (int, error) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if offset >= int64(len(mb.content)) {
		return 0, nil // No data available yet
	}

	n := copy(p, mb.content[offset:])
	return n, nil
}

// Size returns the current size of the buffer.
func (mb *MemoryBuffer) Size() int64 {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	return int64(len(mb.content))
}

// IsClosed returns whether the buffer has been closed.
func (mb *MemoryBuffer) IsClosed() bool {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	return mb.closed
}

// register adds a reader to receive notifications. Returns the notification channel.
func (mb *MemoryBuffer) register(r *MemoryBufferReader) chan struct{} {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	ch := make(chan struct{}, 1) // Buffered to avoid blocking Write()
	mb.readers[r] = ch
	return ch
}

// unregister removes a reader from notifications.
func (mb *MemoryBuffer) unregister(r *MemoryBufferReader) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	delete(mb.readers, r)
}

// NewReader creates an io.ReadCloser that streams from the buffer.
// The reader respects context cancellation — when ctx is cancelled,
// any blocked Read() returns immediately with ctx.Err().
// Calling Close() on the reader also unblocks any pending Read().
// Each reader starts at offset 0 and maintains its own position.
//
// Important: Each reader must be used by a single goroutine only.
// Concurrent reads from the same reader are not safe.
func (mb *MemoryBuffer) NewReader(ctx context.Context) *MemoryBufferReader {
	r := &MemoryBufferReader{
		buf: mb,
		ctx: ctx,
	}
	r.notify = mb.register(r)
	return r
}

// MemoryBufferReader reads from a MemoryBuffer starting at a specific offset.
// Each reader maintains its own position and reads independently.
// Uses per-reader channel for notifications
// Implements io.ReadCloser for compatibility with io.Copy and similar functions.
//
// Thread Safety: A single MemoryBufferReader must only be used by one goroutine
// at a time. This follows the standard io.Reader convention.
type MemoryBufferReader struct {
	buf    *MemoryBuffer
	ctx    context.Context
	notify chan struct{} // Per-reader notification channel

	mu     sync.Mutex
	offset int64 // Current read position
	closed bool  // Reader closed flag
}

// Read implements io.Reader. Blocks until data is available, buffer closes,
// reader is closed, or context is cancelled.
//
// Returns:
//   - (n, nil) when data is read successfully
//   - (0, io.EOF) when buffer is closed and all data has been read
//   - (0, ErrReaderClosed) when reader.Close() was called
//   - (0, context.Canceled/DeadlineExceeded) when context is cancelled
func (r *MemoryBufferReader) Read(p []byte) (int, error) {
	for {
		// Check stop conditions
		r.mu.Lock()
		if r.closed {
			r.mu.Unlock()
			return 0, ErrReaderClosed
		}
		currentOffset := r.offset
		r.mu.Unlock()

		if r.ctx.Err() != nil {
			return 0, r.ctx.Err()
		}

		// Try to read available data
		n, err := r.buf.readAt(p, currentOffset)
		if err != nil {
			return 0, err
		}
		if n > 0 {
			r.mu.Lock()
			r.offset += int64(n)
			r.mu.Unlock()
			return n, nil
		}

		// Check for EOF (buffer closed and we've read everything)
		if r.buf.IsClosed() && currentOffset >= r.buf.Size() {
			return 0, io.EOF
		}

		// Wait for notification (each reader has own channel)
		select {
		case <-r.notify:
			// New data or buffer closed, loop and check
		case <-r.ctx.Done():
			return 0, r.ctx.Err()
		}

		// Re-check closed after waking
		r.mu.Lock()
		if r.closed {
			r.mu.Unlock()
			return 0, ErrReaderClosed
		}
		r.mu.Unlock()
	}
}

// Close implements io.Closer. Unregisters from buffer notifications and
// unblocks any pending Read() calls. Safe to call multiple times.
func (r *MemoryBufferReader) Close() error {
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return nil // Already closed, idempotent
	}
	r.closed = true
	r.mu.Unlock()

	r.buf.unregister(r)

	// Send notification to unblock any waiting Read()
	select {
	case r.notify <- struct{}{}:
	default:
	}

	return nil
}

// Offset returns the current read position.
func (r *MemoryBufferReader) Offset() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.offset
}
