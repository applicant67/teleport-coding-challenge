package output

import (
	"context"
	"io"
	"os"
)

// Subscriber reads output from a Log.
//
// Key behaviors:
//   - Starts reading from offset 0
//   - Blocks when caught up and waiting for new data
type Subscriber struct {
	log    *Log
	file   *os.File
	offset int64
}

// Read returns the next frame.
//
// Blocking behavior:
//   - If data is available, returns immediately
//   - If caught up, blocks until: new data, log closed, or context canceled
//
// Returns io.EOF when log is closed and all data has been read.
func (s *Subscriber) Read(ctx context.Context) (*Frame, error) {
	for {
		// Only run if context is not canceled
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		// Get current state
		notify, committedLen, closed := s.log.getState()

		// Case 1: Data available - read it
		if s.offset < committedLen {
			// Seek before reading to guard against partial reads advancing the file cursor
			if _, err := s.file.Seek(s.offset, io.SeekStart); err != nil {
				return nil, err
			}

			frame, err := ReadFrame(s.file)
			if err != nil {
				return nil, err
			}

			s.offset += int64(frame.Size())
			return frame, nil
		}

		// Case 2: No data and log closed - done
		if closed {
			return nil, io.EOF
		}

		// Case 3: No data and log open - wait for notification
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-notify:
			// New data may be available, loop and check
			continue
		}
	}
}

// Close closes the subscriber's file handle.
func (s *Subscriber) Close() error {
	return s.file.Close()
}
